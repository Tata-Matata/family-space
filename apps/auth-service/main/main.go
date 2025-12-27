package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	api "github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/http"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/jwt"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/password"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/refresh"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/service"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage/sqlite"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (using real env)")
	}

	// Load RSA private key for JWT signing
	path := os.Getenv("JWT_PRIVATE_KEY_PATH")
	privateKey, err := jwt.LoadRSAPrivateKey(path)

	if err != nil {
		log.Fatalf("failed to load JWT private key: %v", err)
	}
	signer := jwt.NewRS256Signer(
		privateKey,
		"family-space-auth", // iss
		"family-space-api",  // aud
		15*time.Minute,
	)

	// Initialize SQLite storage
	db, err := sqlite.Open("auth.db")
	if err != nil {
		log.Fatal(err)
	}
	transactionMgr := storage.NewTransactionMgr(db)

	// REGISTER SERVICE
	hasher := password.NewBcryptHasher(0)

	registrationService := service.NewRegistrationService(
		transactionMgr,
		hasher,
		//function reference
		sqlite.NewUserStore,
		sqlite.NewFamilyStore,
		sqlite.NewMembershipStore,
	)
	registerHandler := api.NewRegisterHandler(
		registrationService,
	)

	// LOGIN SERVICE
	loginService := service.NewLoginService(
		transactionMgr,
		hasher,

		sqlite.NewUserStore,
		sqlite.NewMembershipStore,
		signer,
	)
	loginHandler := api.NewLoginHandler(
		loginService,
		15*time.Minute,
	)

	// REFRESH SERVICE
	refreshKey := os.Getenv("REFRESH_TOKEN_HMAC_KEY")
	if refreshKey == "" {
		log.Fatal("REFRESH_TOKEN_HMAC_KEY must be set")
	}
	refreshService := service.NewRefreshService(
		transactionMgr,
		sqlite.NewRefreshTokenStore,
		sqlite.NewUserStore,
		sqlite.NewMembershipStore,
		refresh.NewHMACRefreshTokenHasher([]byte(refreshKey)),
		&refresh.SecureRefreshTokenGenerator{},
		signer,
		30*24*time.Hour,
	)
	refreshHandler := api.NewRefreshHandler(
		refreshService,
		15*time.Minute,
	)

	// SETUP HTTP SERVER
	mux := http.NewServeMux()
	mux.Handle("/register", registerHandler)
	mux.Handle("/login", loginHandler)

	mux.Handle("/refresh", refreshHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())

}

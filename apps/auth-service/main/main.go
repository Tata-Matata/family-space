package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	api "github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/http"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/jwt"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/password"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/refresh"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/service"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage/sqlite"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

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
	var db *sql.DB

	switch os.Getenv("DB_DRIVER") {
	case "sqlite":
		db, err = initSqlite()
	case "postgres":
		db, err = initPostgres()

	default:
		log.Fatal("DB_DRIVER must be set to select database (sqlite or postgres)")
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
	mux.Handle("/health", api.NewHealthHandler())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())

}

func initSqlite() (*sql.DB, error) {

	var db *sql.DB

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH not set")
	}

	db, err := sqlite.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("using SQLite database")
	return db, nil
}

func initPostgres() (*sql.DB, error) {
	var db *sql.DB
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("using Postgres database")
	return db, nil
}

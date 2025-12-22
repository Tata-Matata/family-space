package server

/*
privateKey, err := jwt.LoadRSAPrivateKey("keys/jwt_private.pem")
if err != nil {
	log.Fatalf("failed to load JWT private key: %v", err)
}
signer := jwt.NewRS256Signer(
	privateKey,
	"family-space-auth",   // iss
	"family-space-api",    // aud
	15*time.Minute,
)

db, err := sqlite.Open("auth.db")
if err != nil {
    log.Fatal(err)
}

registrationService := service.NewRegistrationService(
    db,
    //function reference
    sqlite.NewUserStore,
) */

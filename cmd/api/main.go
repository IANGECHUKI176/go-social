package main

import (
	"gopher_social/internal/auth"
	"gopher_social/internal/db"
	"gopher_social/internal/env"
	"gopher_social/internal/mailer"
	"gopher_social/internal/store"
	"time"

	"go.uber.org/zap"
)

const version = "0.0.2"

//	@title	Gopher Social API

//	@description	API for GopherSocial : a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

func main() {

	cfg := config{
		addr:        env.GetString("ADDR", ":8000"),
		apiURL:      env.GetString("API_URL", "localhost:8081"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/gopher_social?sslmode=disable"),

			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("BASIC_AUTH_USER", "admin"),
				pass: env.GetString("BASIC_AUTH_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "secret"),
				exp:    time.Hour * 24 * 3,
				iss:    "gophersocial",
			},
		},
		version: version,
	}
	//Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("✅ Connected to database")

	store := store.NewPostgresStorage(db)

	// Mailer
	// mailer := mailer.NewSendgridMailer(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	mailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	JWTAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)
	if err != nil {
		logger.Fatal(err)
	}
	app := application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailtrap,
		authenticator: JWTAuthenticator,
	}
	mux := app.mount()

	logger.Fatal(app.run(mux))
}

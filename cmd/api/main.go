package main

import (
	"log"
	"time"

	"github.com/alejandro-cardenas-g/social/internal/db"
	"github.com/alejandro-cardenas-g/social/internal/env"
	"github.com/alejandro-cardenas-g/social/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Gopher social API specs
//	@description	API for social network

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:    env.GetString("ADDR", ":8080"),
		apiHost: env.GetString("HOST_API", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:password@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3,
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)

	if err != nil {
		log.Panic(err)
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	defer db.Close()
	logger.Info("DB Connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))

}

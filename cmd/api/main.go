package main

import (
	"expvar"
	"log"
	"runtime"
	"time"

	"github.com/alejandro-cardenas-g/social/internal/auth"
	"github.com/alejandro-cardenas-g/social/internal/db"
	"github.com/alejandro-cardenas-g/social/internal/env"
	"github.com/alejandro-cardenas-g/social/internal/mailer"
	ratelimiter "github.com/alejandro-cardenas-g/social/internal/rateLimiter"
	"github.com/alejandro-cardenas-g/social/internal/store"
	"github.com/alejandro-cardenas-g/social/internal/store/cache"
	"github.com/go-redis/redis/v8"
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
		addr:        env.GetString("ADDR", ":8080"),
		apiHost:     env.GetString("HOST_API", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:password@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: SendGridConfig{
				apikey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user:     env.GetString("AUTH_BASIC_USER", ""),
				password: env.GetString("AUTH_BASIC_PASSWORD", ""),
			},
			token: tokenConfig{
				secret: env.GetString("TOKEN_SECRET", "exampleToken"),
				exp:    time.Hour * 24 * 3,
				iss:    "socalPostsApp",
			},
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATE_LIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
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

	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("redis cache connection established")
	}

	cacheStorage := cache.NewRedisStorage(rdb)

	store := store.NewStorage(db)

	mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apikey, cfg.mail.fromEmail)

	jwtAuthenticator := auth.NewJWtAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	rateLimiter := ratelimiter.NewFixedWindowRateLimiter(cfg.rateLimiter.RequestsPerTimeFrame, cfg.rateLimiter.TimeFrame)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		cacheStorage:  cacheStorage,
		rateLimiter:   rateLimiter,
	}

	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("routines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	log.Fatal(app.run(mux))

}

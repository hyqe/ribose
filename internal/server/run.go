package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/hyqe/ribose/internal/database"
	"github.com/hyqe/ribose/internal/fit"
	"github.com/hyqe/ribose/internal/users"
)

func Run(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	log.SetFlags(0)

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := connectDB(ctx, cfg.PostgresURL, cfg.MigrationsURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return cfg.Env == "DEV"
		},
	}))
	app.Use(compress.New())

	queries := database.New(db)

	userSvc := users.NewService(queries)

	fit.NewRPC(userSvc).MountFiberApp(app)

	go app.Listen(cfg.Addr())

	<-ctx.Done()
	log.Println("shutting down")
	err = app.ShutdownWithTimeout(time.Second * 10)
	if err != nil {
		log.Fatalf("app.ShutdownWithTimeout(): %v", err)
	}
}

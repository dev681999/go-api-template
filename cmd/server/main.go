package main

import (
	"context"
	"flag"
	"fmt"
	"go-api-template/internal/config"
	"go-api-template/internal/mail"
	"go-api-template/internal/mailer"
	"go-api-template/internal/openapi"
	"go-api-template/internal/security"
	"go-api-template/internal/transport"
	"go-api-template/internal/user"
	"go-api-template/pkg/db"
	"go-api-template/pkg/log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"
)

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()

	logger := log.Setup()

	cfg, err := config.New(*flagConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	logger.Debug().Msgf("%+v", cfg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db, err := db.Connect(ctx, logger.With().Str("layer", "database").Logger(), cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	ms := mail.NewService(logger, mailer.NewService(logger, "admin@admin.com", "Admin"))

	userRepo := user.NewRepository(logger.With().Str("svc", "user").Str("layer", "repo").Logger(), db)
	userSvc := user.NewService(
		logger.With().Str("svc", "user").Str("layer", "service").Logger(),
		userRepo,
		ms,
		cfg.Server.JWTKey,
		cfg.Server.ActivationURL,
	)
	userTransport := user.NewTransport(logger.With().Str("svc", "user").Str("layer", "transport").Logger(), userSvc)

	swagger, err := openapi.GetSwagger()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	swagger.Servers = nil

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	e := transport.NewEchoEngine(logger)

	apiGroup := e.Group("/api/v1")

	apiGroup.Use(openapi.NewPrefixEchoMiddleware("/api/v1"))
	apiGroup.Use(security.ValidationMiddleware(swagger, userSvc))

	openapi.RegisterHandlers(apiGroup, transport.New(userTransport))

	/* router := openapi3filter.NewRouter().WithSwagger(swagger)

	r, d, err := router.FindRoute("POST", &url.URL{
		Scheme: "http",
		Host:   "http://localhost:8000",
		Path:   "/api/v1/user/register",
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	logger.Debug().Msgf("%v %v", r, d) */

	var g run.Group
	{
		g.Add(func() error {
			logger.Info().Str("msg", "serving http").Str("addr", addr).Msg("server")

			return e.Start(addr)
		}, func(error) {
			logger.Info().Str("msg", "stopping http").Msg("server")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := e.Shutdown(ctx); err != nil {
				e.Logger.Fatal(err)
			}

			logger.Info().Str("msg", "stopping database connection").Msg("server")
			if err := db.Close(); err != nil {
				e.Logger.Fatal(err)
			}
		})
	}
	{
		// set-up our signal handler
		var (
			cancelInterrupt = make(chan struct{})
			c               = make(chan os.Signal, 2)
		)
		defer close(c)

		g.Add(func() error {
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	logger.Err(g.Run()).Msg("exit")
}

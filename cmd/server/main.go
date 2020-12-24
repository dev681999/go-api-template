package main

import (
	"context"
	"flag"
	"fmt"
	"go-api-template/internal/config"
	"go-api-template/internal/constants"
	"go-api-template/internal/mail"
	"go-api-template/internal/mailer"
	"go-api-template/internal/openapi"
	"go-api-template/internal/security"
	"go-api-template/internal/storage"
	"go-api-template/internal/transport"
	"go-api-template/internal/user"
	"go-api-template/pkg/db"
	"go-api-template/pkg/log"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"github.com/oklog/run"
	"github.com/rs/zerolog"
)

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func getLoggerWithSvcAndLayer(logger zerolog.Logger, svc string, layer string) zerolog.Logger {
	return logger.With().Str("svc", svc).Str("layer", layer).Logger()
}

func getLoggerWithLayer(logger zerolog.Logger, layer string) zerolog.Logger {
	return logger.With().Str("layer", layer).Logger()
}

func main() {
	flag.Parse()

	logger := log.Setup(log.EnvLocal)

	cfg, err := config.New(*flagConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	logger.Debug().Msgf("config: %+v", cfg)

	backend := s3mem.New()
	faker := gofakes3.New(backend)
	ts := httptest.NewServer(faker.Server())
	defer ts.Close()

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("key", "secret", ""),
		Endpoint:         aws.String(ts.URL),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db, err := db.Connect(ctx, getLoggerWithLayer(logger, constants.LayerDatabase), cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	ms := mail.NewService(logger, mailer.NewService(logger, "admin@admin.com", "Admin"))

	userRepo := user.NewRepository(getLoggerWithSvcAndLayer(logger, constants.SvcUser, constants.LayerRepo), db)
	userSvc := user.NewService(
		getLoggerWithSvcAndLayer(logger, constants.SvcUser, constants.LayerService),
		userRepo,
		ms,
		cfg.Server.JWTKey,
		cfg.Server.ActivationURL,
	)
	userTransport := user.NewTransport(getLoggerWithSvcAndLayer(logger, constants.SvcUser, constants.LayerTransport), userSvc)

	storageRepo := storage.NewRepository(getLoggerWithSvcAndLayer(logger, constants.SvcStorage, constants.LayerRepo), db)
	storageSvc := storage.NewService(
		getLoggerWithSvcAndLayer(logger, constants.SvcUser, constants.LayerService),
		storageRepo,
		s3Client,
	)
	storageTransport := storage.NewTransport(getLoggerWithSvcAndLayer(logger, constants.SvcStorage, constants.LayerTransport), storageSvc)

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

	openapi.RegisterHandlers(apiGroup, transport.New(userTransport, storageTransport))

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

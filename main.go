package main

import (
	"context"
	"embed"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/huzaifa678/Crypto-currency-web-app-project/api"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/gapi"
	"github.com/huzaifa678/Crypto-currency-web-app-project/mail"
	"github.com/huzaifa678/Crypto-currency-web-app-project/oauth2"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/huzaifa678/Crypto-currency-web-app-project/worker"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

//go:embed docs/*
var docsFS embed.FS

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {

	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	oauth2.InitGoogleOAuth(config)

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop() 

	connPool, err := pgxpool.New(ctx, config.Dbsource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.Dbsource)

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddr,
	}


	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	waitGroup, ctx := errgroup.WithContext(ctx)
	
	runTaskProcessor(ctx, waitGroup, config, redisOpt, store)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor)
	runGrpcServer(ctx, waitGroup, config, store,taskDistributor)

	err = waitGroup.Wait()

	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, config utils.Config, redisOpt asynq.RedisClientOpt, store db.Store_interface) {
	mailer := mail.NewGmailSender(config.SenderName, config.SenderEmail, config.SenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

	log.Info().Msg("Task processor started successfully")

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown task processor")

		taskProcessor.Shutdown()
		log.Info().Msg("task processor is stopped")

		return nil
	})
}

func runGrpcServer(ctx context.Context, waitGroup *errgroup.Group, config utils.Config, store db.Store_interface, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(store, config, taskDistributor)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create the GRPC server")
	}

	gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(gprcLogger)
	pb.RegisterCryptoWebAppServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("gRPC server failed to serve")
			return err
		}

		return nil
	})


	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC server")

		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")

		return nil
	})
}

func runGatewayServer(
    ctx context.Context,
    waitGroup *errgroup.Group,
    config utils.Config,
    store db.Store_interface,
    taskDistributor worker.TaskDistributor,
) {
    _, err := gapi.NewServer(store, config, taskDistributor)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to create the Gateway handler server")
    }

    jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
        MarshalOptions: protojson.MarshalOptions{
            UseProtoNames:   true,
            UseEnumNumbers:  false,
			EmitUnpopulated: true,
        },
        UnmarshalOptions: protojson.UnmarshalOptions{
            DiscardUnknown: true,
        },
    })

    grpcMux := runtime.NewServeMux(jsonOption)

    dialOpts := []grpc.DialOption{
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    }

    err = pb.RegisterCryptoWebAppHandlerFromEndpoint(ctx, grpcMux, config.GRPCServerAddr, dialOpts)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to start grpc-gateway")
    }

    mux := http.NewServeMux()
    mux.Handle("/", grpcMux)

    fs := http.FileServer(http.FS(docsFS))
    mux.Handle("/docs/", http.StripPrefix("/", fs))

	mux.HandleFunc("/oauth/google/login", oauth2.GoogleLoginHandler)
	mux.HandleFunc("/oauth/google/callback", oauth2.GoogleCallbackHandler)

    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        ExposedHeaders:   []string{"Content-Length"},
        AllowCredentials: true,
    })
    handler := c.Handler(gapi.HttpLogger(mux))

    httpServer := &http.Server{
        Handler: handler,
        Addr:    config.HTTPServerAddr,
    }

    waitGroup.Go(func() error {
        log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
        err = httpServer.ListenAndServe()
        if err != nil {
            if errors.Is(err, http.ErrServerClosed) {
                return nil
            }
            log.Error().Err(err).Msg("HTTP gateway server failed to serve")
            return err
        }
        return nil
    })

    waitGroup.Go(func() error {
        <-ctx.Done()
        log.Info().Msg("graceful shutdown HTTP gateway server")

        err := httpServer.Shutdown(context.Background())
        if err != nil {
            log.Error().Err(err).Msg("failed to shutdown HTTP gateway server")
            return err
        }

        log.Info().Msg("HTTP gateway server is stopped")
        return nil
    })
}


func runGinServer(config utils.Config, store db.Store_interface) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create the Gin server")
	}
	
	err = server.Start(config.HTTPServerAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start the Gin server")
	}
}
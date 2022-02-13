package main

import (
	"context"
	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thorgnir-go-study/yp-diploma/internal/api"
	"github.com/thorgnir-go-study/yp-diploma/internal/config"
	orderRepository "github.com/thorgnir-go-study/yp-diploma/internal/infrastructure/repository/order"
	userRepository "github.com/thorgnir-go-study/yp-diploma/internal/infrastructure/repository/user"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/order"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/user"

	shopspring "github.com/jackc/pgtype/ext/shopspring-numeric"

	systemLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		systemLog.Fatalf("error while getting configuration: %v", err)
	}
	configureLogger(*cfg)

	dbpool, err := createDBPool(*cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while creating dbpool")
	}
	userService, err := createUserService(dbpool)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while creating user service")
	}

	orderService, err := createOrdersService(dbpool)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while creating orders service")
	}

	srv := api.NewServer(*cfg, userService, orderService)

	errC, err := run(srv)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't run")
	}

	if err = <-errC; err != nil {
		log.Fatal().Err(err).Msg("Error while running")
	}

}

func createDBPool(cfg config.Config) (*pgxpool.Pool, error) {
	dbconfig, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &shopspring.Numeric{},
			Name:  "numeric",
			OID:   pgtype.NumericOID,
		})
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &pgtypeuuid.UUID{},
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
		})
		return nil
	}
	dbpool, err := pgxpool.ConnectConfig(context.Background(), dbconfig)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

func createUserService(dbpool *pgxpool.Pool) (user.UseCase, error) {
	userRepo := userRepository.NewPostgresUserRepository(dbpool)

	srv := user.NewService(userRepo)
	return srv, nil
}

func createOrdersService(dbpool *pgxpool.Pool) (order.UseCase, error) {
	orderRepo := orderRepository.NewPostgresOrderRepository(dbpool)
	srv := order.NewService(orderRepo)

	return srv, nil
}

func run(srv *http.Server) (<-chan error, error) {
	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		log.Info().Msg("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			stop()
			cancel()
			close(errC)
		}()

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}
		log.Info().Msg("Shutdown completed")

	}()

	go func() {
		log.Info().Msg("Server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()
	return errC, nil
}

func configureLogger(_ config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// в дальнейшем можно добавить в конфиг требуемый уровень логирования, аутпут (файл или еще чего) и т.д.
	// пока пишем в консоль красивенько
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

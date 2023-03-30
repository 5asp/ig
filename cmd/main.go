package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/aszeta/ig/pb"

	"github.com/aszeta/ig/endpoints"
	"github.com/aszeta/ig/service"
	transport "github.com/aszeta/ig/transports"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var db *sql.DB
	{
		var err error
		// Connect to the "ordersdb" database
		db, err = sql.Open("postgres", "postgresql://shijuvar@localhost:26257/ordersdb?sslmode=disable")
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	addservice := service.NewService(db, "", logger)
	addendpoint := endpoints.MakeEndpoints(addservice)
	grpcServer := transport.NewGRPCServer(addendpoint, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log("during", "Listen", "err", err)
		os.Exit(1)
	}

	go func() {
		baseServer := grpc.NewServer()
		pb.RegisterMathServiceServer(baseServer, grpcServer)
		level.Info(logger).Log("msg", "Server started successfully 🚀")
		baseServer.Serve(grpcListener)
	}()

	level.Error(logger).Log("exit", <-errs)
}

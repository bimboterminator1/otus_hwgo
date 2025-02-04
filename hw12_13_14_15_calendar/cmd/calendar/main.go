package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	//nolint:depguard
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/app"
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/config"
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/server/http"
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	logg, err := logger.NewLogger(config.Components.Logging.FilePath,
		config.Components.Logging.Level,
		logger.LogFormat(config.Components.Logging.Type))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't initialize logger: %s", err.Error())
		os.Exit(1)
	}

	var storage storage.Storage
	if config.Components.Storage.Type == "memory" {
		storage = memorystorage.New()
	} else if config.Components.Storage.Type == "postgres" {
		storage, err = sqlstorage.New(config.Components.Storage)
		if err != nil {
			logg.Error(err.Error())
			os.Exit(1)
		}
	}

	calendar := app.New(*logg, storage)

	server := internalhttp.NewServer(logg, *calendar, config.Components.Server)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

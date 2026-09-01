package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"discord-diplomacy/internal/bot"
	"discord-diplomacy/internal/config"
	"discord-diplomacy/internal/features/game"
	"discord-diplomacy/internal/features/orders"
	"discord-diplomacy/internal/sqsconsumer"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.AWSRegion))
	if err != nil {
		logger.Error("load AWS configuration", "error", err)
		os.Exit(1)
	}

	application, err := bot.New(
		cfg,
		logger,
		game.New(),
		orders.New(),
	)
	if err != nil {
		logger.Error("create bot", "error", err)
		os.Exit(1)
	}

	consumer := sqsconsumer.New(sqs.NewFromConfig(awsCfg), logger, application, cfg.SQSQueueURL)

	errs := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		errs <- application.Run(runCtx)
	}()
	go func() {
		defer wg.Done()
		errs <- consumer.Run(runCtx)
	}()
	go func() {
		wg.Wait()
		close(errs)
	}()

	var runErr error
	for err := range errs {
		if err != nil && runErr == nil {
			runErr = err
			cancel()
		}
	}
	if runErr != nil {
		logger.Error("application stopped with an error", "error", runErr)
		os.Exit(1)
	}

	cfg.ActiveGame = application.GetActiveGame()
	if err := cfg.Save(); err != nil {
		logger.Error("failed to save the config", "error", err)
	}
}

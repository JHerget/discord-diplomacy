package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"discord-diplomacy/internal/bot"
	"discord-diplomacy/internal/config"
	"discord-diplomacy/internal/features/feedback"
	"discord-diplomacy/internal/features/ping"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	sharedCfg, err := config.LoadShared(cfg.GuildID)
	if err != nil {
		logger.Error("invalid state file", "error", err)
		os.Exit(1)
	}

	application, err := bot.New(
		cfg,
		logger,
		ping.New(sharedCfg),
		feedback.New(sharedCfg),
	)
	if err != nil {
		logger.Error("create bot", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		logger.Error("bot stopped with an error", "error", err)
		os.Exit(1)
	}
}

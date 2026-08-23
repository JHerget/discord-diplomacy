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
	"discord-diplomacy/internal/features/game"
	"discord-diplomacy/internal/features/ping"
	"discord-diplomacy/internal/utils"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	cctx, err := utils.NewCommandContext(cfg.GuildID)
	if err != nil {
		logger.Error("invalid state file", "error", err)
		os.Exit(1)
	}

	application, err := bot.New(
		cfg,
		logger,
		ping.New(),
		feedback.New(),
		game.New(),
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

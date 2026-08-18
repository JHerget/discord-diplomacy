package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"discord-diplomacy/internal/config"
	"discord-diplomacy/internal/interactions"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session  *discordgo.Session
	guildID  string
	registry *interactions.Registry
	logger   *slog.Logger
}

// New creates the Discord session and registers every supplied feature module.
func New(cfg config.Config, logger *slog.Logger, modules ...interactions.Module) (*Bot, error) {
	registry := interactions.NewRegistry()
	for _, module := range modules {
		if module == nil {
			return nil, errors.New("interaction module cannot be nil")
		}
		if err := module.Register(registry); err != nil {
			return nil, fmt.Errorf("register interaction module: %w", err)
		}
	}

	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("create Discord session: %w", err)
	}
	session.Identify.Intents = discordgo.IntentsGuilds

	bot := &Bot{
		session:  session,
		guildID:  cfg.GuildID,
		registry: registry,
		logger:   logger,
	}
	bot.session.AddHandler(bot.handleInteraction)
	bot.session.AddHandler(bot.handleReady)

	return bot, nil
}

// Run connects to Discord, synchronizes guild commands, and waits for cancellation.
func (b *Bot) Run(ctx context.Context) error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("open Discord session: %w", err)
	}

	applicationID := b.session.State.User.ID
	commands, err := b.session.ApplicationCommandBulkOverwrite(applicationID, b.guildID, b.registry.Commands())
	if err != nil {
		_ = b.session.Close()
		return fmt.Errorf("synchronize guild commands: %w", err)
	}
	b.logger.Info("guild commands synchronized", "guild_id", b.guildID, "count", len(commands))

	<-ctx.Done()
	b.logger.Info("shutting down")
	if err := b.session.Close(); err != nil {
		return fmt.Errorf("close Discord session: %w", err)
	}

	return nil
}

func (b *Bot) handleReady(_ *discordgo.Session, ready *discordgo.Ready) {
	b.logger.Info("connected to Discord", "user", ready.User.String())
}

func (b *Bot) handleInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if err := b.registry.Handle(session, interaction); err != nil {
		if errors.Is(err, interactions.ErrHandlerNotFound) {
			b.logger.Warn("ignored unknown interaction", "type", interaction.Type, "error", err)
			return
		}
		b.logger.Error("interaction handler failed", "type", interaction.Type, "error", err)
	}
}

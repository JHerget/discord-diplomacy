package interactions

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var ErrHandlerNotFound = errors.New("interaction handler not found")

type Handler func(*discordgo.Session, *discordgo.InteractionCreate) error

type Module interface {
	Register(*Registry) error
}

type Registry struct {
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]Handler
	modalHandlers   map[string]Handler
}

func NewRegistry() *Registry {
	return &Registry{
		commandHandlers: make(map[string]Handler),
		modalHandlers:   make(map[string]Handler),
	}
}

func (r *Registry) RegisterCommand(command *discordgo.ApplicationCommand, handler Handler) error {
	if command == nil || command.Name == "" {
		return errors.New("command name is required")
	}
	if handler == nil {
		return fmt.Errorf("command %q has no handler", command.Name)
	}
	if _, exists := r.commandHandlers[command.Name]; exists {
		return fmt.Errorf("command %q is already registered", command.Name)
	}

	r.commands = append(r.commands, command)
	r.commandHandlers[command.Name] = handler
	return nil
}

func (r *Registry) RegisterModal(customID string, handler Handler) error {
	if customID == "" {
		return errors.New("modal custom ID is required")
	}
	if handler == nil {
		return fmt.Errorf("modal %q has no handler", customID)
	}
	if _, exists := r.modalHandlers[customID]; exists {
		return fmt.Errorf("modal %q is already registered", customID)
	}

	r.modalHandlers[customID] = handler
	return nil
}

func (r *Registry) Commands() []*discordgo.ApplicationCommand {
	return append([]*discordgo.ApplicationCommand(nil), r.commands...)
}

func (r *Registry) Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	var (
		key     string
		handler Handler
		exists  bool
	)

	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		key = interaction.ApplicationCommandData().Name
		handler, exists = r.commandHandlers[key]
	case discordgo.InteractionModalSubmit:
		key = interaction.ModalSubmitData().CustomID
		handler, exists = r.modalHandlers[key]
	default:
		return fmt.Errorf("%w: unsupported interaction type %d", ErrHandlerNotFound, interaction.Type)
	}

	if !exists {
		return fmt.Errorf("%w: %q", ErrHandlerNotFound, key)
	}

	return handler(session, interaction)
}

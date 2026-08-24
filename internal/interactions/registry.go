package interactions

import (
	"discord-diplomacy/internal/utils"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var ErrHandlerNotFound = errors.New("interaction handler not found")

type Handler interface {
	Handle(*utils.CommandContext) error
}

type Module interface {
	Handler
	Registration() *discordgo.ApplicationCommand
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

func (r *Registry) AddCommand(command Module) error {
	registration := command.Registration()

	if registration == nil || registration.Name == "" {
		return errors.New("command name is required")
	}
	if _, exists := r.commandHandlers[registration.Name]; exists {
		return fmt.Errorf("command %q is already registered", registration.Name)
	}

	r.commands = append(r.commands, registration)
	r.commandHandlers[registration.Name] = command
	return nil
}

func (r *Registry) AddModal(customID string, handler Handler) error {
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

func (r *Registry) Handle(cctx *utils.CommandContext) error {
	var (
		key     string
		handler Handler
		exists  bool
	)

	switch cctx.Interaction.Type {
	case discordgo.InteractionApplicationCommand:
		key = cctx.Interaction.ApplicationCommandData().Name
		handler, exists = r.commandHandlers[key]
	case discordgo.InteractionModalSubmit:
		key = cctx.Interaction.ModalSubmitData().CustomID
		handler, exists = r.modalHandlers[key]
	default:
		return fmt.Errorf("%w: unsupported interaction type %d", ErrHandlerNotFound, cctx.Interaction.Type)
	}

	if !exists {
		return fmt.Errorf("%w: %q", ErrHandlerNotFound, key)
	}

	return handler.Handle(cctx)
}

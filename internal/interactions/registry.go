package interactions

import (
	"discord-diplomacy/internal/types"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var ErrHandlerNotFound = errors.New("interaction handler not found")

type Handler interface {
	Handle(*types.CommandContext) error
}

type Register interface {
	Register(*Registry) error
}

type Submitter interface {
	Submit(*types.CommandContext) error
}

type Module interface {
	Handler
	Register
}

type Registry struct {
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]Handler
	modalHandlers   map[string]Submitter
}

func NewRegistry() *Registry {
	return &Registry{
		commandHandlers: make(map[string]Handler),
		modalHandlers:   make(map[string]Submitter),
	}
}

func (r *Registry) AddCommand(command *discordgo.ApplicationCommand, handler Handler) error {
	if command == nil || command.Name == "" {
		return errors.New("command name is required")
	}
	if _, exists := r.commandHandlers[command.Name]; exists {
		return fmt.Errorf("command %q is already registered", command.Name)
	}

	r.commands = append(r.commands, command)
	r.commandHandlers[command.Name] = handler
	return nil
}

func (r *Registry) AddModal(customID string, submitter Submitter) error {
	if customID == "" {
		return errors.New("modal custom ID is required")
	}
	if submitter == nil {
		return fmt.Errorf("modal %q has no handler", customID)
	}
	if _, exists := r.modalHandlers[customID]; exists {
		return fmt.Errorf("modal %q is already registered", customID)
	}

	r.modalHandlers[customID] = submitter
	return nil
}

func (r *Registry) Commands() []*discordgo.ApplicationCommand {
	return append([]*discordgo.ApplicationCommand(nil), r.commands...)
}

func (r *Registry) Handle(cctx *types.CommandContext) error {
	switch cctx.Interaction.Type {
	case discordgo.InteractionApplicationCommand:
		key := cctx.Interaction.ApplicationCommandData().Name
		handler, ok := r.commandHandlers[key]
		if ok {
			return handler.Handle(cctx)
		}
	case discordgo.InteractionModalSubmit:
		key := cctx.Interaction.ModalSubmitData().CustomID
		submitter, ok := r.modalHandlers[key]
		if ok {
			return submitter.Submit(cctx)
		}
	}

	return fmt.Errorf("%w: unsupported interaction type %d", ErrHandlerNotFound, cctx.Interaction.Type)
}

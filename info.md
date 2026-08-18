# Project Architecture and How It Works

This project is a small Discord bot written in Go. It uses
[DiscordGo](https://github.com/bwmarrin/discordgo) to connect to Discord's
Gateway, register slash commands, receive interactions, display a modal, and
respond to users.

The project is deliberately split into a small reusable foundation and
self-contained feature packages. The foundation knows how to start the bot,
synchronize commands, and route interactions. It does not know the details of
`/ping`, `/feedback`, or any future feature. Feature packages provide those
details through a common registration interface.

## Project layout

```text
.
├── cmd/
│   └── bot/
│       └── main.go                 Executable entry point and module wiring
├── internal/
│   ├── bot/
│   │   └── bot.go                  Discord session lifecycle
│   ├── config/
│   │   └── config.go               Environment configuration
│   ├── features/
│   │   ├── feedback/
│   │   │   └── feedback.go         /feedback command and modal
│   │   └── ping/
│   │       └── ping.go             /ping command
│   └── interactions/
│       ├── modal.go                Modal input helper
│       └── registry.go             Module API and interaction router
├── .env.example                    Required environment variable names
├── .gitignore                      Local files and build output to ignore
├── go.mod                          Module and dependency declarations
├── go.sum                          Dependency integrity checksums
├── Makefile                        Format, build, and run recipes
├── README.md                       Setup and quick-start guide
└── info.md                         This architectural guide
```

Go's `internal` directory has special meaning. Packages below it can be
imported by code in this module, but cannot be imported by unrelated external
modules. That keeps the bot's implementation private while leaving `cmd/bot`
as the executable entry point.

## End-to-end startup sequence

The bot starts in `cmd/bot/main.go`. Its startup sequence is:

1. Create a standard-library `slog.Logger` that writes text logs to standard
   output.
2. Call `config.Load` to read and validate the required environment variables.
3. Construct each feature module with `ping.New()` and `feedback.New()`.
4. Pass the configuration, logger, and feature modules to `bot.New`.
5. Create a context that is canceled by `SIGINT` or `SIGTERM`.
6. Call `Bot.Run` and keep the process alive until that context is canceled.

The corresponding composition code is:

```go
cfg, err := config.Load()
if err != nil {
    logger.Error("invalid configuration", "error", err)
    os.Exit(1)
}

application, err := bot.New(
    cfg,
    logger,
    ping.New(),
    feedback.New(),
)
if err != nil {
    logger.Error("create bot", "error", err)
    os.Exit(1)
}

ctx, stop := signal.NotifyContext(
    context.Background(),
    os.Interrupt,
    syscall.SIGTERM,
)
defer stop()

if err := application.Run(ctx); err != nil {
    logger.Error("bot stopped with an error", "error", err)
    os.Exit(1)
}
```

If configuration, construction, connection, command synchronization, or
shutdown fails, `main` logs the error and exits with a nonzero status.

### Why `main` contains the module list

`main.go` is the composition root: the one place where concrete features are
selected and assembled. The generic bot and router packages never import the
feature packages. Adding a feature requires adding it to this module list, but
does not require modifying session setup, command synchronization, or routing.

This explicit wiring is intentional. It avoids hidden registration through
package-level `init` functions and makes it immediately visible which features
are enabled in a particular executable.

## Configuration

`internal/config/config.go` defines this configuration value:

```go
type Config struct {
    BotToken string
    GuildID  string
}
```

`Load` reads two environment variables:

- `DISCORD_BOT_TOKEN` authenticates the Gateway and REST API session.
- `DISCORD_GUILD_ID` identifies the development server where commands are
  registered.

Both values are trimmed and must be nonempty. Validation errors are collected
with `errors.Join`, so a startup attempt with both values missing reports both
problems at once.

The important part of the loader is:

```go
cfg := Config{
    BotToken: strings.TrimSpace(os.Getenv("DISCORD_BOT_TOKEN")),
    GuildID:  strings.TrimSpace(os.Getenv("DISCORD_GUILD_ID")),
}

var validationErrors []error
if cfg.BotToken == "" {
    validationErrors = append(
        validationErrors,
        fmt.Errorf("DISCORD_BOT_TOKEN is required"),
    )
}
if cfg.GuildID == "" {
    validationErrors = append(
        validationErrors,
        fmt.Errorf("DISCORD_GUILD_ID is required"),
    )
}

if err := errors.Join(validationErrors...); err != nil {
    return Config{}, err
}
```

The project does not load `.env` files. `.env.example` is only a template of
the required names. Export the variables in the shell, inject them with a
process manager, or use a separate environment-loading tool. The real token
must not be committed; `.env` is ignored by Git.

## Building the bot object

`internal/bot/bot.go` owns Discord connection and lifecycle behavior. `bot.New`
performs construction in two stages.

First, it creates an interaction registry and asks each supplied module to
register itself:

```go
registry := interactions.NewRegistry()
for _, module := range modules {
    module.Register(registry)
}
```

Nil modules, duplicate command names, duplicate modal IDs, missing handlers,
and invalid definitions cause construction to fail before the bot connects to
Discord.

Second, it creates a DiscordGo session with the authentication scheme expected
for bots:

```go
discordgo.New("Bot " + cfg.BotToken)
```

The session requests only `discordgo.IntentsGuilds`. This is a non-privileged
Gateway intent and is enough for this interaction-focused starter. The bot does
not request message content, member lists, presences, or other privileged data.

Two DiscordGo event handlers are attached to the session:

- `handleReady` logs the connected Discord user when the Gateway becomes ready.
- `handleInteraction` forwards interaction events to the registry.

Session creation and handler attachment look like this:

```go
session, err := discordgo.New("Bot " + cfg.BotToken)
if err != nil {
    return nil, fmt.Errorf("create Discord session: %w", err)
}

session.Identify.Intents = discordgo.IntentsGuilds
session.AddHandler(bot.handleInteraction)
session.AddHandler(bot.handleReady)
```

The constructed `Bot` retains the session, guild ID, registry, and logger.

## Connection and command synchronization

`Bot.Run` opens the DiscordGo session. Opening establishes the Gateway
connection and populates session state, including the authenticated bot user.
The bot user's ID is also the Discord application ID used for command
registration.

After connecting, `Run` calls `ApplicationCommandBulkOverwrite` with:

- the authenticated application ID;
- `DISCORD_GUILD_ID`;
- every slash-command definition currently stored in the registry.

This replaces this application's complete command set in the selected guild
with the registry's definitions. It has two useful development properties:

- newly added or modified commands appear quickly because they are guild
  commands;
- commands removed from the code are also removed from that guild.

It also means that manually registered guild commands for the same application
which are not present in this program's registry will be removed at startup.
This project does not create global commands.

The synchronization call is concise because the feature modules already put
their definitions into the registry:

```go
applicationID := b.session.State.User.ID
commands, err := b.session.ApplicationCommandBulkOverwrite(
    applicationID,
    b.guildID,
    b.registry.Commands(),
)
if err != nil {
    _ = b.session.Close()
    return fmt.Errorf("synchronize guild commands: %w", err)
}

b.logger.Info(
    "guild commands synchronized",
    "guild_id", b.guildID,
    "count", len(commands),
)
```

Once synchronization succeeds, `Run` blocks on the cancellation context. A
Ctrl+C sends `SIGINT`; container and service managers commonly send `SIGTERM`.
Either signal cancels the context, logs shutdown, and closes the Discord
session cleanly. Shutdown does not unregister commands, so they remain visible
while the bot is offline.

## The interaction registry

`internal/interactions/registry.go` is the extension boundary between the bot
foundation and feature code.

### Handler

All slash commands and modal submissions use one handler signature:

```go
type Handler func(
    *discordgo.Session,
    *discordgo.InteractionCreate,
) error
```

The session lets the handler reply through Discord's API. The interaction
contains the event type, command or modal data, user, guild, channel, and the
token needed to respond. Returning an error keeps logging policy outside the
feature.

### Module

A feature implements this interface:

```go
type Module interface {
    Register(*Registry) error
}
```

A module may register one command, several commands, modal handlers, or any
combination. This allows related interactions—such as the command that opens a
modal and the handler that processes it—to live in the same package.

### Registry storage

The registry contains three collections:

- an ordered slice of `ApplicationCommand` definitions;
- a command handler map keyed by slash-command name;
- a modal handler map keyed by modal custom ID.

The definition slice is sent to Discord during bulk synchronization. The maps
are used for constant-time dispatch when interactions arrive. `Commands`
returns a new slice so callers cannot append to or reorder the registry's own
slice accidentally. The command definition pointers themselves are shared and
should be treated as immutable after registration.

### Registration validation

`RegisterCommand` requires a non-nil definition, a nonempty command name, and a
non-nil handler. It rejects duplicate names. A successful registration stores
both the public Discord definition and its local handler.

`RegisterModal` requires a nonempty custom ID and a non-nil handler. It rejects
duplicate custom IDs. Modal definitions are not registered with Discord at
startup because a modal is created dynamically as an interaction response;
only its future submission handler needs to be stored.

For example, command registration validates the definition before storing the
definition and handler together:

```go
func (r *Registry) RegisterCommand(
    command *discordgo.ApplicationCommand,
    handler Handler,
) error {
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
```

### Dispatch

Every `InteractionCreate` event reaches `Registry.Handle`:

1. For `InteractionApplicationCommand`, it reads the command name from
   `ApplicationCommandData` and looks it up in the command map.
2. For `InteractionModalSubmit`, it reads the modal custom ID from
   `ModalSubmitData` and looks it up in the modal map.
3. Unsupported interaction types and missing keys return an error wrapping
   `ErrHandlerNotFound`.
4. A matching handler is called and its error is returned unchanged.

The central type switch is what lets one Gateway event handler support both
interaction kinds:

```go
switch interaction.Type {
case discordgo.InteractionApplicationCommand:
    key = interaction.ApplicationCommandData().Name
    handler, exists = r.commandHandlers[key]
case discordgo.InteractionModalSubmit:
    key = interaction.ModalSubmitData().CustomID
    handler, exists = r.modalHandlers[key]
default:
    return fmt.Errorf(
        "%w: unsupported interaction type %d",
        ErrHandlerNotFound,
        interaction.Type,
    )
}

if !exists {
    return fmt.Errorf("%w: %q", ErrHandlerNotFound, key)
}

return handler(session, interaction)
```

`Bot.handleInteraction` treats missing handlers as warnings and ignores them.
Other handler failures are logged as errors. The current foundation does not
automatically send an error response to the user; a handler that fails before
responding may therefore leave Discord showing that the interaction failed.
This avoids accidentally trying to respond twice after a partially completed
handler.

## The `/ping` feature

`internal/features/ping/ping.go` is the smallest complete module.

`ping.New` returns a stateless `Module`. Its `Register` method supplies both
parts of a slash command:

- the Discord-visible definition, named `ping` with a description;
- the local `handle` function that runs when `/ping` is invoked.

The handler calls `Session.InteractionRespond` with
`InteractionResponseChannelMessageWithSource`. That response type creates a
message tied to the triggering interaction. Its content is `Pong!`, and the
`MessageFlagsEphemeral` flag means only the user who invoked the command can
see it.

The complete response is:

```go
return session.InteractionRespond(
    interaction.Interaction,
    &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: "Pong!",
            Flags:   discordgo.MessageFlagsEphemeral,
        },
    },
)
```

## The `/feedback` and modal feature

`internal/features/feedback/feedback.go` demonstrates a multi-step feature.
The same module registers the `/feedback` command and the submission handler
for the modal that command opens.

### Custom IDs

Discord returns component identifiers with later interactions, so stable IDs
connect the modal UI to local handlers:

```text
feedback:submit   identifies the submitted feedback modal
feedback:message  identifies its text input
```

Namespacing the IDs with `feedback:` lowers the chance of collisions as more
features are added. They are package constants so creation and parsing use the
same values.

### Opening the modal

When `/feedback` is invoked, `openModal` responds with
`InteractionResponseModal`. The response data defines:

- the modal custom ID and `Feedback` title;
- one action row;
- one required paragraph-style text input;
- a label, placeholder, and maximum length of 1,000 characters.

Discord renders that response as a modal rather than a channel message. When
the user submits it, Discord sends a new `InteractionModalSubmit` event whose
custom ID is `feedback:submit`. The registry routes that event to
`submitModal`.

The modal response data is built from Discord message components:

```go
Data: &discordgo.InteractionResponseData{
    CustomID: modalCustomID,
    Title:    "Feedback",
    Components: []discordgo.MessageComponent{
        discordgo.ActionsRow{
            Components: []discordgo.MessageComponent{
                discordgo.TextInput{
                    CustomID:    inputCustomID,
                    Label:       "What would you like to share?",
                    Style:       discordgo.TextInputParagraph,
                    Placeholder: "Enter your feedback",
                    Required:    true,
                    MaxLength:   1000,
                },
            },
        },
    },
},
```

### Reading the submitted value

Modal components arrive nested inside action rows. The shared
`interactions.TextInputValue` helper walks each row, inspects its text inputs,
and returns the value whose custom ID matches `feedback:message`. It returns an
error if Discord's submission does not contain the expected input.

The helper accounts for the nested action-row structure:

```go
for _, component := range data.Components {
    row, ok := component.(*discordgo.ActionsRow)
    if !ok {
        continue
    }

    for _, rowComponent := range row.Components {
        input, ok := rowComponent.(*discordgo.TextInput)
        if ok && input.CustomID == customID {
            return input.Value, nil
        }
    }
}
```

`submitModal` also trims the value for validation and rejects a blank message.
Discord already marks the input as required, but the server-side check keeps
the handler's assumption explicit. The actual feedback is currently neither
logged nor persisted. A valid submission receives an ephemeral `Thanks for
your feedback!` response.

The submission handler combines extraction, validation, and acknowledgment:

```go
value, err := interactions.TextInputValue(
    interaction.ModalSubmitData(),
    inputCustomID,
)
if err != nil {
    return err
}
if strings.TrimSpace(value) == "" {
    return fmt.Errorf("feedback message is empty")
}

return session.InteractionRespond(
    interaction.Interaction,
    &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: "Thanks for your feedback!",
            Flags:   discordgo.MessageFlagsEphemeral,
        },
    },
)
```

## Adding another slash command

Create a new package under `internal/features`, define a module, and register a
command definition with its handler. A minimal module looks like this:

```go
package hello

import (
    "discord-diplomacy/internal/interactions"
    "github.com/bwmarrin/discordgo"
)

type Module struct{}

func New() Module { return Module{} }

func (Module) Register(registry *interactions.Registry) error {
    return registry.RegisterCommand(&discordgo.ApplicationCommand{
        Name:        "hello",
        Description: "Say hello",
    }, handle)
}

func handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
    return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{Content: "Hello!"},
    })
}
```

Then import the package in `cmd/bot/main.go` and add `hello.New()` to the
arguments passed to `bot.New`. On the next startup, bulk synchronization adds
`/hello` to the development guild automatically.

Command options, subcommands, permissions, localization, and other Discord
features belong in the `discordgo.ApplicationCommand` definition. The routing
foundation does not need to change because dispatch still uses the top-level
command name.

## Adding another modal

A modal needs two paths:

1. A command or component handler responds with
   `discordgo.InteractionResponseModal` and assigns a unique custom ID.
2. The module registers a modal handler under that exact custom ID with
   `Registry.RegisterModal`.

Keep the opener and submission handler in the same feature package. Use unique,
feature-prefixed custom IDs and stable IDs for every input. For ordinary text
inputs, `interactions.TextInputValue` can be reused. A modal with several
inputs can call the helper once per input ID.

The current router supports application commands and modal submissions. Adding
buttons or select menus would require a third registry map and a new
`InteractionMessageComponent` branch in `Registry.Handle`, after which feature
packages could register component handlers through the same pattern.

## Makefile workflow

The Makefile provides three targets with dependency chaining:

```text
run -> build -> format
```

- `make format` runs `gofmt -w .` to rewrite Go source into canonical format.
- `make build` first formats, creates `bin/`, and compiles `./cmd/bot` to
  `bin/bot`.
- `make run` first builds and then starts `./bin/bot` with the current
  environment.

Because `run` depends on `build` and `build` depends on `format`, every normal
run uses freshly formatted and compiled code. The `bin/` directory is ignored
by Git.

## Go module and dependencies

`go.mod` declares the local module name `discord-diplomacy`, targets Go 1.24,
and directly requires DiscordGo v0.29.0. The additional entries are indirect
dependencies used by DiscordGo for WebSocket communication and supporting
functionality.

`go.sum` records cryptographic checksums for downloaded module versions. It is
not a lockfile in the package-manager sense, but it lets the Go tool verify
that future downloads match the dependency content used by this project. Both
`go.mod` and `go.sum` should be committed.

## Current behavior and boundaries

- Commands are registered only in the configured development guild.
- Commands are synchronized every time the process starts.
- `/ping` and feedback confirmation messages are ephemeral.
- Feedback input is validated but not stored, forwarded, or logged.
- The bot handles slash commands and modal submissions, not prefix commands,
  ordinary messages, buttons, or select menus.
- Handler errors are logged; there is no automatic user-facing error response.
- The project has no tests, database, web server, Docker setup, or CI workflow.
- No privileged Gateway intents are requested.

These boundaries keep the starter small while leaving the interaction registry
and module interface as stable extension points for future bot functionality.

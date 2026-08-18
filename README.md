# Discord Diplomacy Bot

A small Discord bot written in Go with [DiscordGo](https://github.com/bwmarrin/discordgo). It includes an extensible interaction registry, a `/ping` command, and a `/feedback` command that opens a modal.

## Prerequisites

- Go 1.24 or newer
- A Discord application with a bot user
- A development Discord server where you can install the application

## Discord setup

1. Create an application in the [Discord Developer Portal](https://discord.com/developers/applications).
2. On the **Bot** page, create the bot and copy its token.
3. On **OAuth2 > URL Generator**, select the `bot` and `applications.commands` scopes. No privileged gateway intents or bot permissions are required for these examples.
4. Open the generated URL and install the bot in your development server.
5. In Discord, enable Developer Mode under **User Settings > Advanced**, then right-click the development server and copy its ID.

## Configuration

The bot reads configuration directly from environment variables. The project does not automatically load `.env` files.

```sh
export DISCORD_BOT_TOKEN="your-bot-token"
export DISCORD_GUILD_ID="your-development-guild-id"
```

See `.env.example` for the required variable names. Never commit a real bot token.

## Run

```sh
make run
```

Use `make format` to format the Go source or `make build` to create `bin/bot`.

On startup, the bot replaces the development guild's application command definitions with the commands in its registry. This makes command changes appear quickly and removes stale guild commands. Commands remain registered when the process shuts down.

Try `/ping` for a basic ephemeral response. Run `/feedback` to open the example modal; submitting it produces an ephemeral confirmation and does not persist the entered text.

## Add a feature

Create a package under `internal/features` with a type that implements:

```go
type Module interface {
	Register(*interactions.Registry) error
}
```

Within `Register`, call `RegisterCommand` and/or `RegisterModal`. Add the new module to the module list in `cmd/bot/main.go`. The session lifecycle, command synchronization, and interaction dispatcher do not need to change.

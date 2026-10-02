# Discord Diplomacy Bot

A small Discord bot written in Go with [DiscordGo](https://github.com/bwmarrin/discordgo). It includes an extensible interaction registry and commands for managing a Diplomacy game.

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
export AWS_REGION="us-west-2"
export SQS_QUEUE_URL="https://sqs.us-west-2.amazonaws.com/123456789012/your-queue"
```

See `.env.example` for the required variable names. Never commit a real bot token.

The SQS consumer expects plain Discord notification messages:

```json
{
    "type": "discord_message",
    "channel_id": "123456789012345678",
    "content": "Message text to post"
}
```

## Run

```sh
make run
```

Use `make format` to format the Go source or `make build` to create `bin/bot`.

## Run as a systemd service

The supplied unit runs as `jherget` from `/home/jherget/code/discord-diplomacy`.
Create `.env.local` in that directory with the required variables shown above,
using `NAME=value` assignments without `export`. Systemd loads this file; its
AWS values override the defaults in the unit. The service uses `jherget`'s home
directory for saved game state and AWS credentials.

Run setup as `jherget` (do not run `sudo make setup`):

```sh
make setup
```

Setup builds the bot, initializes the saved state file if missing, and uses
`sudo` to install, enable, and restart the service. Running it again applies
updated binaries and unit settings while preserving saved game state.

Check service status and recent logs with:

```sh
systemctl status discord-diplomacy --no-pager -l
sudo journalctl -u discord-diplomacy -n 50 --no-pager
```

On startup, the bot replaces the development guild's application command definitions with the commands in its registry. This makes command changes appear quickly and removes stale guild commands. Commands remain registered when the process shuts down.

## Add a feature

Create a package under `internal/features` with a type that implements:

```go
type Module interface {
	Register(*interactions.Registry) error
}
```

Within `Register`, call `RegisterCommand` and/or `RegisterModal`. Add the new module to the module list in `cmd/bot/main.go`. The session lifecycle, command synchronization, and interaction dispatcher do not need to change.

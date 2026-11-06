package game

import (
	"fmt"
	"strings"

	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/apis/diplomacy/models"
	"discord-diplomacy/internal/interactions"
	"discord-diplomacy/internal/types"

	"github.com/bwmarrin/discordgo"
)

const (
	joinModalCustomID       = "game:join:submit"
	joinPlayerInputCustomID = "game:join:player"
)

var playerCombos = map[int][]string{
	2: {"Austria-Hungary/Germany/Turkey", "England/France/Russia"},
	3: {"England/Germany/Austria-Hungary", "Russia/Italy", "France/Turkey"},
	4: {"England", "Austria-Hungary/France", "Germany/Turkey", "Italy/Russia"},
	5: {"Russia", "England", "Turkey", "France", "Austria-Hungary"},
	6: {"Russia", "England", "Turkey", "France", "Austria-Hungary", "Germany"},
	7: {"Russia", "England", "Turkey", "France", "Austria-Hungary", "Germany", "Italy"},
}

var JoinSubcommand = types.Subcommand{
	Name:        "join",
	Description: "Join the current game.",
	Handler: func(cctx *types.CommandContext) error {
		if cctx.ActiveGame == nil {
			return cctx.BasicEphemeralResponse("There is no active game.")
		}

		userID, ok := cctx.InteractionUserID()
		if !ok {
			return cctx.BasicEphemeralResponse("Could not identify the invoking user.")
		}

		API := diplomacy.NewAPI()
		game, err := API.GetGame(*cctx.ActiveGame)
		if err != nil {
			return cctx.BasicEphemeralResponse(err.Error())
		}

		if game.RandomAssignments {
			players, err := API.CreatePlayers(game.ID, models.CreatePlayersRequest{UserID: userID})
			if err != nil {
				return cctx.BasicEphemeralResponse(err.Error())
			}

			names := make([]string, 0)
			for _, player := range players {
				names = append(names, player.Name)
			}
			playerNames := strings.Join(names, "/")

			if game.ExternalID != nil {
				cctx.BasicChannelResponse(
					*game.ExternalID,
					fmt.Sprintf("%s joined the game as %s.", cctx.InteractionUserName(), playerNames),
				)
			}

			return cctx.BasicEphemeralResponse(fmt.Sprintf("You joined the game as %s.", playerNames))
		}

		options := availablePlayerOptions(game.Players, game.PlayerCount)
		if len(options) == 0 {
			return cctx.BasicEphemeralResponse("The game is full.")
		}

		return cctx.Session.InteractionRespond(cctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: joinModalCustomID,
				Title:    "Join Game",
				Components: []discordgo.MessageComponent{
					discordgo.Label{
						Label: "Player",
						Component: discordgo.SelectMenu{
							CustomID:    joinPlayerInputCustomID,
							Placeholder: "Choose a player",
							MinValues:   new(1),
							MaxValues:   1,
							Required:    new(true),
							Options:     options,
						},
					},
				},
			},
		})
	},
}

type joinModalSubmitter struct{}

func (joinModalSubmitter) Submit(cctx *types.CommandContext) error {
	if cctx.ActiveGame == nil {
		return cctx.BasicEphemeralResponse("There is no active game.")
	}

	userID, ok := cctx.InteractionUserID()
	if !ok {
		return cctx.BasicEphemeralResponse("Could not identify the invoking user.")
	}

	API := diplomacy.NewAPI()
	game, err := API.GetGame(*cctx.ActiveGame)
	if err != nil {
		return cctx.BasicEphemeralResponse(err.Error())
	}
	if game.RandomAssignments {
		return cctx.BasicEphemeralResponse("Player selection is not available for this game.")
	}

	values := interactions.ModalValues(cctx.Interaction.ModalSubmitData())
	playerNames := values[joinPlayerInputCustomID]

	_, err = API.CreatePlayers(game.ID, models.CreatePlayersRequest{
		UserID:      userID,
		PlayerNames: strings.Split(playerNames, "/"),
	})
	if err != nil {
		return cctx.BasicEphemeralResponse(err.Error())
	}

	if game.ExternalID != nil {
		cctx.BasicChannelResponse(
			*game.ExternalID,
			fmt.Sprintf("%s joined the game as %s.", cctx.InteractionUserName(), playerNames),
		)
	}

	return cctx.BasicEphemeralResponse(fmt.Sprintf("You joined the game as %s.", playerNames))
}

func availablePlayerOptions(players []models.Player, playerCount int) []discordgo.SelectMenuOption {
	lookup := make(map[string]models.Player)
	for _, player := range players {
		lookup[player.Name] = player
	}

	options := make([]discordgo.SelectMenuOption, 0)
	combo := playerCombos[playerCount]

	for _, c := range combo {
		names := strings.Split(c, "/")
		skipCombo := false

		for _, name := range names {
			if lookup[name].UserID != nil {
				skipCombo = true
				break
			}
		}

		if skipCombo {
			continue
		}

		options = append(options, discordgo.SelectMenuOption{
			Label: c,
			Value: c,
		})
	}

	return options
}

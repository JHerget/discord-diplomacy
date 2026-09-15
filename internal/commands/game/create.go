package game

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"discord-diplomacy/internal/apis/diplomacy"
	"discord-diplomacy/internal/apis/diplomacy/models"
	"discord-diplomacy/internal/interactions"
	"discord-diplomacy/internal/types"
	"discord-diplomacy/internal/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	createModalCustomID    = "game:create:submit"
	startDateInputCustomID = "game:create:start-date"
	startHourInputCustomID = "game:create:start-hour"
	timezoneInputCustomID  = "game:create:timezone"
	daysPerTurnCustomID    = "game:create:days-per-turn"
)

var timezoneOffsets = map[string]int{
	"Eastern":  -5,
	"Central":  -6,
	"Mountain": -7,
	"Pacific":  -8,
	"Alaska":   -9,
	"Hawaii":   -10,
}
var daysPerTurnValues = []int{1, 2, 3, 7, 14, 21, 30}
var dateFormat = "01/02/2006"

var CreateSubcommand = types.Subcommand{
	Name:        "create",
	Description: "Create a new game.",
	Handler: func(cctx *types.CommandContext) error {
		return cctx.Session.InteractionRespond(cctx.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: createModalCustomID,
				Title:    "Create Game",
				Components: []discordgo.MessageComponent{
					discordgo.Label{
						Label: "Start Date",
						Component: discordgo.TextInput{
							CustomID:    startDateInputCustomID,
							Style:       discordgo.TextInputShort,
							Placeholder: "mm/dd/yyyy",
							Required:    utils.BoolPtr(true),
							MinLength:   10,
							MaxLength:   10,
						},
					},
					discordgo.Label{
						Label: "Start Hour",
						Component: discordgo.SelectMenu{
							CustomID:    startHourInputCustomID,
							Placeholder: "Choose start hour",
							MinValues:   utils.IntPtr(1),
							MaxValues:   1,
							Required:    utils.BoolPtr(true),
							Options:     startHourOptions(),
						},
					},
					discordgo.Label{
						Label: "Timezone",
						Component: discordgo.SelectMenu{
							CustomID:    timezoneInputCustomID,
							Placeholder: "Choose timezone",
							MinValues:   utils.IntPtr(1),
							MaxValues:   1,
							Required:    utils.BoolPtr(true),
							Options:     timezoneOptions(),
						},
					},
					discordgo.Label{
						Label: "Days Per Turn",
						Component: discordgo.SelectMenu{
							CustomID:    daysPerTurnCustomID,
							Placeholder: "Choose days per turn",
							MinValues:   utils.IntPtr(1),
							MaxValues:   1,
							Required:    utils.BoolPtr(true),
							Options:     daysPerTurnOptions(),
						},
					},
				},
			},
		})
	},
}

func (c Command) Submit(cctx *types.CommandContext) error {
	API := diplomacy.NewAPI()

	if cctx.ActiveGame != nil {
		g, err := API.GetGame(*cctx.ActiveGame)
		if err == nil && g.InProgress {
			return cctx.BasicEphemeralResponse("There is already a game in progress.")
		}
	}

	values := interactions.ModalValues(cctx.Interaction.ModalSubmitData())

	startDate, err := parseStartDate(values[startDateInputCustomID])
	if err != nil {
		return cctx.BasicEphemeralResponse(err.Error())
	}

	startHour, err := parseStartHour(values[startHourInputCustomID])
	if err != nil {
		return cctx.BasicEphemeralResponse(err.Error())
	}

	timezoneName := values[timezoneInputCustomID]
	timezoneOffset, ok := timezoneOffsets[timezoneName]
	if !ok {
		return cctx.BasicEphemeralResponse("Select a valid timezone.")
	}

	location := time.FixedZone(timezoneName, timezoneOffset*60*60)
	startTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startHour, 0, 0, 0, location)
	if !startTime.After(time.Now()) {
		return cctx.BasicEphemeralResponse("Start date and time must be in the future.")
	}

	daysPerTurn, err := parseDaysPerTurn(values[daysPerTurnCustomID])
	if err != nil {
		return cctx.BasicEphemeralResponse(err.Error())
	}

	g, err := API.CreateGame(models.CreateGameRequest{
		ExternalID:    &cctx.GuildID,
		MapID:         "6956498133c5739468982b62",
		DaysPerTurn:   daysPerTurn,
		TurnStartHour: startHour,
		Timezone:      timezoneOffset,
		StartDate:     int(startTime.UTC().Unix()),
	})
	if err != nil {
		return cctx.BasicEphemeralResponse(err.Error())
	}

	cctx.SetActiveGame(g.ID)

	return cctx.BasicResponse(fmt.Sprintf("New game starting on %s!", startDate.Format(dateFormat)))
}

func parseStartDate(value string) (time.Time, error) {
	startDate, err := time.Parse(dateFormat, value)
	if err != nil || startDate.Format(dateFormat) != value {
		return time.Time{}, fmt.Errorf("Enter a valid start date in MM/DD/YYYY format.")
	}

	return startDate, nil
}

func parseStartHour(value string) (int, error) {
	startHour, err := strconv.Atoi(value)
	if err != nil || startHour < 0 || startHour > 23 {
		return 0, fmt.Errorf("Select a valid start hour.")
	}

	return startHour, nil
}

func parseDaysPerTurn(value string) (int, error) {
	daysPerTurn, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("Select a valid days per turn.")
	}

	if slices.Contains(daysPerTurnValues, daysPerTurn) {
		return daysPerTurn, nil
	}

	return 0, fmt.Errorf("Select a valid days per turn.")
}

func startHourOptions() []discordgo.SelectMenuOption {
	options := make([]discordgo.SelectMenuOption, 0, 24)
	for hour := 0; hour <= 23; hour++ {
		value := strconv.Itoa(hour)
		options = append(options, discordgo.SelectMenuOption{
			Label: value,
			Value: value,
		})
	}

	return options
}

func timezoneOptions() []discordgo.SelectMenuOption {
	return []discordgo.SelectMenuOption{
		{Label: "Eastern", Value: "Eastern", Description: "UTC-5"},
		{Label: "Central", Value: "Central", Description: "UTC-6"},
		{Label: "Mountain", Value: "Mountain", Description: "UTC-7"},
		{Label: "Pacific", Value: "Pacific", Description: "UTC-8"},
		{Label: "Alaska", Value: "Alaska", Description: "UTC-9"},
		{Label: "Hawaii", Value: "Hawaii", Description: "UTC-10"},
	}
}

func daysPerTurnOptions() []discordgo.SelectMenuOption {
	options := make([]discordgo.SelectMenuOption, 0, len(daysPerTurnValues))
	for _, day := range daysPerTurnValues {
		value := strconv.Itoa(day)
		options = append(options, discordgo.SelectMenuOption{
			Label: value,
			Value: value,
		})
	}

	return options
}

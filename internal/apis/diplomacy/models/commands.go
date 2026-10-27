package models

import (
	"discord-diplomacy/internal/utils"
	"fmt"
	"strings"
)

type CommandValidationResult struct {
	PlayerName string `json:"playerName"`
	Value      string `json:"value"`
	Valid      bool   `json:"valid"`
	Error      string `json:"error"`
}
type CommandValidationResults []CommandValidationResult

func (c CommandValidationResults) PrettyPrint() string {
	playerCommands := make(map[string][]CommandValidationResult)

	for _, result := range c {
		playerCommands[result.PlayerName] = append(playerCommands[result.PlayerName], result)
	}

	var result strings.Builder
	for player, commands := range playerCommands {
		fmt.Fprintf(&result, "### %s\n", player)

		for _, c := range commands {
			valid := "❌"
			if c.Valid {
				valid = "✅"
			}

			err := ""
			if c.Error != "" {
				err = fmt.Sprintf(" (%s)", c.Error)
			}

			fmt.Fprintf(&result, "%s %s%s\n", valid, utils.ToTitle(c.Value), err)
		}
	}

	return result.String()
}

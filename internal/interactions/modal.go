package interactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func TextInputValue(data discordgo.ModalSubmitInteractionData, customID string) (string, error) {
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

	return "", fmt.Errorf("text input %q was not found in modal submission", customID)
}

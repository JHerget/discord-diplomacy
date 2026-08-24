package interactions

import "github.com/bwmarrin/discordgo"

func TextInputValues(data discordgo.ModalSubmitInteractionData) map[string]string {
	values := make(map[string]string)

	for _, component := range data.Components {
		row, ok := component.(*discordgo.ActionsRow)
		if !ok {
			continue
		}

		for _, rowComponent := range row.Components {
			input, ok := rowComponent.(*discordgo.TextInput)
			if ok {
				values[input.CustomID] = input.Value
			}
		}
	}

	return values
}

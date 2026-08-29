package interactions

import "github.com/bwmarrin/discordgo"

func TextInputValues(data discordgo.ModalSubmitInteractionData) map[string]string {
	return ModalValues(data)
}

func ModalValues(data discordgo.ModalSubmitInteractionData) map[string]string {
	values := make(map[string]string)

	for _, component := range data.Components {
		collectModalValues(values, component)
	}

	return values
}

func collectModalValues(values map[string]string, component discordgo.MessageComponent) {
	switch c := component.(type) {
	case *discordgo.ActionsRow:
		for _, rowComponent := range c.Components {
			collectModalValues(values, rowComponent)
		}
	case *discordgo.Label:
		collectModalValues(values, c.Component)
	case *discordgo.TextInput:
		values[c.CustomID] = c.Value
	case *discordgo.SelectMenu:
		if len(c.Values) > 0 {
			values[c.CustomID] = c.Values[0]
		}
	}
}

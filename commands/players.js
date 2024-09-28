import { SlashCommandBuilder } from "discord.js";
import StateManager from "../state/state-manager.js";

const players = {
  data: new SlashCommandBuilder()
    .setName("players")
    .setDescription("List all players in the game."),
  async execute(interaction) {
    const allPlayers = StateManager.getAllPlayers();
    let message = "";

    for (const player of allPlayers) {
      message += `* ${player.globalName} is playing as ${player.power} (${player.active ? "active" : "inactive"})\n`;
    }

    if (message.length > 0) {
      await interaction.reply({
        content: message,
        ephemeral: true,
      });
    } else {
      await interaction.reply({
        content: "No players are registered in this game.",
        ephemeral: true,
      });
    }
  },
};

export default players;

import { SlashCommandBuilder } from "discord.js";
import StateManager from "../managers/state/state-manager.js";

const deregister = {
  data: new SlashCommandBuilder()
    .setName("deregister")
    .setDescription("Remove yourself from the game."),
  async execute(interaction) {
    if (!StateManager.isRegistered(interaction.user.id)) {
      await interaction.reply({
        content: "Silly goose, you're not even registered in this game!",
        ephemeral: true,
      });
      return;
    }

    StateManager.deregisterPlayer(interaction.user.id);
    await interaction.reply({
      content: "You have successfully been removed from the game.",
      ephemeral: true,
    });
  },
};

export default deregister;

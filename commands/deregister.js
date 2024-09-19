import { SlashCommandBuilder } from "discord.js";
import StateManager from "../state/state-manager.js";

const deregister = {
  data: new SlashCommandBuilder()
    .setName("deregister")
    .setDescription("Remove yourself from the game."),
  async execute(interaction) {
    if (!StateManager.isRegistered(interaction.user.id)) {
      await interaction.reply(
        "Silly goose, you're not even registered in this game!",
      );
      return;
    }

    StateManager.deregisterPlayer(interaction.user.id);
    await interaction.reply(
      "You have successfully been removed from the game.",
    );
  },
};

export default deregister;

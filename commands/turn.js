import { SlashCommandBuilder } from "discord.js";
import StateManager from "../state/state-manager.js";

const turn = {
  data: new SlashCommandBuilder()
    .setName("turn")
    .setDescription("Gives the details about the current turn."),
  async execute(interaction) {
    const turn = StateManager.getCurrentTurn();

    await interaction.reply({
      content: `## ${turn.name}\nPhase: ${turn.phase}\nOrders received: ${turn.orders.join(", ")}`,
      ephemeral: true,
    });
  },
};

export default turn;

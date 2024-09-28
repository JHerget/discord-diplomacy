import { SlashCommandBuilder } from "discord.js";
import StateManager from "../state/state-manager.js";

const turn = {
  data: new SlashCommandBuilder()
    .setName("turn")
    .setDescription("Gives the details about the current turn."),
  async execute(interaction) {
    const turn = StateManager.getCurrentTurn();
    const people = StateManager.getOrders().map((order) => order.player.power);

    await interaction.reply({
      content: `## ${turn.name}\nPhase: ${turn.phase}\nOrders received: ${people.join(", ")}`,
      ephemeral: true,
    });
  },
};

export default turn;

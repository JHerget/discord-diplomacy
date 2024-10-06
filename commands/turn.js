import { SlashCommandBuilder } from "discord.js";
import StateManager from "../managers/state/state-manager.js";

const turn = {
  data: new SlashCommandBuilder()
    .setName("turn")
    .setDescription("Gives the details about the current turn."),
  async execute(interaction) {
    const turn = StateManager.getCurrentTurn();
    const orders = StateManager.getOrders();
    const people = orders.map((order) => order.player.power);
    const playerOrders = orders
      .filter(
        (order) =>
          order.phase == turn.phase &&
          order.player.userId == interaction.user.id,
      )
      .map((order) => order.orders.join("\n"));

    await interaction.reply({
      content: `## ${turn.name}\nPhase: ${turn.phase}\nOrders received: ${people.join(", ")}\n\nYour orders:\n${playerOrders.join("\n")}`,
      ephemeral: true,
    });
  },
};

export default turn;

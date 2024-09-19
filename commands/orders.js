import { SlashCommandBuilder } from "discord.js";
import StateManager from "../state/state-manager.js";
import ordersModal from "../modals/orders-modal.js";

const orders = {
  data: new SlashCommandBuilder()
    .setName("orders")
    .setDescription("Produces a form to submit the orders for your units."),
  async execute(interaction) {
    if (!StateManager.isRegistered(interaction.user.id)) {
      await interaction.reply("You are not a registered player!");
      return;
    }

    await interaction.showModal(ordersModal.data);
  },
};

export default orders;

import {
  ActionRowBuilder,
  ModalBuilder,
  TextInputBuilder,
  TextInputStyle,
} from "discord.js";
import StateManager from "../state/state-manager.js";

const orders = new TextInputBuilder()
  .setCustomId("orders")
  .setLabel("What are your orders?")
  .setStyle(TextInputStyle.Paragraph)
  .setRequired(true)
  .setPlaceholder("F Nap-Tys\nA Rom-Tus");

const ordersActionRow = new ActionRowBuilder().addComponents(orders);
const modal = new ModalBuilder()
  .setCustomId("orders-modal")
  .setTitle("Orders")
  .addComponents(ordersActionRow);

const ordersModal = {
  name: "orders-modal",
  data: modal,
  async execute(interaction) {
    const powerName = StateManager.getPowerName(interaction.user.id);
    const turn = StateManager.getCurrentTurn();
    const orders = interaction.fields.getTextInputValue("orders").split("\n");

    if (StateManager.addOrders(interaction.user.id, orders)) {
      await interaction.reply(
        `${powerName}'s orders have been submitted for the '${turn.phase}' phase of ${turn.name}!`,
      );
    } else {
      await interaction.reply(
        "One or more of your orders were formatted incorrectly. Please try again...or don't and just let Italy win.",
      );
    }
  },
};

export default ordersModal;

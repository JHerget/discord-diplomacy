import {
  ActionRowBuilder,
  ModalBuilder,
  TextInputBuilder,
  TextInputStyle,
} from "discord.js";
import StateManager from "../managers/state/state-manager.js";
import { OrderManager } from "../managers/order/order-manager.js";

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
    const turn = StateManager.getCurrentTurn();
    const rawOrders = interaction.fields
      .getTextInputValue("orders")
      .trim()
      .split("\n");
    let formattedOrders = [];

    for (const order of rawOrders) {
      if (OrderManager.isOrderFormatValid(order)) {
        formattedOrders.push(OrderManager.parseOrder(order));
      } else {
        await interaction.reply({
          content: `This order is invalid for the current phase or is formatted incorrectly: \`${order}\`. Use the \`/format\` command to review the correct formats for all order types.`,
          ephemeral: true,
        });

        return;
      }
    }

    StateManager.addOrders(interaction.user.id, formattedOrders);

    await interaction.reply({
      content: `Your orders have been submitted for the '${turn.phase}' phase of ${turn.name}!`,
      ephemeral: true,
    });
  },
};

export default ordersModal;

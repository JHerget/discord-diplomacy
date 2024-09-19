import { Client, Collection, GatewayIntentBits } from "discord.js";
import config from "./config.json" with { type: "json" };
import orders from "./commands/orders.js";
import register from "./commands/register.js";
import deregister from "./commands/deregister.js";
import turn from "./commands/turn.js";
import ordersModal from "./modals/orders-modal.js";
import ready from "./events/ready.js";
import commandInteraction from "./events/command-interaction.js";
import modalInteraction from "./events/modal-interaction.js";
import cron from "node-cron";
import moment from "moment-timezone";
import StateManager from "./state/state-manager.js";

const client = new Client({
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.MessageContent,
  ],
});

client.commands = new Collection();
client.commands.set(orders.data.name, orders);
client.commands.set(register.data.name, register);
client.commands.set(deregister.data.name, deregister);
client.commands.set(turn.data.name, turn);

client.modals = new Collection();
client.modals.set(ordersModal.name, ordersModal);

client.once(ready.name, (...args) => ready.execute(...args));
client.on(commandInteraction.name, (...args) =>
  commandInteraction.execute(...args),
);
client.on(modalInteraction.name, (...args) =>
  modalInteraction.execute(...args),
);

client.login(config.token);

cron.schedule("* * * * *", () => {
  const currentTime = moment().tz("America/Denver").format("HH:mm");

  // DIPLOMATIC PHASE END
  if (currentTime === "23:30") {
    const orders = StateManager.getOrders();
    const channel = client.channels.cache.find(
      (ch) => ch.name === "orders" && ch.isTextBased(),
    );

    if (channel) {
      let message = "### Unit Orders: \n\n";

      for (let order of orders) {
        message += `${order.powerName}\n`;
        message += `${order.ordersDate}\n`;
        message += `${order.orders.join("\n")}\n\n`;
      }

      channel.send(message);
    }

    stateManager.nextPhase();
  }

  // RETREAT PHASE END
  if (currentTime === "11:00") {
    const orders = StateManager.getOrders();
    const channel = client.channels.cache.find(
      (ch) => ch.name === "orders" && ch.isTextBased(),
    );

    if (channel) {
      let message = "### Retreat Orders: \n\n";

      for (let order of orders) {
        message += `${order.powerName}\n`;
        message += `${order.ordersDate}\n`;
        message += `${order.orders.join("\n")}\n\n`;
      }

      channel.send(message);
    }

    stateManager.nextPhase();
  }

  // REINFORCE PHASE END
  if (currentTime === "12:00") {
    const orders = StateManager.getOrders();
    const channel = client.channels.cache.find(
      (ch) => ch.name === "orders" && ch.isTextBased(),
    );

    if (channel) {
      let message = "### Reinforcement Orders: \n\n";

      for (let order of orders) {
        message += `${order.powerName}\n`;
        message += `${order.ordersDate}\n`;
        message += `${order.orders.join("\n")}\n\n`;
      }

      channel.send(message);
    }

    stateManager.nextPhase();
  }
});

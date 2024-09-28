import { Client, Collection, GatewayIntentBits } from "discord.js";
import config from "./config.json" with { type: "json" };
import orders from "./commands/orders.js";
import register from "./commands/register.js";
import deregister from "./commands/deregister.js";
import turn from "./commands/turn.js";
import players from "./commands/players.js";
import ordersModal from "./modals/orders-modal.js";
import onReady from "./events/on-ready.js";
import commandInteraction from "./events/command-interaction.js";
import modalInteraction from "./events/modal-interaction.js";

const client = new Client({
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMembers,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.MessageContent,
  ],
});

client.commands = new Collection();
client.commands.set(orders.data.name, orders);
client.commands.set(register.data.name, register);
client.commands.set(deregister.data.name, deregister);
client.commands.set(turn.data.name, turn);
client.commands.set(players.data.name, players);

client.modals = new Collection();
client.modals.set(ordersModal.name, ordersModal);

client.once(onReady.name, (...args) => onReady.execute(...args));
client.on(commandInteraction.name, (...args) =>
  commandInteraction.execute(...args),
);
client.on(modalInteraction.name, (...args) =>
  modalInteraction.execute(...args),
);

client.login(config.token);

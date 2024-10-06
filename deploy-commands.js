import { REST, Routes } from "discord.js";
import config from "./config.json" with { type: "json" };
import orders from "./commands/orders.js";
import register from "./commands/register.js";
import deregister from "./commands/deregister.js";
import turn from "./commands/turn.js";
import players from "./commands/players.js";
import format from "./commands/format.js";

const commands = [
  orders.data.toJSON(),
  register.data.toJSON(),
  deregister.data.toJSON(),
  turn.data.toJSON(),
  players.data.toJSON(),
  format.data.toJSON(),
];
const rest = new REST().setToken(config.token);

(async () => {
  try {
    console.log(
      `Started refreshing ${commands.length} application (/) commands.`,
    );

    const data = await rest.put(Routes.applicationCommands(config.clientId), {
      body: commands,
    });

    console.log(
      `Successfully reloaded ${data.length} application (/) commands.`,
    );
  } catch (error) {
    console.error(error);
  }
})();

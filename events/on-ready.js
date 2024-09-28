import { Events } from "discord.js";
import ServerManager from "../server/server-manager.js";
import phaseJob from "../server/cron.js";

const onReady = {
  name: Events.ClientReady,
  once: true,
  execute(client) {
    ServerManager.setClient(client);
    ServerManager.setGameStateChannel(
      client.channels.cache.find((ch) => ch.name === "game-state"),
    );
    ServerManager.setOrdersChannel(
      client.channels.cache.find((ch) => ch.name === "orders"),
    );
    ServerManager.setGreatPowerRole(
      client.guilds.cache.first().roles.cache.find(role => role.name == "Great Power")
    );

    phaseJob.start();

    console.log(`Ready! Logged in as ${client.user.tag}`);
  },
};

export default onReady;

import { SlashCommandBuilder } from "discord.js";
import StateManager from "../state/state-manager.js";
import ServerManager from "../server/server-manager.js";

const register = {
  data: new SlashCommandBuilder()
    .setName("register")
    .setDescription("Register yourself as a player in the game."),
  async execute(interaction) {
    const availablePowers = StateManager.getAvailablePowers();
    const randomIndex = Math.floor(Math.random() * availablePowers.length);

    if (StateManager.inProgress()) {
      await interaction.reply({
        content: "The game has already started.",
        ephemeral: true,
      });
      return;
    }

    if (availablePowers.length <= 0) {
      await interaction.reply({
        content: "No more players can be added to the game.",
        ephemeral: true,
      });
      return;
    }

    if (StateManager.isRegistered(interaction.user.id)) {
      await interaction.reply({
        content: "You are already registered. You trynna cheat??",
        ephemeral: true,
      });
      return;
    }

    const chosenPower = availablePowers[randomIndex];

    if (StateManager.registerPlayer(chosenPower, interaction.user)) {
      await interaction.reply(
        `${interaction.user.globalName} is registered as ${chosenPower}!`,
      );

      try {
        await interaction.member.roles.add(ServerManager.greatPowerRole);
      } catch (error) {
        console.log(error);
        await interaction.followUp({
          content: `I was unable to grant you the role of 'Great Power'. You'll have to ask the game admin to help you 😔.`,
          ephemeral: true,
        });
      }

      try {
        await interaction.member.setNickname(chosenPower);
      } catch {
        await interaction.followUp({
          content: `I was unable to change your nickname to ${chosenPower} so you'll have to do it yourself 😔.`,
          ephemeral: true,
        });
      }
    }
  },
};

export default register;

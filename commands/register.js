import { SlashCommandBuilder } from "discord.js";
import StateManager from "../state/state-manager.js";

const register = {
  data: new SlashCommandBuilder()
    .setName("register")
    .setDescription("Register yourself as a player in the game.")
    .addStringOption((option) =>
      option
        .setName("great-power-name")
        .setDescription(
          "This is the name of the great power you wish to register as.",
        )
        .setRequired(true)
        .addChoices(
          { name: "Austria-Hungary", value: "Austria-Hungary" },
          { name: "England", value: "England" },
          { name: "France", value: "France" },
          { name: "Germany", value: "Germany" },
          { name: "Italy", value: "Italy" },
          { name: "Russia", value: "Russia" },
          { name: "Turkey", value: "Turkey" },
        ),
    ),
  async execute(interaction) {
    const chosenPower = interaction.options.getString("great-power-name");
    const availablePowers = StateManager.getAvailablePowers();

    if (StateManager.inProgress()) {
      await interaction.reply({
        content: "The game has already started.",
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

    if (StateManager.registerPlayer(chosenPower, interaction.user.id)) {
      await interaction.reply({
        content: `You are registered as ${chosenPower}!`,
        ephemeral: true,
      });
    } else {
      await interaction.reply({
        content: `Someone is already registered as ${chosenPower} *womp womp*\n\nAvailable great powers: ${availablePowers.join(", ")}`,
        ephemeral: true,
      });
    }
  },
};

export default register;

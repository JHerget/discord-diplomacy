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
      await interaction.reply("The game has already started.");
      return;
    }

    if (StateManager.isRegistered(interaction.user.id)) {
      await interaction.reply("You are already registered. You trynna cheat??");
      return;
    }

    if (StateManager.registerPlayer(chosenPower, interaction.user.id)) {
      await interaction.reply(`You are registered as ${chosenPower}!`);
    } else {
      await interaction.reply(
        `Someone is already registered as ${chosenPower} *womp womp*\n\nAvailable great powers: ${availablePowers.join(", ")}`,
      );
    }
  },
};

export default register;

import { Events } from "discord.js";

const commandInteraction = {
  name: Events.InteractionCreate,
  async execute(interaction) {
    if (!interaction.isChatInputCommand()) return;

    const command = interaction.client.commands.get(interaction.commandName);

    if (!command) {
      console.error(`No command matching ${interaction.commandName}`);
    }

    try {
      await command.execute(interaction);
    } catch (error) {
      console.error(error);
    }
  },
};

export default commandInteraction;

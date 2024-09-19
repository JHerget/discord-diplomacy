import { Events } from "discord.js";

const modalInteraction = {
  name: Events.InteractionCreate,
  async execute(interaction) {
    if (!interaction.isModalSubmit()) return;

    const modal = interaction.client.modals.get(interaction.customId);

    if (!modal) {
      console.error(`No modal matching ${interaction.customId}`);
    }

    try {
      await modal.execute(interaction);
    } catch (error) {
      console.error(error);
    }
  },
};

export default modalInteraction;

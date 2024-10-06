import { SlashCommandBuilder } from "discord.js";

const format = {
  data: new SlashCommandBuilder()
    .setName("format")
    .setDescription("Lists the format for each type of order."),
  async execute(interaction) {
    await interaction.reply({
      content: `### Hold\nFormat: \`<unit-type> <location>-<hold or h>\`\nExample: \`A Par-Hold\` or \`A Par-H\`\n### Move\nFormat: \`<unit-type> <location>-<destination>\`\nExample: \`A Par-Bur\`\n### Support\nFormat: \`<unit-type> <location> S <unit-type> <location>-<destination or hold or h>\`\nExample: \`A Mar S A Par-Bur\` or \`A Bur S A Par-Hold\` or \`A Bur S A Par-H\`\n### Convoy\nFormat: \`F <location> C A <location>-<destination>\`\nExample: \`F GoL C A Mar-Spa\`\n### Retreat\nFormat: \`<unit-type> <location>-<destination>\`\nExample: \`A Par-Bur\`\n### Disband\nFormat: \`D <location>\`\nExample: \`D Par\`\n### Reinforce\nFormat: \`<unit-type> <location>\`\nExample: \`A Par\``,
      ephemeral: true,
    });
  },
};

export default format;

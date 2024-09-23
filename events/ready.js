import { Events } from "discord.js";
import cron from "node-cron";
import moment from "moment-timezone";
import StateManager from "../state/state-manager.js";

const onReady = {
  name: Events.ClientReady,
  once: true,
  execute(client) {
    const boardStateChannel = client.channels.cache.find(
      (ch) => ch.name === "board-state",
    );
    const ordersChannel = client.channels.cache.find(
      (ch) => ch.name === "orders",
    );

    cron.schedule("* * * * *", () => {
      const turn = StateManager.getCurrentTurn();
      const orders = StateManager.getOrders();

      if (ordersChannel) {
        let message = `## ${turn.name}\nUnit Orders: \n\n`;

        for (let order of orders) {
          message += `${order.player.power}\n`;
          message += `${order.orders.join("\n")}\n\n`;
        }

        ordersChannel.send(message);
      }

      StateManager.nextPhase();
    });

    console.log(`Ready! Logged in as ${client.user.tag}`);
  },
};

export default onReady;

// cron.schedule("* * * * *", () => {
//   const currentTime = moment().tz("America/Denver").format("HH:mm");

//   // DIPLOMATIC PHASE END
//   if (currentTime === "23:30") {
//     const orders = StateManager.getOrders();
//     const channel = client.channels.cache.find(
//       (ch) => ch.name === "orders" && ch.isTextBased(),
//     );

//     if (channel) {
//       let message = "### Unit Orders: \n\n";

//       for (let order of orders) {
//         message += `${order.powerName}\n`;
//         message += `${order.ordersDate}\n`;
//         message += `${order.orders.join("\n")}\n\n`;
//       }

//       channel.send(message);
//     }

//     stateManager.nextPhase();
//   }

//   // RETREAT PHASE END
//   if (currentTime === "11:00") {
//     const orders = StateManager.getOrders();
//     const channel = client.channels.cache.find(
//       (ch) => ch.name === "orders" && ch.isTextBased(),
//     );

//     if (channel) {
//       let message = "### Retreat Orders: \n\n";

//       for (let order of orders) {
//         message += `${order.powerName}\n`;
//         message += `${order.ordersDate}\n`;
//         message += `${order.orders.join("\n")}\n\n`;
//       }

//       channel.send(message);
//     }

//     stateManager.nextPhase();
//   }

//   // REINFORCE PHASE END
//   if (currentTime === "12:00") {
//     const orders = StateManager.getOrders();
//     const channel = client.channels.cache.find(
//       (ch) => ch.name === "orders" && ch.isTextBased(),
//     );

//     if (channel) {
//       let message = "### Reinforcement Orders: \n\n";

//       for (let order of orders) {
//         message += `${order.powerName}\n`;
//         message += `${order.ordersDate}\n`;
//         message += `${order.orders.join("\n")}\n\n`;
//       }

//       channel.send(message);
//     }

//     stateManager.nextPhase();
//   }
// });

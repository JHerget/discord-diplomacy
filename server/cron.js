import cron from "node-cron";
import moment from "moment-timezone";
import StateManager from "../state/state-manager.js";
import { Phase } from "../state/state-models.js";
import ServerManager from "./server-manager.js";

const phaseJob = cron.schedule("0 * * * *", () => {
  const currentTime = moment().tz("America/Denver").format("HH:mm");
  const phase = StateManager.getCurrentPhase();

  if (currentTime === "12:00" && phase === Phase.Diplomatic) {
    StateManager.incrementDaysOnTurn();

    const daysOnTurn = StateManager.getCurrentTurn().daysOnTurn;
    const daysPerTurn = StateManager.getDaysPerTurn();
    const allOrdersSubmitted = StateManager.allOrdersSubmitted();

    if (daysOnTurn == daysPerTurn || allOrdersSubmitted) {
      ServerManager.ordersChannel.send(StateManager.getOrdersMessage());
      StateManager.nextPhase();
      ServerManager.gameStateChannel.send(
        `We are now on the '${StateManager.getCurrentPhase()}' phase of ${turn.name}! Orders for this phase are due by midnight.`,
      );
    } else if (daysOnTurn == daysPerTurn - 1) {
      ServerManager.ordersChannel.send(
        `Orders for the '${Phase.Diplomatic}' phase are due tomorrow at noon!`,
      );
    }
  }

  if (currentTime === "00:00" && phase === Phase.RetreatAndReinforce) {
    ServerManager.ordersChannel.send(StateManager.getOrdersMessage());
    StateManager.nextPhase();
    ServerManager.gameStateChannel.send(
      `We are now on the '${StateManager.getCurrentPhase()}' phase of ${turn.name}!`,
    );
  }
});

export default phaseJob;

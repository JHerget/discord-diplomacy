import { OrdersRegex } from "./regex.js";
import {
  ConvoyOrder,
  DisbandOrder,
  MoveOrder,
  ReinforceOrder,
  SupportOrder,
} from "./order-models.js";
import StateManager from "../state/state-manager.js";
import { Phase } from "../state/state-models.js";

export class OrderManager {
  static isOrderFormatValid(orderString) {
    const order = orderString.toLowerCase().trim();

    switch (StateManager.getCurrentPhase()) {
      case Phase.Diplomatic:
        if (OrdersRegex.Support.test(order)) {
          return true;
        }
        if (OrdersRegex.Convoy.test(order)) {
          return true;
        }
        if (OrdersRegex.Move.test(order)) {
          return true;
        }
        if (OrdersRegex.Reinforce.test(order)) {
          return false;
        }
        if (OrdersRegex.Disband.test(order)) {
          return false;
        }
        break;
      case Phase.RetreatAndReinforce:
        if (OrdersRegex.Support.test(order)) {
          return false;
        }
        if (OrdersRegex.Convoy.test(order)) {
          return false;
        }
        if (OrdersRegex.Move.test(order)) {
          return true;
        }
        if (
          OrdersRegex.Reinforce.test(order) &&
          this.game.currentTurn % 2 == 0
        ) {
          return true;
        }
        if (OrdersRegex.Disband.test(order)) {
          return true;
        }
        break;
    }

    return false;
  }

  static parseOrder(orderString) {
    const order = orderString.toLowerCase().trim();

    if (OrdersRegex.Support.test(order)) {
      return new SupportOrder(order);
    }
    if (OrdersRegex.Convoy.test(order)) {
      return new ConvoyOrder(order);
    }
    if (OrdersRegex.Move.test(order)) {
      return new MoveOrder(order);
    }
    if (OrdersRegex.Reinforce.test(order)) {
      return new ReinforceOrder(order);
    }
    if (OrdersRegex.Disband.test(order)) {
      return new DisbandOrder(order);
    }
  }
}

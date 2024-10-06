import { OrdersRegex } from "./regex.js";

export const UnitType = {
  Army: "Army",
  Fleet: "Fleet",
};

export const OrderType = {
    Move: "Move",
    Support: "Support",
    Convoy: "Convoy",
    Reinforce: "Reinforce",
    Disband: "Disband"
}

export class MoveOrder {
  constructor(orderString) {
    let data = orderString.match(OrdersRegex.Move).slice(1);

    this.orderType = OrderType.Move;
    this.unitType = data[0] == "a" ? UnitType.Army : UnitType.Fleet;
    this.location = data[1];
    this.destination = ["h", "hold"].includes(data[2]) ? data[1] : data[2];
  }
}

export class SupportOrder {
  constructor(orderString) {
    let data = orderString.match(OrdersRegex.Support).slice(1);

    this.orderType = OrderType.Support;
    this.unitType = data[0] == "a" ? UnitType.Army : UnitType.Fleet;
    this.location = data[1];
    this.moveOrder = new MoveOrder(`${data[2]} ${data[3]}-${data[4]}`);
  }
}

export class ConvoyOrder {
  constructor(orderString) {
    let data = orderString.match(OrdersRegex.Convoy).slice(1);

    this.orderType = OrderType.Convoy;
    this.unitType = UnitType.Fleet;
    this.location = data[0];
    this.moveOrder = new MoveOrder(`a ${data[1]}-${data[2]}`);
  }
}

export class ReinforceOrder {
  constructor(orderString) {
    let data = orderString.match(OrdersRegex.Reinforce).slice(1);

    this.orderType = OrderType.Reinforce;
    this.unitType = data[0] == "a" ? UnitType.Army : UnitType.Fleet;
    this.location = data[1];
  }
}

export class DisbandOrder {
  constructor(orderString) {
    let data = orderString.match(OrdersRegex.Disband).slice(1);

    this.orderType = OrderType.Move;
    this.location = data[0];
  }
}

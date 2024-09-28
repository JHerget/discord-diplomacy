export const Phase = {
  Diplomatic: "Diplomacy and Order Writing",
  RetreatAndReinforce: "Retreat and Reinforce",
};

export class Game {
  constructor() {
    this.currentTurn = 1;
    this.inProgress = false;
    this.daysPerTurn = 3;
    this.players = [];
    this.turns = [new Turn("Spring 1901", this.currentTurn)];
  }
}

export class Player {
  constructor(power, globalName, userId) {
    this.power = power;
    this.globalName = globalName;
    this.userId = userId;
    this.active = true;
  }
}

export class Turn {
  constructor(name, number) {
    this.name = name;
    this.number = number;
    this.daysOnTurn = 0;
    this.phase = Phase.Diplomatic;
    this.orders = [];
  }
}

export class Orders {
  constructor(player, phase, orders) {
    this.player = player;
    this.phase = phase;
    this.orders = orders;
  }
}

export const Phase = {
  Diplomatic: "Diplomacy and Order Writing",
  Retreat: "Retreat and Disband",
  Reinforce: "Gain and Lose Units",
};

export class Game {
  constructor() {
    this.currentTurn = 1;
    this.inProgress = false;
    this.players = [];
    this.turns = [new Turn("Spring 1901", this.currentTurn)];
  }
}

export class Player {
  constructor(power, userId) {
    this.power = power;
    this.userId = userId;
  }
}

export class Turn {
  constructor(name, number) {
    this.name = name;
    this.number = number;
    this.phase = Phase.Diplomatic;
    this.orders = [];
    this.retreats = [];
    this.reinforcements = [];
  }
}

export class Orders {
  constructor(player, phase, orders) {
    this.player = player;
    this.phase = phase;
    this.orders = orders;
  }
}

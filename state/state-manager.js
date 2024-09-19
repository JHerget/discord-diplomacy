import fs from "node:fs";
import { Game, Orders, Phase, Player, Turn } from "./state-models.js";

class GameStateManager {
  constructor() {
    this.stateFile = "./state/state.json";
    this.game = new Game();

    try {
      const gameData = fs.readFileSync(this.stateFile, "utf8");
      this.game = JSON.parse(gameData);
    } catch {
      console.log("Starting a new game!");
    }
  }

  saveState() {
    const newData = JSON.stringify(this.game, null, 4);
    fs.writeFileSync(this.stateFile, newData);
  }

  addOrders(userId, orders) {
    const turn = this.game.turns[this.game.currentTurn - 1];
    const newOrders = new Orders(
      new Player(this.getPowerName(), userId),
      turn.phase,
      orders,
    );
    let ordersUpdated = false;

    for (let i = 0; i < turn.orders.length; i++) {
      const order = turn.orders[i];

      if (order.player.userId == userId && order.phase == turn.phase) {
        turn.orders[i] = newOrders;
        ordersUpdated = true;
      }
    }

    if (!ordersUpdated) {
      turn.orders.push(newOrders);
    }

    this.saveState();

    return true;
  }

  getOrders() {
    const turn = this.game.turns[this.game.currentTurn - 1];

    return turn.orders.filter((order) => order.phase == turn.phase);
  }

  isRegistered(userId) {
    let isRegistered = false;

    this.game.players.forEach((player) => {
      if (player.userId === userId) {
        isRegistered = true;
        return;
      }
    });

    return isRegistered;
  }

  registerPlayer(power, userId) {
    let canRegister = true;

    this.game.players.forEach((player) => {
      if (player.power === power) {
        canRegister = false;
        return;
      }
    });

    if (!canRegister) {
      return false;
    }

    this.game.players.push(new Player(power, userId));
    this.saveState();

    return true;
  }

  deregisterPlayer(userId) {
    const players = this.game.players;

    for (let i = 0; i < players.length; i++) {
      if (players[i].userId === userId) {
        this.game.players.splice(i, 1);
        this.saveState();
        return;
      }
    }
  }

  getPowerName(userId) {
    let power = null;

    this.game.players.forEach((player) => {
      if (player.userId === userId) {
        power = player.power;
        return;
      }
    });

    return power;
  }

  getTurnName() {
    const turn = this.game.currentTurn;
    const year = 1901 + Math.floor(turn / 2);
    const season = turn % 2 == 1 ? "Spring" : "Fall";

    return `${season} ${year}`;
  }

  getAvailablePowers() {
    const allPowers = [
      "Austria-Hungary",
      "England",
      "France",
      "Germany",
      "Italy",
      "Russia",
      "Turkey",
    ];
    const currentPowers = this.game.players.map((player) => player.power);
    let availablePowers = [];

    allPowers.forEach((power) => {
      if (!currentPowers.includes(power)) {
        availablePowers.push(power);
      }
    });

    return availablePowers;
  }

  getCurrentTurn() {
    const turn = this.game.turns[this.game.currentTurn - 1];
    let turnData = {
      name: turn.name,
      phase: turn.phase,
      orders: [],
    };

    switch (turn.phase) {
      case Phase.Diplomatic:
        turnData.orders = turn.orders.map((order) => order.power);
        break;
      case Phase.Retreat:
        turnData.orders = turn.retreats.map((order) => order.power);
        break;
      case Phase.Reinforce:
        turnData.orders = turn.reinforcements.map((order) => order.power);
        break;
    }

    return turnData;
  }

  getCurrentPhase() {
    return this.game.turns[this.game.currentTurn - 1].phase;
  }

  nextPhase() {
    const turn = this.game.turns[this.game.currentTurn - 1];

    switch (turn.phase) {
      case Phase.Diplomatic:
        this.game.phase = Phase.Retreat;
        this.saveState();
        break;
      case Phase.Retreat:
        this.game.phase = Phase.Reinforce;
        this.saveState();
        break;
      case Phase.Reinforce:
        this.newTurn();
        break;
    }
  }

  newTurn() {
    this.game.currentTurn += 1;
    this.game.turns.push(new Turn(this.getTurnName(), this.game.currentTurn));

    this.saveState();
  }

  inProgress() {
    return this.game.inProgress;
  }
}

const StateManager = new GameStateManager();
export default StateManager;

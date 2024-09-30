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
      new Player(this.getPowerName(userId), userId),
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

  registerPlayer(power, user) {
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

    this.game.players.push(new Player(power, user.globalName, user.id));
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

  getAllPlayers() {
    return this.game.players;
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
    const year = 1900 + Math.ceil(turn / 2);
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
    return this.game.turns[this.game.currentTurn - 1];
  }

  getCurrentPhase() {
    return this.getCurrentTurn().phase;
  }

  setCurrentPhase(phase) {
    this.game.turns[this.game.currentTurn - 1].phase = phase;
    this.saveState();
  }

  nextPhase() {
    switch (this.getCurrentPhase()) {
      case Phase.Diplomatic:
        this.setCurrentPhase(Phase.RetreatAndReinforce);
        break;
      case Phase.RetreatAndReinforce:
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

  getDaysPerTurn() {
    return this.game.daysPerTurn;
  }

  allOrdersSubmitted() {
    const ordersReceived = new Set(
      this.getOrders()
        .filter((order) => order.player.active)
        .map((order) => order.player.userId),
    );
    const registeredPlayers = new Set(
      this.game.players
        .filter((player) => player.active)
        .map((player) => player.userId),
    );

    return ordersReceived.size >= registeredPlayers.size;
  }

  incrementDaysOnTurn() {
    this.game.turns[this.game.currentTurn - 1].daysOnTurn += 1;
    this.saveState();
  }

  getOrdersMessage() {
    const turn = this.getCurrentTurn();

    const orders = this.getOrders();
    let message = `## ${turn.name}\n${turn.phase}\n\n`;

    for (const order of orders) {
      message += `**${order.player.power}**\n`;
      message += `${order.orders.join("\n")}\n\n`;
    }

    return message;
  }
}

const StateManager = new GameStateManager();
export default StateManager;

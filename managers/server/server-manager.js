class ServerStateManager {
  constructor() {
    this.client = null;
    this.gameStateChannel = null;
    this.ordersChannel = null;
    this.greatPowerRole = null;
  }

  setClient(client) {
    this.client = client;
  }

  setGameStateChannel(gameStateChannel) {
    this.gameStateChannel = gameStateChannel;
  }

  setOrdersChannel(ordersChannel) {
    this.ordersChannel = ordersChannel;
  }

  setGreatPowerRole(greatPowerRole) {
    this.greatPowerRole = greatPowerRole;
  }
}

const ServerManager = new ServerStateManager();
export default ServerManager;

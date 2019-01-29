
Web3 = require('web3');

module.exports = {
  contracts_directory: "contracts",

  networks: {
    dev: {
      host: "localhost",
      port: 8545,
      network_id: "*"
    },
    devui: {
      host: "localhost",
      port: 7545,
      network_id: '*'
    },
    production: {
      network_id: "*",
      provider: () => {
        return new Web3.providers.HttpProvider(process.env.PRODUCTION_URL);
      }
    }
  }
};

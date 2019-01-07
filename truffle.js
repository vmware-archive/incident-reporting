/*
 * NB: since truffle-hdwallet-provider 0.0.5 you must wrap HDWallet providers in a 
 * function when declaring them. Failure to do so will cause commands to hang. ex:
 * ```
 * mainnet: {
 *     provider: function() { 
 *       return new HDWalletProvider(mnemonic, 'https://mainnet.infura.io/<infura-key>') 
 *     },
 *     network_id: '1',
 *     gas: 4500000,
 *     gasPrice: 10000000000,
 *   },
 */
Web3 = require('web3');

module.exports = {
  // See <http://truffleframework.com/docs/advanced/configuration>
  // to customize your Truffle configuration!

  networks: {
    dev: {
      host: "localhost",
      port: 8545,
      network_id: "*" // Match any network id
    },
    devui: {
      host: "localhost",
      port: 7545,
      network_id: '*'
    },
    vmware: {
      network_id: "*",
      provider: () => {
        return new Web3.providers.HttpProvider("https://dev@blockchain.local:XXXXXXXX@mgmt.blockchain.vmware.com/blockchains/b20bef33-b35d-44a6-b5b9-15423a25d31f/api/concord/eth");
      }
    }
  }
};

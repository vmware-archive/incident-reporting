var IncidentLog = artifacts.require("./IncidentLog.sol");

module.exports = function(deployer) {
  deployer.deploy(IncidentLog);
};

// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: BSD-2
var IncidentLog = artifacts.require("./IncidentLog.sol");

module.exports = function(deployer) {
  deployer.deploy(IncidentLog);
};

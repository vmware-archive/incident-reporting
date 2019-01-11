# Incident reporting blockchain demo

## Scenario

MedCo is using large data sets that live in AWS s3 for a research application.  The results of the application should be kept encrypted and moved into the corporate private cloud for storage and later use.  This data will be sold to a large pharmaceutical company.

Loss of integrity for this data is critical to the business’s IP and should be reported and tracked.  In addition to MedCo needing to verify integrity, the purchasing Pharma company will need to verify their purchase has not been compromised.

It could be possible for an employee or a bad actor to thwart the privacy of this data and wipe record of that event taking place.  In order to ensure compliance, we need some method of tracking and preventing alteration of reported incidents.

Blockchain allows for all parties to verify no incidents were reported, or that the incidents reported were handled correctly, and that no tampering of they system has occurred.

Possible con: latency.  Will incidents need to be known in real time, or is some delay reasonable?

On detection of a security incident, a log entry should be made.

### Example incidents

* VMs detected to be running vulnerable software
* Ssh access to VM running the app
* Breach of encryption key security
* AppDefense identifies a problem

## Deploy the contracts

This only needs to be done once for a demo.  Save the ID for IncidentLog that is output by the deploy below.  You'll need it when running the server.

## Install the dependencies:

* install npm
* install go
* install go-ethereum from https://geth.ethereum.org/install/

``` bash
npm install -g truffle@4.1.14
```

edit truffle.js to add a section like:

``` javascript
    vmware: {
      network_id: "*",
      provider: () => {
        return new Web3.providers.HttpProvider("https://dev@blockchain.local:XXXX@$@mgmt.blockchain.vmware.com/blockchains/XXXX/api/concord/eth");
      }
```

Then run the deploy
``` bash
truffle deploy --network vmware --reset

... output ...
Compiling ./contracts/IncidentLog.sol...
Compiling ./contracts/Migrations.sol...
Writing artifacts to ./build/contracts

Using network 'vmware'.

Running migration: 1_initial_migration.js
  Deploying Migrations...
  ... 0x6eed7a412e8761730dc032eb0586aaaac943ed959dfff8024211f27afa3df4ed
  Migrations: 0xd4b650638c525758e57a99e7a85508a708bace8c
Saving successful migration to network...
  ... 0x7a8d864d0d7b2c4a5b85a05bc9e931407edd30a9b83d1c500c3a272e6aec9e14
Saving artifacts...
Running migration: 2_deploy_contracts.js
  Deploying IncidentLog...
  ... 0x99a050fbad55ad44644d1a094781d9f72380fbe5f8e98688dd3f87b4d6c172c6
  IncidentLog: 0xe30267a60f8bd9103d5988734aab9c3b0acb011f
Saving successful migration to network...
  ... 0x22eafa16d629a5a0be06e447605ecef2748379760ba9aede4475565c376f6357
Saving artifacts...
```
Save the ID following `IncidentLog:` for use when setting your environment for running the server below.

## Running the service manually

``` bash
npm install

# create a new account and keys.
# Save the passhrase and account file for next steps
geth account new
geth account list
... output ...
# Account #0: {d55010663ffd36aba75d2b29eba24015d6e20671} keystore:///Users/xxx/Library/Ethereum/keystore/UTC--2018-11-26T20-34-57.266338000Z--d55010663ffd36aba75d2b29eba24015d6e20671
# Account #1: {fe00bb37a56282d33680542ae1cd6763660b4812} keystore:///Users/xxx/Library/Ethereum/keystore/UTC--2019-01-03T19-49-59.140692000Z--fe00bb37a56282d33680542ae1cd6763660b4812
# Account #2: {1bf6a382b68444bf1578d0647bed98990367d5cd} keystore:///Users/xxx/Library/Ethereum/keystore/UTC--2019-01-11T20-12-23.689636000Z--1bf6a382b68444bf1578d0647bed98990367d5cd
...

# set up your environment
export CLIENT_PASSPHRASE=password_for_account_created_above
export CLIENT_KEYFILE=/Users/xxx/Library/Ethereum/keystore/UTC--2019-01-03T19-49-59.140692000Z--fe00bb37a56282d33680542ae1cd6763660b4812
export CLIENT_CONTRACT_ADDRESS=id_from_above_output_from_truffle_deploy # address for the deployed contract
export CLIENT_URL=mgmt.blockchain.vmware.com/blockchains/xxxxx/api/concord/eth
export CLIENT_USER=dev@blockchain.local
export CLIENT_PASSWORD=xxxxx

go generate
go build
./incident-reporting
```

### post an incident

    curl --header "Content-Type: application/json"   --request POST   --data '{"Reporter":"0xFE00BB37A56282d33680542Ae1CD6763660b5555","Message":"automatic reporting"}' localhost/rest/log

### check it

    curl --header "Content-Type: application/json"   --request GET localhost/rest/log/0

### use the ui

Open a browser to http://localhost/log
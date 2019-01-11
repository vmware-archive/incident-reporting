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

Dependencies

* install npm
* install go
* install go-ethereum https://geth.ethereum.org/install/

``` bash
npm install -g truffle
npm install

# create a new account and keys.
# Save the passhrase and account id for next steps
geth account new
geth account list


# deploy the contracts
truffle deploy --network vmware --reset

export CLIENT_PASSPHRASE=password_for_account_created_above
export CLIENT_CONTRACT_ADDRESS=0xdb3d71898f878bc5e6ef6e0de985a55ca483c0c0 # address for the deployed contract
export CLIENT_URL=mgmt.blockchain.vmware.com/blockchains/xxxxx/api/concord/eth
export CLIENT_USER=dev@blockchain.local
export CLIENT_PASSWORD=xxxxx
export CLIENT_KEYFILE=/Users/xxx/Library/Ethereum/keystore/UTC--2019-01-03T19-49-59.140692000Z--fe00bb37a56282d33680542ae1cd6763660b4812


```
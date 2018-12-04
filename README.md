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

Setting up and running

* install npm

``` bash

install npm
npm install truffle -g
npm install solc@0.4.24+commit.e67f0147.Emscripten.clang
npm install -g solium
npm install -g ganache-cli
npm install drizzle
start ganach ui
cd client
npm start

```
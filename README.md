# Incident Reporting

## Overview

Reporting and tracking of security incidents is important so that all parties affected by the event know the status and can respond in concert.
Current methods to track incidents are generally paper-based manual processes.
More recent systems are based on a centralized database with a web interface to report incidents and interact with them.

Manual incident reporting process is not amenable to scenarios where security incidents may affect more than a single entity or organization.
For example, a security breach in the supply chain for a food manufacturer could result in several retail businesses with products on shelf that contain a pathogen.
Current methods of notifying the proper authorities and consumers require a phone tree to call all the correct parties, which then react as individuals or local committees.

This code repository implements a system that allows all interested parties (retail, government, public, etc) to

* report incidents in an automated, real-time fashion
* see the incident as soon as it is reported
* respond in parallel

Additionally, the record it resistant to malicious attacks and alteration of the records.

## Scenarios this app helps

### Pharma Data Sales

MedCo is using large data sets that live in AWS s3 for a research application.  The results of the application should be kept encrypted and moved into the corporate private cloud for storage and later use.  This data will be sold to a large pharmaceutical company.

Loss of integrity for this data is critical to the business’s IP and should be reported and tracked.  In addition to MedCo needing to verify integrity, the purchasing Pharma company will need to verify their purchase has not been compromised.

It could be possible for an employee or a bad actor to thwart the privacy of this data and wipe record of that event taking place.  In order to ensure compliance, we need some method of tracking and preventing alteration of reported incidents.

This incident reporting application allows for all parties to verify that no incidents were reported, or that the incidents reported were handled correctly, and that no tampering or removal of incidents has occurred.

Possible con: latency.  Will incidents need to be known in real time, or is some delay reasonable?

### Cross Enterprise Application Security Monitoring

Big-Co has many smaller business units that develop and run business critical applications.
Company wide IT has responsibility for deploying servers, networking equipment, hypervisors, platform stacks and other infrastructure.

Red-BU may be responsible for running their own application in production on top of IT furnished infrastructure, while Blue-BU develops apps with SRE-BU who then runs those apps on the IT provided infrastructure.
IT is responsible for managing security incidents at the network and infrastructure levels, while SRE-BU is responsible for security incidents in Blue-Apps, and Red-BU product management is responsible for Red-Apps security incidents.

Meanwhile, HR may need to be involved if an employee was involved in a security incident, and police/FBI may need to be involved if there is an external security threat.

Security incidents here need to be globally reportable, reacted to by various parties with different concerns, and with some parties gaining advantage if they can alter or prevent recording incidents.

This application of blockchain based incident reporting meets those requirements with regard to all teams, internal and external that can access the system.
Additionally, no single entity needs to be in control of the infrastructure to support this action in concert.

## Running the app

### Install the dependencies

* Docker
* docker-compose
* make

### Build the ui and truffle containers

    make

### Deploy the contracts

    export PRODUCTION_URL=https://dev@blockchain.local:XXXX@$@mgmt.blockchain.vmware.com/blockchains/XXXX/api/concord/eth
    make run-truffle

Sample output:

``` bash
Compiling ./contracts/IncidentLog.sol...
Compiling ./contracts/Migrations.sol...
Writing artifacts to ./build/contracts

Using network 'production'.

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

### Running the UI service

Edit the `env` file to update

* `CLIENT_CONTRACT_ADDRESS` line to match your IncidentLog id from above
* `CLIENT_URL` to be the url to your blockchain's Ethereum API endpoint
* `CLIENT_USER` and `CLIENT_PASSWORD` to be your credentials at that endpoint

Then run

     make run-ui

#### post an incident

    curl --header "Content-Type: application/json"   --request POST   --data '{"Reporter":"0xFE00BB37A56282d33680542Ae1CD6763660b5555","Message":"automatic reporting"}' localhost:8080/rest/log

#### check it

    curl --header "Content-Type: application/json"   --request GET localhost:8080/rest/log/0

#### use the ui

Open a browser to http://localhost:8080/log

## Fast demo run through

    # Edit run-demo.sh
    ./run-demo.sh

## Contributing

The incident-reporting project team welcomes contributions from the community. Before you start working with incident-reporting, please read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Copyright 2019 VMware, Inc.
SPDX-License-Identifier: BSD-2

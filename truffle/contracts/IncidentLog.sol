// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: BSD-2
pragma solidity ^0.4.24;

/// @title A log for recording incidents
/// @author Tom Scanlan
contract IncidentLog {

    struct Incident {
        address reporter;  // what public id reported this incident?
        string message;    // log message for the incident
        uint timestamp;    // time the incident was committed to the blockchain
        string location;   // where did the event take place
        bool resolved;     // has this issue been resolved?
    }

    // This event will fire each time an incident is reported
    event FireIncident (
        address reporter,
        string message
    );

    event GotCalled();

    // A dynamically-sized array of `Incidents` structs.
    Incident[] public incidents;

    function reportIncident (address reporter, string memory message, string memory location) public {
        emit GotCalled();
        uint timestamp = now;
        incidents.push(
            Incident({
                reporter: reporter,
                message: message,
                timestamp: timestamp,
                location: location,
                resolved: false
            })
        );
        emit FireIncident(reporter, message);
    }

    function getCount () public view returns (uint256) {
        emit GotCalled();
        return incidents.length;
    }

    function getIncident (uint256 n) public view returns (address, string memory, uint, string, bool) {
        emit GotCalled();
        require(incidents.length != 0, "no log entries yet");
        require(n < incidents.length, "requested entry doesn't exist");
        Incident storage i = incidents[n];
        return (i.reporter, i.message, i.timestamp, i.location, i.resolved);
    }
}

pragma solidity ^0.4.24;

/// @title A log for recording incidents
/// @author Tom Scanlan
/// @dev Storage is costly and should only be used for critical data
contract IncidentLog {

    struct Incident {
        address reporter;  // what public id reported this incident?
        string message;    // log message for the incident
        uint timestamp;
    }

    // This event will fire each time an incident is reported
    event FireIncident (
        address reporter,
        string message,
        uint timestamp
    );

    // A dynamically-sized array of `Incidents` structs.
    Incident[] public incidents;

    function reportIncident (address reporter, string memory message) public {
        uint timestamp = now;
        incidents.push(
            Incident({
                reporter: reporter,
                message: message,
                timestamp: timestamp
            })
        );
        emit FireIncident(reporter, message, timestamp);
    }

    function getCount () public view returns (uint256) {
        return incidents.length;
    }

    function getIncident (uint256 n) public view returns (address, string memory, uint) {
        require(incidents.length != 0, "no log entries yet");
        require(n < incidents.length, "requested entry doesn't exist");
        Incident storage i = incidents[n];
        return (i.reporter, i.message, i.timestamp);
    }
}

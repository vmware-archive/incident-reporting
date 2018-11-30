var IncidentLog = artifacts.require("IncidentLog")


var il = IncidentLog.at(IncidentLog.address)
il.reportIncident(web3.eth.accounts[0], 'testing')

contract("IncidentLog", function(accounts) {

    let il

    beforeEach('setup contract for each test', async function () {
        il = await IncidentLog.new()
    })

    it("should be deployed", function() {
        return IncidentLog.deployed()
    })

    it("should start out empty", function() {
        return assert.equal(0, il.incidents.length, "didn't start out empty");
    })

    it("should allow adding a message", async function() {
        await il.reportIncident(accounts[0], "here is an entry");
        incident = await il.getIncident(0);
        count = await il.getCount()
        assert.equal("here is an entry", incident[1], "got wrong message");
        assert.equal(count, 1, "got wrong number of messages");
    })
})
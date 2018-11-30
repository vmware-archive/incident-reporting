import React from "react";

class GetLogLine extends React.Component {
    state = { dataKey: null };

    
  handleKeyDown = e => {
    // if the enter key is pressed, set the value with the string
    if (e.keyCode === 13) {
      this.setValue(e.target.value);
    }
  };
  
  setValue = value => {
    const { drizzle, drizzleState } = this.props;
    const contract = drizzle.contracts.IncidentLog;

    const dataKey = contract.methods.getIncident.cacheCall(value);
    this.setState({ dataKey });
  };


    render() {
        // get the contract state from drizzleState
        const { IncidentLog } = this.props.drizzleState.contracts;

        console.log(IncidentLog);
        console.log(this.state.dataKey);
        // using the saved `dataKey`, get the variable we're interested in
        const logline = IncidentLog.getIncident[this.state.dataKey];
        console.log(logline);

        if (logline && logline.value) {
            const address = logline.value[0];
            const message = logline.value[1];
            const time = logline.value[2];

            // if it exists, then we display its value
            return (
                <div>
                    <input type="text" onKeyDown={this.handleKeyDown} />
                    <p>address: {address}</p>
                    <p>message: {message}</p>
                    <p>time: {time}</p>
                </div>);

        } else {
            return (
                <div>
                    <input type="text" onKeyDown={this.handleKeyDown} />
                </div>);
        }
    }
}

export default GetLogLine;
import React from "react";

class ReadString extends React.Component {
  state = { dataKey: null };

  componentDidMount() {
    const { drizzle } = this.props;
    const contract = drizzle.contracts.IncidentLog;

    // let drizzle know we want to watch the `myString` method
    const dataKey = contract.methods["getCount"].cacheCall();
    console.log(contract);
    console.log(dataKey);
    // save the `dataKey` to local component state for later reference
    this.setState({ dataKey });
  }

  render() {
    // get the contract state from drizzleState
    const { IncidentLog } = this.props.drizzleState.contracts;

    // using the saved `dataKey`, get the variable we're interested in
    console.log(IncidentLog);
    const myString = IncidentLog.getCount[this.state.dataKey];

    // if it exists, then we display its value
    return <p>My stored string: {myString && myString.value}</p>;
  }
}

export default ReadString;
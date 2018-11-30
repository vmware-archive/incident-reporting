import React, { Component } from 'react';
import './App.css';
import ReadCount from "./ReadCount";
import ReportIncident from "./ReportIncident";
import GetLogLine from "./GetLogLine"


class App extends Component {
  state = { loading: true, drizzleState: null };

  componentDidMount() {
    const { drizzle } = this.props;
  
    // subscribe to changes in the store
    this.unsubscribe = drizzle.store.subscribe(() => {
  
      // every time the store updates, grab the state from drizzle
      const drizzleState = drizzle.store.getState();
  
      // check to see if it's ready, if so, update local component state
      if (drizzleState.drizzleStatus.initialized) {
        this.setState({ loading: false, drizzleState });
      }
    });
  }
  
  compomentWillUnmount() {
    this.unsubscribe();
  }

  render() {
    if (this.state.loading) return "Loading Drizzle...";
    return (
      <div className="App">
        <ReadCount
          drizzle={this.props.drizzle}
          drizzleState={this.state.drizzleState}
        />

        <ReportIncident
          drizzle={this.props.drizzle}
          drizzleState={this.state.drizzleState}
        />

        <GetLogLine
          drizzle={this.props.drizzle}
          drizzleState={this.state.drizzleState}
        />
      </div>
    );
  }
}

export default App;

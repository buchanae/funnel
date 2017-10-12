import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';

class App extends Component {
  constructor() {
    super()
    this.state = {data: {}}
  }

  componentWillMount() {
    var app = this
    fetch("./v1/tasks/b79houqrl6qts4qurblg?view=FULL").then(
      function(response) {
        response.json().then(function(data) {
          console.log(data, "foo")
          app.setState({ data })

        });
      })
  }

  render() {

    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <p className="App-intro">
          {this.state.data.id}
        </p>
        <p className="App-intro">
          {this.state.data.name}
        </p>
        <p className="App-intro">
          {this.state.data.state}
        </p>
        <p className="App-intro">
          CPUs: {this.state.data.resources.cpu_cores}, RAM: {this.state.data.resources.ram_gb}
        </p>
      </div>
    );
  }
}

export default App;

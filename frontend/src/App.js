import React, { Component } from 'react';
import './App.css';
import User from './User';
import Orgs from './Orgs';
import Donate from './Donate';

class App extends Component {
    render() {
        return (
            <div className="container">
            <div className="row">
            <div className="col-lg-12 text-center">
            <h1>DonorBuddy</h1>
            </div>
            </div>
            <User/>
            <Orgs/>
            <Donate/>
            </div>
        );
    }
}

export default App;

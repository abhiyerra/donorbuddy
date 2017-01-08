import React, { Component } from 'react';
import Org from './Org';

class Orgs extends Component {
  constructor() {
    super();
    this.state = {
      orgs: []
    };
  }

  getOrgs() {
    fetch("/v1/orgs").then(data => data.json()).then(orgs => {
      let orgTags = [];
      for(let i = 0; i < orgs.length; i++) {
        orgTags.push(<Org key={"org"+i} org={orgs[i]} />);
      }
      this.setState({orgs: orgTags});
    });
  }

  componentDidMount() {
    this.getOrgs();
  }

  render() {
    return (
      <div className="row">
      <div className="col-lg-12">
      <div>

      </div>
      <table className="table">
      <tbody>
      {this.state.orgs}
      </tbody>
      </table>
      </div>
      </div>
    );
  }
}

export default Orgs;

import React, { Component } from 'react';
import $ from 'jquery';
import Org from './Org';

class Orgs extends Component {
  constructor() {
    super();
    this.state = {
      State: '',
      City: '',
      Category: '',
      orgs: []
    };

    this.handleChangeCity = this.handleChangeCity.bind(this);
    this.handleChangeState = this.handleChangeState.bind(this);
    this.handleChangeCategory = this.handleChangeCategory.bind(this);
    this.search = this.search.bind(this);
  }

  getOrgs() {
    $.get(global.APIServer+"/v1/orgs", {
      state: this.state.State,
      city: this.state.City,
      category: this.state.Category,
    }, orgs => {
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

  handleChangeCity(event) {
    this.setState({City: event.target.value});
  }
  handleChangeState(event) {
    this.setState({State: event.target.value});
  }
  handleChangeCategory(event) {
    this.setState({Category: event.target.value});
  }

  search(e) {
    e.preventDefault();
    this.getOrgs();
  }


  render() {
    return (
      <div className="row">
        <div className="col-lg-12 text-center">
          <h2>Find Orgs</h2>
          <hr className="star-primary"/>
        </div>

        <div className="col-lg-12">
          <div className="form-group">
            <label for="usr">City:</label>
            <input type="text" className="form-control" value={this.state.City} onChange={this.handleChangeCity}/>
          </div>
          <div className="form-group">
            <label for="pwd">State:</label>
            <input type="text" className="form-control"  value={this.state.State} onChange={this.handleChangeState}/>
          </div>
          <div className="form-group">
            <label for="pwd">Category:</label>
            <input type="text" className="form-control" value={this.state.Category} onChange={this.handleChangeCategory}/>
          </div>
          <a className="btn btn-lg btn-success" onClick={this.search}>Search</a>

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

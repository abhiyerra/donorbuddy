import React, { Component } from 'react';
import $ from 'jquery';
import Donate from './Donate';

class DonationList extends Component {
  constructor(props) {
    super(props);

    this.handleRemove = this.handleRemove.bind(this);
  }

  getOrgs() {
    let orgs = [];

    if(this.props.user == null) {
      return orgs;
    }


    for(let i = 0; i < this.props.user.Orgs.length; i++) {
      orgs.push(
        <tr>
          <td>{this.props.user.Orgs[i].Name}</td>
          <td>{93 / this.props.user.Orgs.length}%</td>
          <td><a href="#" onClick={this.handleRemove(`${global.APIServer}/v1/user/org/${this.props.user.Orgs[i].Id}`)}>Remove</a></td>
        </tr>
      );
    }

    orgs.push(
      <tr>
        <td>DonorBuddy Fee</td>
        <td>7%</td>
        <td></td>
      </tr>
    );

      return orgs;
  }

  handleRemove(link) {
    return function(e) {
      e.preventDefault();

      $.ajax({
        url: link,
        type: 'DELETE',
        success: data => {
          console.log(data);
        }
      });
    }
  }

  render() {
    if(this.props.user == null) {
      return null;
    }

    if(this.props.user.Orgs == null) {
      return (
      <div className="row">
        <div className="col-lg-12 text-center">
          <h2>Your Donation List</h2>
          <hr className="star-primary"/>
        </div>

        <div className="col-lg-12 text-center">
          <p>Add some places you want to donate to down below.</p>
        </div>
      </div>

      );
    }

    let orgs = this.getOrgs();
    return (
      <div className="row">
        <div className="col-lg-12 text-center">
          <h2>Your Donation List</h2>
          <hr className="star-primary"/>
        </div>

        <div className="col-lg-12">
          <table className="table">
            <tbody>
              {orgs}
            </tbody>
          </table>
        </div>
      </div>
    );
  }
}


class User extends Component {
  constructor() {
    super();

    this.state = {
      user: null
    };
  }

  componentDidMount() {
    $.get(global.APIServer+"/v1/user", (user) => {
      if(user.Error !== undefined) {
        user = null;
      }

      this.setState({
        user: user
      });
    });
  }

  render() {
    return (
      <div>
        <Donate user={this.state.user} />
        <DonationList user={this.state.user} />
      </div>
    );
  }
}

export default User;

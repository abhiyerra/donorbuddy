import React, { Component } from 'react';
import $ from 'jquery';

class User extends Component {
  constructor() {
    super();

    this.state = {
      loginButton: null,
      user: null
    };


    this.handleRemove = this.handleRemove.bind(this);
  }


  componentDidMount() {
    $.get("/v1/user", user => {
      console.log(user);
      if(user.Error !== undefined) {
        this.setState({
          loginButton: (
            <a href={"/auth/login"} className="btn btn-lg btn-primary">Login with Facebook</a>
          )
        });
      } else {
        this.setState({
          user: user
        });
      }
    });
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

  orgInfo() {
    if(this.state.user == null) {
      return null
    }

    if(this.state.user.Orgs == null) {
      return null;
    }

    let orgs = [];
    for(let i = 0; i < this.state.user.Orgs.length; i++) {
      orgs.push(
        <tr>
          <td>{this.state.user.Orgs[i].Name}</td>
          <td>{100 / this.state.user.Orgs.length}%</td>
          <td><a href="#" onClick={this.handleRemove(`/v1/user/org/${this.state.user.Orgs[i].Id}`)}>Remove</a></td>
        </tr>
      );
    }

    return (
      <div className="col-lg-12">
        <h2>Current Donation Breakdown</h2>
        <table className="table">
          <tbody>
            {orgs}
          </tbody>
        </table>
      </div>
    );
  }

  render() {
    return (
      <div className="row">
        <div className="col-lg-12 text-center">
          {this.state.loginButton}
        </div>

        {this.orgInfo()}
      </div>
    );
  }
}

export default User;

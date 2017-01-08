import React, { Component } from 'react';

class User extends Component {

  componentDidMount() {
    fetch("/v1/user").then(data => data.json()).then(user => {
      console.log(user);
    });
  }

  render() {
    return (
      <div className="row">
      <div className="col-lg-12 text-center">
      <a href={global.APIServer+"/auth/login"} className="btn btn-lg btn-primary">Login with Facebook</a>
      </div>
      </div>
    );
  }
}

export default User;

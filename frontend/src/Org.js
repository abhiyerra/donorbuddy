import React, { Component } from 'react';

class Org extends Component {
  render() {
    return (
      <tr className="col-lg-3">
            <td>{this.props.org.Name}</td>
            <td>{this.props.org.Category}</td>
            <td><a href={`/v1/user/org/${this.props.org.Id}`}>Add</a></td>
      </tr>
    );
  }
}

export default Org;

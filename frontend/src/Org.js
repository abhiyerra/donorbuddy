import React, { Component } from 'react';
import $ from 'jquery';

class Org extends Component {
  constructor() {
    super();

    this.handleAdd = this.handleAdd.bind(this);
  }

  handleAdd(e) {
    e.preventDefault();

    $.ajax({
      url: `${global.APIServer}/v1/user/org/${this.props.org.Id}`,
      type: 'PUT',
      success: data => {
        console.log(data);
      }
    });
  }

  render() {
    return (
      <tr>
        <td>{this.props.org.Name}</td>
        <td>{this.props.org.Category}</td>
        <td><a href="#" onClick={this.handleAdd}>Add</a></td>
      </tr>
    );
  }
}

export default Org;

import React from 'react';
import { Router, Route } from 'react-router';

import App from './App';
import Home from './Home';

const Routes = (props) => (
  <Router {...props}>
    <Route path="/" component={Home} />
    <Route path="/donate" component={App} />
  </Router>
);

export default Routes;

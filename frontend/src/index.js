import React from 'react';
import ReactDOM from 'react-dom';
import { browserHistory } from 'react-router';
import $ from 'jquery';

import Routes from './routes';

import './index.css';
import 'bootstrap/dist/css/bootstrap.css';

global.APIServer = "http://nanoyak.com:8080";
global.Stripe.setPublishableKey('pk_test_wHqW7mw8SlROUV95cvvaOUZD');

$.ajaxSetup({
    xhrFields: {
       withCredentials: true
    },
    crossDomain: true
});

ReactDOM.render(<Routes history={browserHistory} />, document.getElementById('root'));

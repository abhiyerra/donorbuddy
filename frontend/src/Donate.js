import React, { Component } from 'react';
import $ from 'jquery';
import './Donate.css';

class Login extends Component {
  render() {
    return (
      <div className="row">
        <div className="col-lg-12 text-center">
          <h2>Get Started</h2>
          <hr className="star-primary"/>
        </div>

        <div className="col-lg-12 text-center">
          <p>Get started quickly. Login with Facebook and then you will be able to add to your donation list.</p>

          <a href={global.APIServer+"/auth/login"} className="btn btn-lg btn-primary">Login with Facebook</a>
        </div>
      </div>
    );
  }
}


class Donate extends Component {
  render() {
    console.log(this.props);
    if(this.props.user == null) {
      return <Login />;
    }

    // TODO: Show cancel donate.
    if(this.props.user.HasSubscription) {
      return null;
    }

    return (
      <div className="row">
        <div className="col-lg-12 text-center">
          <h2>Donate</h2>
          <hr className="star-primary"/>
        </div>

        <div className="col-lg-12">
          <div className="panel panel-default credit-card-box">
            <div className="panel-heading display-table" >
              <div className="row display-tr" >
                <h3 className="panel-title display-td" >Payment Details</h3>
                <div className="display-td" >
                  <img className="img-responsive pull-right" src="http://i76.imgup.net/accepted_c22e0.png" />
                </div>
              </div>
            </div>
            <div className="panel-body">
              <form role="form" id="payment-form" method="POST" action="javascript:void(0);" action={global.APIServer+"/v1/payments"}>

                <div className="row">
                  <div className="col-xs-12">
                    <div className="form-group">
                      <label for="cardNumber">Amount to Donate Per Month</label>
                      <div className="input-group">
                        <label class="radio-inline"><input type="radio" name="amount" value="1000" /> $10</label>
                        <label class="radio-inline"><input type="radio" name="amount" value="2000" /> $20</label>
                        <label class="radio-inline"><input type="radio" name="amount" value="5000" /> $50</label>
                        <label class="radio-inline"><input type="radio" name="amount" value="10000" /> $100</label>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="row">
                  <div className="col-xs-12">
                    <div className="form-group">
                      <label for="cardNumber">NAME ON CARD</label>
                      <div className="input-group">
                        <input
                            type="text"
                            className="form-control"
                            name="name"
                            placeholder="Name on Card"
                            autoComplete="cc-number"
                            required autoFocus
                        />
                        <span className="input-group-addon"><i className="fa"></i></span>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="row">
                  <div className="col-xs-12">
                    <div className="form-group">
                      <label for="cardNumber">EMAIL</label>
                      <div className="input-group">
                        <input
                            type="email"
                            className="form-control"
                            name="email"
                            placeholder="Email"
                            autoComplete="cc-number"
                            required autoFocus
                        />
                        <span className="input-group-addon"><i className="fa"></i></span>
                      </div>
                    </div>
                  </div>
                </div>


                <div className="row">
                  <div className="col-xs-12">
                    <div className="form-group">
                      <label for="cardNumber">CARD NUMBER</label>
                      <div className="input-group">
                        <input
                            type="tel"
                            className="form-control"
                            placeholder="Valid Card Number"
                            autoComplete="cc-number"
                            required autoFocus
                            data-stripe="number"
                        />
                        <span className="input-group-addon"><i className="fa fa-credit-card"></i></span>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="row">
                  <div className="col-xs-7 col-md-7">
                    <div className="form-group">
                      <label for="cardExpiry"><span className="hidden-xs">EXPIRATION</span><span className="visible-xs-inline">EXP</span> DATE</label>
                      <input
                          type="tel"
                          className="form-control"
                          placeholder="MM / YY"
                          autoComplete="cc-exp"
                          required
                          data-stripe="exp"
                      />
                    </div>
                  </div>
                  <div className="col-xs-5 col-md-5 pull-right">
                    <div className="form-group">
                      <label for="cardCVC">CV CODE</label>
                      <input
                          type="tel"
                          className="form-control"
                          placeholder="CVC"
                          autoComplete="cc-csc"
                          required
                          data-stripe="cvc"
                      />
                    </div>
                  </div>
                </div>
                <div className="row">
                  <div className="col-xs-12">
                    <button className="subscribe btn btn-success btn-lg btn-block" type="button">Start Donations</button>
                  </div>
                </div>
                <div className="row" style={{display: "none"}}>
                  <div className="col-xs-12">
                    <p className="payment-errors"></p>
                  </div>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

function stripeResponseHandler(status, response) {
  // Grab the form:
  var $form = $('#payment-form');

  if (response.error) { // Problem!

    // Show the errors on the form:
    $form.find('.payment-errors').text(response.error.message);
    $form.find('.submit').prop('disabled', false); // Re-enable submission

  } else { // Token was created!

    // Get the token ID:
    var token = response.id;

    // Insert the token ID into the form so it gets submitted to the server:
    $form.append($('<input type="hidden" name="stripeToken">').val(token));

    // Submit the form:
    $form.get(0).submit();
  }
};

$(function() {
  var $form = $('#payment-form');
  $form.submit(function(event) {
    // Disable the submit button to prevent repeated clicks:
    $form.find('.submit').prop('disabled', true);

    // Request a token from Stripe:
    global.Stripe.card.createToken($form, stripeResponseHandler);

    // Prevent the form from being submitted:
    return false;
  });
});

export default Donate;

import React, { Component } from 'react';
import { Link } from 'react-router';

class Home extends Component {
  render() {
    return (
      <div>
        <nav id="mainNav" className="navbar navbar-default navbar-fixed-top navbar-custom">
          <div className="container">
            <div className="navbar-header page-scroll">
              <button type="button" className="navbar-toggle" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1">
                <span className="sr-only">Toggle navigation</span> Menu <i className="fa fa-bars"></i>
              </button>
              <a className="navbar-brand" href="#page-top">DonorDonkey</a>
            </div>

            <div className="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
              <ul className="nav navbar-nav navbar-right">
                <li className="hidden">
                  <a href="#page-top"></a>
                </li>
                <li className="page-scroll">
                  <a href="#how-it-works">How it Works</a>
                </li>
                <li className="page-scroll">
                  <Link to={`/donate`}>Donate</Link>
                </li>
              </ul>
            </div>
          </div>
        </nav>

        <header>
          <div className="container">
            <div className="row">
              <div className="col-lg-12">
                {/*<img className="img-responsive" src="img/profile.png" alt=""/>*/}
                <div className="intro-text">
                  <span className="name">DonorDonkey</span>
                  <hr className="star-light"/>
                  <span className="skills">Split Your Donations to Multiple Non-Profits</span>
                  <br/>
                  <Link className="btn btn-danger btn-lg" to={`/donate`}>Donate Now</Link>
                </div>
              </div>
            </div>
          </div>
        </header>

        <section id="how-it-works">
          <div className="container">
            <div className="row">
              <div className="col-lg-12 text-center">
                <h2>How it Works</h2>
                <hr className="star-primary"/>
              </div>
            </div>
            <div className="row">
              <div className="col-sm-4 portfolio-item">
                <h3>1. Pick Your Non-Profits</h3>
              </div>

              <div className="col-sm-4 portfolio-item">
                <h3>2. Specify Contribution</h3>
              </div>

              <div className="col-sm-4 portfolio-item">
                <h3>3. Donate</h3>
              </div>

            </div>
          </div>
        </section>

        <footer className="text-center">
          <div className="footer-above">
            <div className="container">
              <div className="row">
                <div className="footer-col col-md-4">
                  <h3>Location</h3>
                  <p>3481 Melrose Place
                    <br/>Beverly Hills, CA 90210</p>
                </div>
                <div className="footer-col col-md-4">
                  <h3>Around the Web</h3>
                  <ul className="list-inline">
                    <li>
                      <a href="#" className="btn-social btn-outline"><i className="fa fa-fw fa-facebook"></i></a>
                    </li>
                    <li>
                      <a href="#" className="btn-social btn-outline"><i className="fa fa-fw fa-github"></i></a>
                    </li>
                  </ul>
                </div>
                <div className="footer-col col-md-4">
                  <h3>About BuildLeft</h3>
                  <p>BuildLeft provides tools for progressive causes.</p>
                </div>
              </div>
            </div>
          </div>
          <div className="footer-below">
            <div className="container">
              <div className="row">
                <div className="col-lg-12">
                  &copy; 2017 <a href="https://www.opszero.com">Acksin, LLC</a>
                </div>
              </div>
            </div>
          </div>
        </footer>
      </div>

    );
  }
}

export default Home;

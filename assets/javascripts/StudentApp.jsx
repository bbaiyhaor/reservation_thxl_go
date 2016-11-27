/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import ReactDOM from 'react-dom';
import { Router, Route, IndexRoute, hashHistory } from 'react-router';

import StudentLoginPage from '#pages/student/StudentLoginPage';
import StudentRegisterPage from '#pages/student/StudentRegisterPage';
import StudentProtocolPage from '#pages/student/StudentProtocolPage';
import StudentReservationListPage from '#pages/student/StudentReservationListPage';
import StudentMakeReservationPage from '#pages/student/StudentMakeReservationPage';
import StudentFeedbackPage from '#pages/student/StudentFeedbackPage';

const routes = (
  <Route path="/">
    <IndexRoute component={StudentReservationListPage} />
    <Route path="login" component={StudentLoginPage} />
    <Route path="register" component={StudentRegisterPage} />
    <Route path="protocol" component={StudentProtocolPage} />
    <Route path="reservation" component={StudentReservationListPage} />
    <Route path="reservation/make" component={StudentMakeReservationPage} />
    <Route path="reservation/feedback" component={StudentFeedbackPage} />
  </Route>
);

ReactDOM.render(
  <Router history={hashHistory}>
    {routes}
  </Router>,
  document.getElementById('content'),
);

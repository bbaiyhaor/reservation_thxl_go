import { BrowserRouter, Route } from 'react-router-dom';
import React from 'react';
import ReactDOM from 'react-dom';
import StudentFeedbackPage from '#pages/student/StudentFeedbackPage';
import StudentLoginPage from '#pages/student/StudentLoginPage';
import StudentMakeReservationPage from '#pages/student/StudentMakeReservationPage';
import StudentMakeReservationSuccessPage from '#pages/student/StudentMakeReservationSuccessPage';
import StudentProtocolPage from '#pages/student/StudentProtocolPage';
import StudentRegisterPage from '#pages/student/StudentRegisterPage';
import StudentReservationListPage from '#pages/student/StudentReservationListPage';

const routes = (
  <div>
    <Route exact path="/" component={StudentLoginPage} />
    <Route exact path="/login" component={StudentLoginPage} />
    <Route exact path="/register" component={StudentRegisterPage} />
    <Route exact path="/protocol" component={StudentProtocolPage} />
    <Route exact path="/reservation" component={StudentReservationListPage} />
    <Route exact path="/reservation/make" component={StudentMakeReservationPage} />
    <Route exact path="/reservation/make/success" component={StudentMakeReservationSuccessPage} />
    <Route exact path="/reservation/feedback" component={StudentFeedbackPage} />
  </div>
);

ReactDOM.render(
  <BrowserRouter basename="/m/student">
    {routes}
  </BrowserRouter>,
  document.getElementById('content'),
);

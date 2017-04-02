import { BrowserRouter, Route } from 'react-router-dom';
import React from 'react';
import ReactDOM from 'react-dom';
import TeacherChangePasswordPage from '#pages/teacher/TeacherChangePasswordPage';
import TeacherFeedbackPage from '#pages/teacher/TeacherFeedbackPage';
import TeacherLoginPage from '#pages/teacher/TeacherLoginPage';
import TeacherReservationListPage from '#pages/teacher/TeacherReservationListPage';
import TeacherResetPasswordPage from '#pages/teacher/TeacherResetPasswordPage';
import TeacherViewStudentInfoPage from '#pages/teacher/TeacherViewStudentInfoPage';

const routes = (
  <div>
    <Route exact path="/" component={TeacherReservationListPage} />
    <Route exact path="/login" component={TeacherLoginPage} />
    <Route exact path="/password/change" component={TeacherChangePasswordPage} />
    <Route exact path="/password/reset" component={TeacherResetPasswordPage} />
    <Route exact path="/reservation" component={TeacherReservationListPage} />
    <Route exact path="/reservation/feedback" component={TeacherFeedbackPage} />
    <Route exact path="/student" component={TeacherViewStudentInfoPage} />
  </div>
);

ReactDOM.render(
  <BrowserRouter basename="/m/teacher">
    {routes}
  </BrowserRouter>,
  document.getElementById('content'),
);

/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import ReactDOM from 'react-dom';
import { Router, Route, IndexRoute, hashHistory } from 'react-router';

import TeacherLoginPage from '#pages/teacher/TeacherLoginPage';
import TeacherChangePasswordPage from '#pages/teacher/TeacherChangePasswordPage';
import TeacherResetPasswordPage from '#pages/teacher/TeacherResetPasswordPage';
import TeacherReservationListPage from '#pages/teacher/TeacherReservationListPage';
import TeacherFeedbackPage from '#pages/teacher/TeacherFeedbackPage';
import TeacherViewStudentInfoPage from '#pages/teacher/TeacherViewStudentInfoPage';

const routes = (
  <Route path="/">
    <IndexRoute component={TeacherReservationListPage} />
    <Route path="login" component={TeacherLoginPage} />
    <Route path="password/change" component={TeacherChangePasswordPage} />
    <Route path="password/reset" component={TeacherResetPasswordPage} />
    <Route path="reservation" component={TeacherReservationListPage} />
    <Route path="reservation/feedback" component={TeacherFeedbackPage} />
    <Route path="student" component={TeacherViewStudentInfoPage} />
  </Route>
);

ReactDOM.render(
  <Router history={hashHistory}>
    {routes}
  </Router>,
  document.getElementById('content'),
);

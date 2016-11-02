/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import ReactDOM from 'react-dom';
import { Router, Route, IndexRoute, hashHistory } from 'react-router';

import TeacherLoginPage from '#pages/teacher/TeacherLoginPage';
import TeacherReservationListPage from '#pages/teacher/TeacherReservationListPage';

const routes = (
  <Route path="/">
    <IndexRoute component={TeacherReservationListPage} />
    <Route path="login" component={TeacherLoginPage} />
    <Route path="reservation" component={TeacherReservationListPage} />
  </Route>
);

ReactDOM.render(
  <Router history={hashHistory}>
    {routes}
  </Router>,
  document.getElementById('content')
);

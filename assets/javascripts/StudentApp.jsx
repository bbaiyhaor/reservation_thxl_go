/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import ReactDOM from 'react-dom';
import {Router, Route, IndexRoute, hashHistory} from 'react-router';

import StudentLoginPage from '#pages/StudentLoginPage';
import StudentRegisterPage from '#pages/StudentRegisterPage';
import StudentProtocolPage from '#pages/StudentProtocolPage';
import StudentReservationListPage from '#pages/StudentReservationListPage';
import StudentMakeReservationPage from '#pages/StudentMakeReservationPage';

class StudentApp extends React.Component{
    render() {
        return (
            <div id="student-app">
                {this.props.children}
            </div>
        );
    }
}

const routes = (
    <Route path="/" component={StudentApp}>
        <IndexRoute component={StudentReservationListPage}/>
        <Route path="login" component={StudentLoginPage}/>
        <Route path="register" component={StudentRegisterPage}/>
        <Route path="protocol" component={StudentProtocolPage}/>
        <Route path="reservation" component={StudentReservationListPage}/>
        <Route path="reservation/make" component={StudentMakeReservationPage}/>
    </Route>
);

ReactDOM.render(
    <Router history={hashHistory}>{routes}</Router>,
    document.getElementById('content')
);
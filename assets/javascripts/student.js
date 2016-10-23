/**
 * Created by shudi on 2016/10/22.
 */
import React from 'react';
import ReactDOM from 'react-dom';
import {Router, Route, IndexRoute, hashHistory} from 'react-router';

import StudentLoginPage from '#pages/student-login';
import StudentRegisterPage from '#pages/student-register';
import StudentProtocolPage from '#pages/student-protocol';
import StudentReservationPage from '#pages/student-reservation';

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
        <IndexRoute component={StudentReservationPage}/>
        <Route path="login" component={StudentLoginPage}/>
        <Route path="register" component={StudentRegisterPage}/>
        <Route path="protocol" component={StudentProtocolPage}/>
        <Route path="reservation" component={StudentReservationPage}/>
    </Route>
);

ReactDOM.render(
    <Router history={hashHistory}>{routes}</Router>,
    document.getElementById('content')
);
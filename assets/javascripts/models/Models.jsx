/**
 * Created by shudi on 2016/10/22.
 */
/* eslint no-use-before-define: "off" */
import $ from 'jquery';

let apiServer = '/api';
if (process.env.NODE_ENV === 'development') {
  apiServer = 'http://localhost:9000/api';
}

const apiStudentLogin = `${apiServer}/user/student/login`;
const apiStudentRegister = `${apiServer}/user/student/register`;
// const apiTeacherLogin = `${apiServer}/user/teacher/login`;
// const apiAdminLogin = `${apiServer}/user/admin/login`;
const apiLogout = `${apiServer}/user/logout`;
const apiViewReservationsByStudent = `${apiServer}/student/reservation/view`;

function fetch(url, method, payload, succCallback, errCallback) {
  $.ajax({
    url,
    type: method,
    dataType: 'json',
    data: payload,
  }).done((data) => {
    if (process.env.NODE_ENV === 'development') {
      console.log(`"${method}" ${url}`);
      console.log(payload);
      console.log(data);
    }
    succCallback && succCallback(data);
  }).fail((xhr, errorType, error) => {
    if (process.env.NODE_ENV === 'development') {
      console.error(`fetch ${url} error:`, errorType, error, xhr);
    }
    errCallback && errCallback('服务器开小差了，请稍候重试！');
  });
}

export const User = {
  userId: '',
  username: '',
  userType: -1,
  fullname: '',
  student: null,

  clearUser() {
    this.userId = '';
    this.username = '';
    this.userType = -1;
    this.fullname = '';
    this.student = null;
  },

  studentLogin(username, password, succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        this.userId = data.payload.user_id;
        this.username = data.payload.username;
        this.userType = data.payload.user_type;
        this.fullname = data.payload.fullname;
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    const payload = {
      username,
      password,
    };
    fetch(apiStudentLogin, 'POST', payload, succ, errCallback);
  },

  studentRegister(username, password, succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        this.userId = data.payload.user_id;
        this.username = data.payload.username;
        this.userType = data.payload.user_type;
        this.fullname = data.payload.fullname;
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    const payload = {
      username,
      password,
    };
    fetch(apiStudentRegister, 'POST', payload, succ, errCallback);
  },

  logout(succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        this.clearUser();
        Application.clearApplication();
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.status, data.payload);
      }
    };
    fetch(apiLogout, 'GET', {}, succ, errCallback);
  },
};

export const Application = {
  reservations: null,

  clearApplication() {
    this.reservations = null;
  },

  ViewReservationsByStudent(succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        this.reservations = data.payload.reservations;
        // for (let i = 0; i < 10; i++) {
        //     this.reservations.push(this.reservations[0]);
        // }
        // this.reservations.push(this.reservations[1]);
        User.student = data.payload.student;
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.status, data.payload);
      }
    };
    fetch(apiViewReservationsByStudent, 'GET', {}, succ, errCallback);
  },
};

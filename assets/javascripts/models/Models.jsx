/**
 * Created by shudi on 2016/10/22.
 */
/* eslint no-use-before-define: "off", no-console: "off" */
import $ from 'jquery';

let apiServer = '/api';
if (process.env.NODE_ENV === 'development') {
  apiServer = 'http://localhost:9000/api';
}

// User API
const apiUpdateSession = `${apiServer}/user/session`;
const apiStudentLogin = `${apiServer}/user/student/login`;
const apiStudentRegister = `${apiServer}/user/student/register`;
const apiTeacherLogin = `${apiServer}/user/teacher/login`;
const apiTeacherPasswordChange = `${apiServer}/user/teacher/password/change`;
// const apiAdminLogin = `${apiServer}/user/admin/login`;
const apiLogout = `${apiServer}/user/logout`;

// Category API
const apiGetFeedbackCategories = `${apiServer}/category/feedback`;

// Student API
const apiViewReservationsByStudent = `${apiServer}/student/reservation/view`;
const apiMakeReservationByStudent = `${apiServer}/student/reservation/make`;
const apiGetFeedbackByStudent = `${apiServer}/student/reservation/feedback/get`;
const apiSubmitFeedbackByStudent = `${apiServer}/student/reservation/feedback/submit`;

// Teacher API
const apiViewReservationsByTeacher = `${apiServer}/teacher/reservation/view`;
const apiGetFeedbackByTeacher = `${apiServer}/teacher/reservation/feedback/get`;
const apiSubmitFeedbackByTeacher = `${apiServer}/teacher/reservation/feedback/submit`;
const apiGetStudentInfoByTeacher = `${apiServer}/teacher/student/get`;
const apiQueryStudentInfoByTeacher = `${apiServer}/teacher/student/query`;

function fetch(url, method, payload, succCallback, errCallback) {
  $.ajax({
    url,
    type: method,
    dataType: 'json',
    data: payload,
    traditional: true,
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
  teacher: null,

  clearUser() {
    this.userId = '';
    this.username = '';
    this.userType = -1;
    this.fullname = '';
    this.student = null;
  },

  updateSession(succCallback, errCallback) {
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
    fetch(apiUpdateSession, 'GET', {}, succ, errCallback);
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

  teacherLogin(username, password, succCallback, errCallback) {
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
    fetch(apiTeacherLogin, 'POST', payload, succ, errCallback);
  },

  teacherPasswordChange(username, oldPassword, newPassword, succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    const payload = {
      username,
      old_password: oldPassword,
      new_password: newPassword,
    };
    fetch(apiTeacherPasswordChange, 'POST', payload, succ, errCallback);
  },

  logout(succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        this.clearUser();
        Application.clearApplication();
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
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

  getFeedbackCategories(succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    fetch(apiGetFeedbackCategories, 'GET', {}, succ, errCallback);
  },

  viewReservationsByStudent(succCallback, errCallback) {
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
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    fetch(apiViewReservationsByStudent, 'GET', {}, succ, errCallback);
  },

  makeReservationByStudent(payload, succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    fetch(apiMakeReservationByStudent, 'POST', payload, succ, errCallback);
  },

  getFeedbackByStudent(reservationId, sourceId, succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    const payload = {
      reservation_id: reservationId,
      source_id: sourceId,
    };
    fetch(apiGetFeedbackByStudent, 'POST', payload, succ, errCallback);
  },

  submitFeedbackByStudent(reservationId, sourceId, scores, succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    const payload = {
      reservation_id: reservationId,
      source_id: sourceId,
      scores,
    };
    fetch(apiSubmitFeedbackByStudent, 'POST', payload, succ, errCallback);
  },

  viewReservationsByTeacher(succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        this.reservations = data.payload.reservations;
        // for (let i = 0; i < 10; i += 1) {
        //   this.reservations.push(this.reservations[0]);
        // }
        // this.reservations.push(this.reservations[1]);
        User.teacher = data.payload.teacher;
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    fetch(apiViewReservationsByTeacher, 'GET', {}, succ, errCallback);
  },

  getFeedbackByTeacher(reservationId, sourceId, succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    const payload = {
      reservation_id: reservationId,
      source_id: sourceId,
    };
    fetch(apiGetFeedbackByTeacher, 'POST', payload, succ, errCallback);
  },

  submitFeedbackByTeacher(payload, succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    fetch(apiSubmitFeedbackByTeacher, 'POST', payload, succ, errCallback);
  },

  getStudentInfoByTeacher(studentId, succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    const payload = {
      student_id: studentId,
    };
    fetch(apiGetStudentInfoByTeacher, 'POST', payload, succ, errCallback);
  },

  queryStudentInfoByTeacher(studentUsername, succCallback, errCallback) {
    const succ = (data) => {
      if (data.status === 'OK') {
        succCallback && succCallback(data.payload);
      } else {
        errCallback && errCallback(data.err_msg, data.payload);
      }
    };
    const payload = {
      student_username: studentUsername,
    };
    fetch(apiQueryStudentInfoByTeacher, 'POST', payload, succ, errCallback);
  },
};

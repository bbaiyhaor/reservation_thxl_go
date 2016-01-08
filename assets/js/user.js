var width = $(window).width();
var height = $(window).height();

function optimize(t) {
  $(t).css('left', (width - $(t).width()) / 2 - 11 + 'px');
  $(t).css('top', (height - $(t).height()) / 2 - 11 + 'px');
}

function studentLogin() {
  $.ajax({
    type: 'POST',
    async: false,
    url: '/user/student/login',
    data: {
      username: $('#username').val(),
      password: $('#password').val(),
    },
    dataType: 'json',
    success: function(data) {
      if (data.state === 'SUCCESS') {
        window.location.href = data.url;
      } else {
        alert(data.message);
      }
    }
  });
}

function studentRegister() {
  if (!$('#agree').is(':checked')) {
    alert('请先阅读并同意咨询协议');
    return;
  }
  var username = $('#username').val();
  var password = $('#password').val();
  var passwordConfirm = $('#password_confirm').val();
  if (password !== passwordConfirm) {
    alert('两次密码不一致，请重新输入');
    $('#password').val('');
    $('#password_confirm').val('');
    return;
  }
  $.ajax({
    type: 'POST',
    async: false,
    url: '/user/student/register',
    data: {
      username: username,
      password: password,
    },
    dataType: 'json',
    success: function(data) {
      if (data.state === 'SUCCESS') {
        window.location.href = data.url;
      } else {
        alert(data.message);
      }
    }
  });
}

function teacherLogin() {
  $.ajax({
    type: 'POST',
    async: false,
    url: '/user/teacher/login',
    data: {
      username: $('#username').val(),
      password: $('#password').val(),
    },
    dataType: 'json',
    success: function(data) {
      if (data.state === 'SUCCESS') {
        window.location.href = data.url;
      } else {
        alert(data.message);
      }
    }
  });
}

function adminLogin() {
  $.ajax({
    type: 'POST',
    async: false,
    url: '/user/admin/login',
    data: {
      username: $('#username').val(),
      password: $('#password').val(),
    },
    dataType: 'json',
    success: function(data) {
      if (data.state === 'SUCCESS') {
        window.location.href = data.url;
      } else {
        alert(data.message);
      }
    }
  });
}

function logout() {
  $.ajax({
    type: 'GET',
    async: false,
    url: '/user/logout',
    data: {},
    dataType: 'json',
    success: function(data) {
      if (data.state === 'SUCCESS') {
        window.location.href = data.url;
      }
    },
  });
}

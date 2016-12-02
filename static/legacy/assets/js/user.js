var width = $(window).width();
var height = $(window).height();

function optimize(t) {
  $(t).css('left', (width - $(t).width()) / 2 - 11 + 'px');
  $(t).css('top', (height - $(t).height()) / 2 - 11 + 'px');
}

function adminLogin() {
  $.ajax({
    type: 'POST',
    async: false,
    url: '/api/user/admin/login',
    data: {
      username: $('#username').val(),
      password: $('#password').val(),
    },
    dataType: 'json',
    success: function(data) {
      if (data.status === 'OK') {
        window.location.href = data.payload.redirect_url;
      } else {
        alert(data.err_msg);
      }
    }
  });
}

function logout() {
  $.ajax({
    type: 'GET',
    async: false,
    url: '/api/user/logout',
    data: {},
    dataType: 'json',
    success: function(data) {
      if (data.status === 'OK') {
        window.location.href = data.payload.redirect_url;
      }
    },
  });
}

var width = $(window).width();
var height = $(window).height();
var timedReservations;

function viewTimedReservations() {
  $.getJSON('/api/admin/timetable/view', function(json, textStatus) {
    if (json.status === 'OK') {
      console.log(json.payload);
      timedReservations = json.payload.timed_reservations;
      for (var weekday in timedReservations) {
        if (timedReservations.hasOwnProperty(weekday)) {
          refreshDataTable(weekday, timedReservations[weekday]);
          optimize(weekday);
        }
      }
    } else {
      alert(json.err_msg);
    }
  });
}

function refreshDataTable(weekday, timedReservations) {
  $('#page_maintable_' + weekday).first().html('\
  <div class="table_col" id="col_select_' + weekday + '">\
    <div class="table_head table_cell" id="head_select_' + weekday + '">\
      <button class="btn_select_all" name="all" onclick="selectAll();">全选</button>\
    </div>\
  </div>\
  <div class="table_col" id="col_time_' + weekday + '">\
    <div class="table_head table_cell">时间</div>\
  </div>\
  <div class="table_col" id="col_teacher_fullname_' + weekday + '">\
    <div class="table_head table_cell">咨询师</div>\
  </div>\
  <div class="table_col" id="col_teacher_username_' + weekday + '">\
    <div class="table_head table_cell">咨询师编号</div>\
  </div>\
  <div class="table_col" id="col_teacher_mobile_' + weekday + '">\
    <div class="table_head table_cell">咨询师手机</div>\
  </div>\
  <div class="table_col" id="col_teacher_work_type_' + weekday + '">\
    <div class="table_head table_cell">咨询师工作类型</div>\
  </div>\
  <div class="table_col" id="col_status_' + weekday + '">\
    <div class="table_head table_cell">状态</div>\
  </div>\
  <div class="table_col" id="col_operation_' + weekday + '">\
    <div class="table_head table_cell">操作</div>\
  </div>\
  <div class="clearfix"></div>\
');

  for (var i = 0; i < timedReservations.length; ++i) {
    $('#col_select_' + weekday).append('<div class="table_cell" id="cell_select_' + weekday + '_' + i + '">' + 
      '<input class="checkbox" type="checkbox" id="cell_checkbox_' + weekday + '_' + i + '"></div>');
    $('#col_time_' + weekday).append('<div class="table_cell" id="cell_time_' + weekday + '_' + i + '" \
      name="' + i + '">' + timedReservations[i].start_clock + '　-　' + timedReservations[i].end_clock + '</div>');
    $('#col_teacher_fullname_' + weekday).append('<div class="table_cell" id="cell_teacher_fullname_' + weekday + '_' + 
      i + '">' + timedReservations[i].teacher_fullname + '</div>');
    $('#col_teacher_username_' + weekday).append('<div class="table_cell" id="cell_teacher_username_' + weekday + '_'
      + i + '">' + timedReservations[i].teacher_username + '</div>');
    $('#col_teacher_mobile_' + weekday).append('<div class="table_cell" id="cell_teacher_mobile_' + weekday + '_'
      + i + '">' + timedReservations[i].teacher_mobile + '</div>');
    $('#col_teacher_work_type_' + weekday).append('<div class="table_cell" id="cell_teacher_work_type_' + weekday + '_'
      + i + '"><select id="select_teacher_work_type_' + weekday + '_' + i + '"><option value="0">请选择</option>' +
      '<option value="1">专职咨询师</option><option value="2">兼职咨询师</option><option value="3">实习咨询师</option>' +
      '</select></div>');
    $('#select_teacher_work_type_' + weekday + '_' + i).val(timedReservations[i].teacher_work_type);
    if (timedReservations[i].status === 1) {
      $('#col_status_' + weekday).append('<div class="table_cell" id="cell_status_' + weekday + '_' + i + '">生效</div>');
      $('#col_operation_' + weekday).append('<div class="table_cell" id="cell_operation_' + weekday + '_' + + i + '">'
        + '<button type="button" id="cell_operation_oper_' + weekday + '_' + + i + '" onclick="closeTimedReservation(\''
        + weekday + '\',' + i + ');">关闭</button></div>');
    } else {
      $('#col_status_' + weekday).append('<div class="table_cell" id="cell_status_' + weekday + '_' + i + '">关闭</div>');
      $('#col_operation_' + weekday).append('<div class="table_cell" id="cell_operation_' + weekday + '_' + + i + '">'
        + '<button type="button" id="cell_operation_oper_' + weekday + '_' + + i + '" onclick="openTimedReservation(\''
        + weekday + '\',' + i + ');">打开</button></div>');
    }
  }
  $('#col_select_' + weekday).append('<div class="table_cell" id="cell_select_' + weekday + '_add"><input type="checkbox"></div>');
  $('#col_time_' + weekday).append('<div class="table_cell" id="cell_time_' + weekday + '_add">点击新增</div>');
  $('#col_teacher_fullname_' + weekday).append('<div class="table_cell" id="cell_teacher_fullname_' + weekday + '_add"></div>');
  $('#col_teacher_username_' + weekday).append('<div class="table_cell" id="cell_teacher_username_' + weekday + '_add"></div>');
  $('#col_teacher_mobile_' + weekday).append('<div class="table_cell" id="cell_teacher_mobile_' + weekday + '_add"></div>');
  $('#col_teacher_work_type_' + weekday).append('<div class="table_cell" id="cell_teacher_work_type_' + weekday + '_add"></div>');
  $('#col_status_' + weekday).append('<div class="table_cell" id="cell_status_' + weekday + '_add"></div>');
  $('#col_operation_' + weekday).append('<div class="table_cell" id="cell_operation_' + weekday + '_add"></div>');
  $(function() {
    for (var i = 0; i < timedReservations.length; i++) {
      $('#cell_time_' + weekday + '_' + i).click(function(e) {
        editTimedReservation(weekday, $(e.target).attr('name'));
      });
    }
    $('#cell_time_' + weekday + '_add').click(function() {
      addTimedReservation(weekday);
    });
  });
}

function selectAll() {
  if ($('.btn_select_all').first() && $('.btn_select_all').first().prop('name') && $('.btn_select_all').first().prop('name') === 'all') {
    $('.checkbox').prop('checked', true);
    $('.btn_select_all').prop('name', 'none');
    $('.btn_select_all').text('不选');
  } else {
    $('.checkbox').prop('checked', false);
    $('.btn_select_all').prop('name', 'all');
    $('.btn_select_all').text('全选');
  }
}

function optimize(weekday, t) {
  $('#col_select_' + weekday).width(40);
  $('#col_time_' + weekday).width(300);
  $('#col_teacher_fullname_' + weekday).width(120);
  $('#col_teacher_username_' + weekday).width(160);
  $('#col_teacher_mobile_' + weekday).width(160);
  $('#col_teacher_work_type_' + weekday).width(160);
  $('#col_status_' + weekday).width(85);
  $('#col_operation_' + weekday).width(85);
  for (var i = 0; i < timedReservations[weekday].length; ++i) {
    var maxHeight = Math.max(
      $('#cell_select_' + weekday + '_' + i).height(),
      $('#cell_time_' + weekday + '_' + i).height(),
      $('#cell_teacher_fullname_' + weekday + '_' + i).height(),
      $('#cell_teacher_username_' + weekday + '_' + i).height(),
      $('#cell_teacher_mobile_' + weekday + '_' + i).height(),
      $('#cell_teacher_work_type_' + weekday + '_' + i).height(),
      $('#cell_status_' + weekday + '_' + i).height(),
      $('#cell_operation_' + weekday + '_' + i).height()
    );
    $('#cell_select_' + weekday + '_' + i).height(maxHeight);
    $('#cell_time_' + weekday + '_' + i).height(maxHeight);
    $('#cell_teacher_fullname_' + weekday + '_' + i).height(maxHeight);
    $('#cell_teacher_username_' + weekday + '_' + i).height(maxHeight);
    $('#cell_teacher_mobile_' + weekday + '_' + i).height(maxHeight);
    $('#cell_teacher_work_type_' + weekday + '_' + i).height(maxHeight);
    $('#cell_status_' + weekday + '_' + i).height(maxHeight);
    $('#cell_operation_' + weekday + '_' + i).height(maxHeight);

    if (i % 2 == 1) {
      $('#cell_select_' + weekday + '_' + i).css('background-color', 'white');
      $('#cell_time_' + weekday + '_' + i).css('background-color', 'white');
      $('#cell_teacher_fullname_' + weekday + '_' + i).css('background-color', 'white');
      $('#cell_teacher_username_' + weekday + '_' + i).css('background-color', 'white');
      $('#cell_teacher_mobile_' + weekday + '_' + i).css('background-color', 'white');
      $('#cell_teacher_work_type_' + weekday + '_' + i).css('background-color', 'white');
      $('#cell_status_' + weekday + '_' + i).css('background-color', 'white');
      $('#cell_operation_' + weekday + '_' + i).css('background-color', 'white');
    } else {
      $('#cell_select_' + weekday + '_' + i).css('background-color', '#f3f3ff');
      $('#cell_time_' + weekday + '_' + i).css('background-color', '#f3f3ff');
      $('#cell_teacher_fullname_' + weekday + '_' + i).css('background-color', '#f3f3ff');
      $('#cell_teacher_username_' + weekday + '_' + i).css('background-color', '#f3f3ff');
      $('#cell_teacher_mobile_' + weekday + '_' + i).css('background-color', '#f3f3ff');
      $('#cell_teacher_work_type_' + weekday + '_' + i).css('background-color', '#f3f3ff');
      $('#cell_status_' + weekday + '_' + i).css('background-color', '#f3f3ff');
      $('#cell_operation_' + weekday + '_' + i).css('background-color', '#f3f3ff');
    }

    if (timedReservations[weekday][i].status === 1) {
      $('#cell_status_' + weekday + '_' + i).css('background-color', 'greenyellow');
    } else {
      $('#cell_status_' + weekday + '_' + i).css('background-color', 'red');
    }
  }
  $('#cell_select_' + weekday + '_add').height(28);
  $('#cell_time_' + weekday + '_add').height(28);
  $('#cell_teacher_fullname_' + weekday + '_add').height(28);
  $('#cell_teacher_username_' + weekday + '_add').height(28);
  $('#cell_teacher_mobile_' + weekday + '_add').height(28);
  $('#cell_teacher_work_type_' + weekday + '_add').height(28);
  $('#cell_status_' + weekday + '_add').height(28);
  $('#cell_operation_' + weekday + '_add').height(28);

  $('.table_head').height($('#head_select_' + weekday).height());
  $(t).css('left', (width - $(t).width()) / 2 - 11 + 'px');
  $(t).css('top', (height - $(t).height()) / 2 - 11 + 'px');
  $('#page_maintable_' + weekday).css('margin-left', 0.5 * (width - (40 + 300 + 120 + 85 + 85 + 480)) + 'px');
}

function addTimedReservation(weekday) {
  $('#cell_time_' + weekday + '_add').off();
  $('#cell_time_' + weekday + '_add').first().html('<input style="width:20px" id="start_hour_' + weekday + '_add"/>时' +
    '<input style="width:20px" id="start_minute_' + weekday + '_add"/>分　-　' +
    '<input style="width:20px" id="end_hour_' + weekday + '_add"/>时' +
    '<input style="width:20px" id="end_minute_' + weekday + '_add"/>分');
  $('#cell_teacher_fullname_' + weekday + '_add').first().html('<input id="teacher_fullname_' + weekday + '_add" style="width:60px"/>'
    + '<button type="button" onclick="searchTeacher(\'' + weekday + '\');">搜索</button>');
  $('#cell_teacher_username_' + weekday + '_add').first().html('<input id="teacher_username_' + weekday + '_add" style="width:120px"/>');
  $('#cell_teacher_mobile_' + weekday + '_add').first().html('<input id="teacher_mobile_' + weekday + '_add" style="width:120px"/>');
  $('#cell_teacher_work_type_' + weekday + '_add').first().html('<select id="select_teacher_work_type_' + weekday + '_add">' +
    '<option value="0">请选择</option><option value="1">专职咨询师</option><option value="2">兼职咨询师</option><option value="3">实习咨询师</option></select>');
  $('#cell_status_' + weekday + '_add').first().html('<button type="button" onclick="addTimedReservationConfirm(\'' + weekday + '\');">确认</button>');
  $('#cell_operation_' + weekday + '_add').first().html('<button type="button" onclick="window.location.reload();">取消</button>');
  optimize(weekday);
}

function addTimedReservationConfirm(weekday) {
  var startHour = $('#start_hour_' + weekday + '_add').val();
  var startMinute = $('#start_minute_' + weekday + '_add').val();
  var startClock = (startHour.length < 2 ? '0' : '') + startHour + ':';
  if (startMinute.length == 0) {
    startClock += '00';
  } else if (startMinute.length == 1) {
    startClock += '0' + startMinute;
  } else {
    startClock += startMinute;
  }
  var endHour = $('#end_hour_' + weekday + '_add').val();
  var endMinute = $('#end_minute_' + weekday + '_add').val();
  var endClock = (endHour.length < 2 ? '0' : '') + endHour + ':';
  if (endMinute.length == 0) {
    endClock += '00';
  } else if (endMinute.length == 1) {
    endClock += '0' + endMinute;
  } else {
    endClock += endMinute;
  }
  var payload = {
    weekday: weekday,
    start_clock: startClock,
    end_clock: endClock,
    teacher_username: $('#teacher_username_' + weekday + '_add').val(),
    teacher_fullname: $('#teacher_fullname_' + weekday + '_add').val(),
    teacher_mobile: $('#teacher_mobile_' + weekday + '_add').val(),
    teacher_work_type: $('#select_teacher_work_type_' + weekday + '_add').val(),
  };
  $.post('/api/admin/timetable/add', payload, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      viewTimedReservations();
    } else if (data.status === '%CHECK%') {
      addTimedReservationCheck(weekday, payload);
    } else {
      alert(data.err_msg);
    }
  });
}

function addTimedReservationCheck(weekday, payload) {
  $('body').append('\
    <div id="pop_add_timed_reservation_check_' + weekday + '" class="pop_window" style="width: 50%">\
      咨询师信息有变更，是否更新？\
      <br>\
      <button type="button" name="confirm">确认</button>\
      <button type="button" onclick="$(\'#pop_add_timed_reservation_check_' + weekday + '\').remove();">取消</button>\
    </div>\
  ');
  $(function() {
    $('#pop_add_timed_reservation_check_' + weekday + ' [name=confirm]').click(function(event) {
      $('#pop_add_timed_reservation_check_' + weekday).remove();
      addTimedReservationCheckConfirm(payload);
    });
  });
  optimize('Monday', '#pop_add_timed_reservation_check_' + weekday);
}

function addTimedReservationCheckConfirm(payload) {
  payload['force'] = true;
  $.post('/api/admin/timetable/add', payload, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      viewTimedReservations();
    } else if (data.status === '%CHECK%') {
      addTimedReservationCheck(payload);
    } else {
      alert(data.err_msg);
    }
  });
}

function editTimedReservation(weekday, index) {
  $('#cell_time_' + weekday + '_' + index).off();
  $('#cell_time_' + weekday + '_' + index).first().html('<input style="width:20px" id="start_hour_' + weekday + '_' + index + '"/>时' +
    '<input style="width:20px" id="start_minute_' + weekday + '_' + index + '"/>分　-　' +
    '<input style="width:20px" id="end_hour_' + weekday + '_' + index + '"/>时' +
    '<input style="width:20px" id="end_minute_' + weekday + '_' + index + '"/>分');
  $('#cell_teacher_fullname_' + weekday + '_' + index).first().html('<input id="teacher_fullname_' + weekday + '_' + index + '" style="width:60px" '
    + 'value="' + timedReservations[weekday][index].teacher_fullname + '""></input>'
    + '<button type="button" onclick="searchTeacher(\'' + weekday + '\',' + index + ');">搜索</button>');
  $('#cell_teacher_username_' + weekday + '_' + index).first().html('<input id="teacher_username_' + weekday + '_' + index + '" style="width:120px" '
    + 'value="' + timedReservations[weekday][index].teacher_username + '"/>');
  $('#cell_teacher_mobile_' + weekday + '_' + index).first().html('<input id="teacher_mobile_' + weekday + '_' + index + '" style="width:120px" '
    + 'value="' + timedReservations[weekday][index].teacher_mobile + '"/>');
  $('#cell_teacher_work_type_' + weekday + '_' + index).first().html('<select id="select_teacher_work_type_' + weekday + '_' + index + '">' +
    '<option value="0">请选择</option><option value="1">专职咨询师</option><option value="2">兼职咨询师</option><option value="3">实习咨询师</option></select>');
  $('#select_teacher_work_type_' + weekday + '_' + index).val(timedReservations[weekday][index].teacher_work_type);
  $('#cell_status_' + weekday + '_' + index).first().html('<button type="button" onclick="editTimedReservationConfirm(\'' + weekday + '\',' + index + ');">确认</button>');
  $('#cell_operation_' + weekday + '_' + index).first().html('<button type="button" onclick="window.location.reload();">取消</button>');
  optimize(weekday);
}

function editTimedReservationConfirm(weekday, index) {
  var startHour = $('#start_hour_' + weekday + '_' + index).val();
  var startMinute = $('#start_minute_' + weekday + '_' + index).val();
  var startClock = (startHour.length < 2 ? '0' : '') + startHour + ':';
  if (startMinute.length == 0) {
      startClock += '00';
  } else if (startMinute.length == 1) {
      startClock += '0' + startMinute;
  } else {
      startClock += startMinute;
  }
  var endHour = $('#end_hour_' + weekday + '_' + index).val();
  var endMinute = $('#end_minute_' + weekday + '_' + index).val();
  var endClock = (endHour.length < 2 ? '0' : '') + endHour + ':';
  if (endMinute.length == 0) {
      endClock += '00';
  } else if (endMinute.length == 1) {
      endClock += '0' + endMinute;
  } else {
      endClock += endMinute;
  }
  var payload = {
      timed_reservation_id: timedReservations[weekday][index].id,
      weekday: weekday,
      start_clock: startClock,
      end_clock: endClock,
      teacher_username: $('#teacher_username_' + weekday + '_' + index).val(),
      teacher_fullname: $('#teacher_fullname_' + weekday + '_' + index).val(),
      teacher_mobile: $('#teacher_mobile_' + weekday + '_' + index).val(),
      teacher_work_type: $('#select_teacher_work_type_' + weekday + '_' + index).val(),
  };
  $.post('/api/admin/timetable/edit', payload, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      viewTimedReservations();
    } else if (data.status === '%CHECK%') {
      editTimedReservationCheck(weekday, payload);
    } else {
      alert(data.err_msg);
    }
  });
}

function editTimedReservationCheck(weekday, payload) {
  $('body').append('\
    <div id="pop_edit_timed_reservation_check_' + weekday + '" class="pop_window" style="width: 50%">\
      咨询师信息有变更，是否更新？\
      <br>\
      <button type="button" name="confirm">确认</button>\
      <button type="button" onclick="$(\'#pop_edit_timed_reservation_check_' + weekday + '\').remove();">取消</button>\
    </div>\
  ');
  $(function() {
    $('#pop_edit_timed_reservation_check_' + weekday + ' [name=confirm]').click(function(event) {
      $('#pop_edit_timed_reservation_check_' + weekday).remove();
      editTimedReservationCheckConfirm(payload);
    });
  });
  optimize('Monday', '#pop_edit_timed_reservation_check_' + weekday);
}

function editTimedReservationCheckConfirm(payload) {
  payload['force'] = true;
  $.post('/api/admin/timetable/edit', payload, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      viewTimedReservations();
    } else if (data.status === '%CHECK%') {
      editTimedReservationCheck(payload);
    } else {
      alert(data.err_msg);
    }
  });
}

function searchTeacher(weekday, index) {
  $.post('/api/admin/teacher/search', {
    teacher_username: $('#teacher_username_' + weekday + '_' + (index === undefined ? 'add' : index)).val(),
    teacher_fullname: $('#teacher_fullname_' + weekday + '_' + (index === undefined ? 'add' : index)).val(),
    teacher_mobile: $('#teacher_mobile_' + weekday + '_' + (index === undefined ? 'add' : index)).val(),
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      $('#teacher_username_' + weekday + '_' + (index === undefined ? 'add' : index)).val(data.payload.teacher.username);
      $('#teacher_fullname_' + weekday + '_' + (index === undefined ? 'add' : index)).val(data.payload.teacher.fullname);
      $('#teacher_mobile_' + weekday + '_' + (index === undefined ? 'add' : index)).val(data.payload.teacher.mobile);
    }
  });
}

function removeTimedReservations() {
  $('body').append('\
  <div id="pop_remove_timed_reservations" class="pop_window" style="width: 50%">\
    确认删除选中的预设咨询？\
    <br>\
    <button type="button" onclick="$(\'#pop_remove_timed_reservations\').remove();removeReservationsConfirm();">确认</button>\
    <button type="button" onclick="$(\'#pop_remove_timed_reservations\').remove();">取消</button>\
  </div>\
');
  optimize('Monday', '#pop_remove_timed_reservations');
}

function removeReservationsConfirm() {
  var timedReservationIds = [];
  for (var weekday in timedReservations) {
    if (timedReservations.hasOwnProperty(weekday)) {
      for (var i = 0; i < timedReservations[weekday].length; ++i) {
        if ($('#cell_checkbox_' + weekday + '_' + i)[0].checked) {
          timedReservationIds.push(timedReservations[weekday][i].id);
        }
      }
    }
  }
  $.ajax({
    url: '/api/admin/timetable/remove',
    type: "POST",
    dataType: "json",
    data: {
      timed_reservation_ids: timedReservationIds,
    },
    traditional: true,
  })
  .done(function(data) {
    if (data.status === 'OK') {
      viewTimedReservations();
    } else {
      alert(data.err_msg);
    }
  });
}

function openTimedReservations() {
  $('body').append('\
    <div id="pop_open_timed_reservations" class="pop_window" style="width: 50%">\
      确认打开选中的预设咨询？\
      <br>\
      <button type="button" onclick="$(\'#pop_open_timed_reservations\').remove();openReservationsConfirm();">确认</button>\
      <button type="button" onclick="$(\'#pop_open_timed_reservations\').remove();">取消</button>\
    </div>\
  ');
  optimize('Monday', '#pop_open_timed_reservations');
}

function openReservationsConfirm() {
  var timedReservationIds = [];
  for (var weekday in timedReservations) {
    if (timedReservations.hasOwnProperty(weekday)) {
      for (var i = 0; i < timedReservations[weekday].length; ++i) {
        if ($('#cell_checkbox_' + weekday + '_' + i)[0].checked) {
          timedReservationIds.push(timedReservations[weekday][i].id);
        }
      }
    }
  }
  $.ajax({
    url: '/api/admin/timetable/open',
    type: "POST",
    dataType: "json",
    data: {
      timed_reservation_ids: timedReservationIds,
    },
    traditional: true,
  })
  .done(function(data) {
    if (data.status === 'OK') {
      viewTimedReservations();
    } else {
      alert(data.err_msg);
    }
  });
}

function closeTimedReservations() {
  $('body').append('\
    <div id="pop_close_timed_reservations" class="pop_window" style="width: 50%">\
      确认关闭选中的预设咨询？\
      <br>\
      <button type="button" onclick="$(\'#pop_close_timed_reservations\').remove();closeReservationsConfirm();">确认</button>\
      <button type="button" onclick="$(\'#pop_close_timed_reservations\').remove();">取消</button>\
    </div>\
  ');
  optimize('Monday', '#pop_close_timed_reservations');
}

function closeReservationsConfirm() {
  var timedReservationIds = [];
  for (var weekday in timedReservations) {
    if (timedReservations.hasOwnProperty(weekday)) {
      for (var i = 0; i < timedReservations[weekday].length; ++i) {
        if ($('#cell_checkbox_' + weekday + '_' + i)[0].checked) {
          timedReservationIds.push(timedReservations[weekday][i].id);
        }
      }
    }
  }
  $.ajax({
    url: '/api/admin/timetable/close',
    type: "POST",
    dataType: "json",
    data: {
      timed_reservation_ids: timedReservationIds,
    },
    traditional: true,
  })
  .done(function(data) {
    if (data.status === 'OK') {
      viewTimedReservations();
    } else {
      alert(data.err_msg);
    }
  });
}

function openTimedReservation(weekday, index) {
  var timedReservationIds = [];
  timedReservationIds.push(timedReservations[weekday][index].id);
  $.ajax({
    url: '/api/admin/timetable/open',
    type: "POST",
    dataType: "json",
    data: {
      timed_reservation_ids: timedReservationIds,
    },
    traditional: true,
  })
  .done(function(data) {
    if (data.status === 'OK') {
      viewTimedReservations();
    } else {
      alert(data.err_msg);
    }
  });
}

function closeTimedReservation(weekday, index) {
  var timedReservationIds = [];
  timedReservationIds.push(timedReservations[weekday][index].id);
  $.ajax({
    url: '/api/admin/timetable/close',
    type: "POST",
    dataType: "json",
    data: {
      timed_reservation_ids: timedReservationIds,
    },
    traditional: true,
  })
  .done(function(data) {
    if (data.status === 'OK') {
      viewTimedReservations();
    } else {
      alert(data.err_msg);
    }
  })
}
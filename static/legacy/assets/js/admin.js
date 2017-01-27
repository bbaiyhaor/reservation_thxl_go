var width = $(window).width();
var height = $(window).height();
var reservations;
var firstCategory;
var secondCategory;

function viewReservations() {
  getFeedbackCategories();
  if ($('#query_date').val() !== '') {
    queryReservations();
    return;
  }
  $.getJSON('/api/admin/reservation/view', function(json, textStatus) {
    if (json.status === 'OK') {
      console.log(json.payload);
      reservations = json.payload.reservations;
      refreshDataTable(reservations);
      optimize();
    } else {
      alert(json.err_msg);
    }
  });
}

function queryReservations() {
  $.getJSON('/api/admin/reservation/view/daily', {
    from_date: $('#query_date').val()
  }, function(json, textStatus) {
    if (json.status === 'OK') {
      console.log(json.payload);
      reservations = json.payload.reservations;
      refreshDataTable(reservations);
      optimize();
    } else {
      alert(json.err_msg);
    }
  });
}

function queryReservationsByTeacher() {
  $.getJSON('/api/admin/reservation/view/teacher/username', {
    teacher_username: $('#query_teacher_username').val()
  }, function(json, textStatus) {
    if (json.status === 'OK') {
      console.log(json.payload);
      reservations = json.payload.reservations;
      refreshDataTable(reservations);
      optimize();
    } else {
      alert(json.err_msg);
    }
  });
}

function exportTodayReservations() {
  $.getJSON('/api/admin/reservation/export/today', function(json, textStatus) {
    if (json.status === 'OK') {
      window.open(json.payload.url);
    } else {
      alert(json.err_msg);
    }
  });
}

function queryStudent() {
  $.post('/api/admin/student/query', {
    student_username: $('#query_student').val()
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      showStudent(data.payload.student, data.payload.reservations);
    } else {
      alert(data.err_msg);
    }
  });
}

function getFeedbackCategories() {
  $.getJSON('/api/category/feedback', function(json, textStatus) {
    if (json.status === 'OK') {
      firstCategory = json.payload.first_category;
      secondCategory = json.payload.second_category;
    }
  });
}

function refreshDataTable(reservations) {
  $('#page_maintable').first().html('\
    <div class="table_col" id="col_select">\
      <div class="table_head table_cell" id="head_select">\
        <button id="btn_select_all" name="all" onclick="selectAll();">全选</button>\
      </div>\
    </div>\
    <div class="table_col" id="col_time">\
      <div class="table_head table_cell">时间</div>\
    </div>\
    <div class="table_col" id="col_teacher_fullname">\
      <div class="table_head table_cell">咨询师</div>\
    </div>\
    <div class="table_col" id="col_teacher_username">\
      <div class="table_head table_cell">咨询师编号</div>\
    </div>\
    <div class="table_col" id="col_teacher_mobile">\
      <div class="table_head table_cell">咨询师手机</div>\
    </div>\
    <div class="table_col" id="col_status">\
      <div class="table_head table_cell">状态</div>\
    </div>\
    <div class="table_col" id="col_student">\
      <div class="table_head table_cell">学生</div>\
    </div>\
    <div class="clearfix"></div>\
  ');

  for (var i = 0; i < reservations.length; ++i) {
    $('#col_select').append('<div class="table_cell" id="cell_select_' + i + '">'
      + '<input class="checkbox" type="checkbox" id="cell_checkbox_' + i + '"></div>');
    $('#col_time').append('<div class="table_cell" id="cell_time_' + i + '" name="' + i + '">' + reservations[i].start_time +
      '至' + reservations[i].end_time + '</div>');
    $('#col_teacher_fullname').append('<div class="table_cell" id="cell_teacher_fullname_'
      + i + '">' + reservations[i].teacher_fullname + '</div>');
    $('#col_teacher_username').append('<div class="table_cell" id="cell_teacher_username_'
      + i + '">' + reservations[i].teacher_username + '</div>');
    $('#col_teacher_mobile').append('<div class="table_cell" id="cell_teacher_mobile_'
      + i + '">' + reservations[i].teacher_mobile + '</div>');
    if (reservations[i].status === 1) {
      $('#col_status').append('<div class="table_cell" id="cell_status_' + i + '">未预约</div>');
      $('#col_student').append('<div class="table_cell" id="cell_student_' + i + '">'
        + '<button type="button" id="cell_student_view_' + i + '" onclick="setStudent(' + i + ')">指定'
        + '</button></div>');
    } else if (reservations[i].status === 2) {
      $('#col_status').append('<div class="table_cell" id="cell_status_' + i + '">已预约</div>');
      $('#col_student').append('<div class="table_cell" id="cell_student_' + i + '">'
        + '<button type="button" id="cell_student_view_' + i + '" onclick="getStudent(' + i + ');">查看'
        + '</button></div>');
    } else if (reservations[i].status === 3) {
      $('#col_status').append('<div class="table_cell" id="cell_status_' + i + '">'
        + '<button type="button" id="cell_status_feedback_' + i + '" onclick="getFeedback(' + i + ');">'
        + '反馈</button></div>');
      $('#col_student').append('<div class="table_cell" id="cell_student_' + i + '">'
        + '<button type="button" id="cell_student_view_' + i + '" onclick="getStudent(' + i + ');">查看'
        + '</button></div>');
    }
  }
  $('#col_select').append('<div class="table_cell" id="cell_select_add"><input type="checkbox"></div>');
  $('#col_time').append('<div class="table_cell" id="cell_time_add">点击新增</div>');
  $('#col_teacher_fullname').append('<div class="table_cell" id="cell_teacher_fullname_add"></div>');
  $('#col_teacher_username').append('<div class="table_cell" id="cell_teacher_username_add"></div>');
  $('#col_teacher_mobile').append('<div class="table_cell" id="cell_teacher_mobile_add"></div>');
  $('#col_status').append('<div class="table_cell" id="cell_status_add"></div>');
  $('#col_student').append('<div class="table_cell" id="cell_student_add"></div>');
  $(function() {
    for (var i = 0; i < reservations.length; i++) {
      $('#cell_time_' + i).click(function(e) {
        editReservation($(e.target).attr('name'));
      });
    }
    $('#cell_time_add').click(addReservation);
  });
}

function selectAll() {
  if ($('#btn_select_all').prop('name') && $('#btn_select_all').prop('name') === 'all') {
    $('.checkbox').prop('checked', true);
    $('#btn_select_all').prop('name', 'none');
    $('#btn_select_all').text('不选');
  } else {
    $('.checkbox').prop('checked', false);
    $('#btn_select_all').prop('name', 'all');
    $('#btn_select_all').text('全选');
  }
}

function optimize(t) {
  $('#col_select').width(40);
  $('#col_time').width(405);
  $('#col_teacher_fullname').width(120);
  $('#col_teacher_username').width(160);
  $('#col_teacher_mobile').width(160);
  $('#col_status').width(85);
  $('#col_student').width(85);
  for (var i = 0; i < reservations.length; ++i) {
    var maxHeight = Math.max(
      $('#cell_select_' + i).height(),
      $('#cell_time_' + i).height(),
      $('#cell_teacher_fullname_' + i).height(),
      $('#cell_teacher_username_' + i).height(),
      $('#cell_teacher_mobile_' + i).height(),
      $('#cell_status_' + i).height(),
      $('#cell_student_' + i).height()
    );
    $('#cell_select_' + i).height(maxHeight);
    $('#cell_time_' + i).height(maxHeight);
    $('#cell_teacher_fullname_' + i).height(maxHeight);
    $('#cell_teacher_username_' + i).height(maxHeight);
    $('#cell_teacher_mobile_' + i).height(maxHeight);
    $('#cell_status_' + i).height(maxHeight);
    $('#cell_student_' + i).height(maxHeight);

    if (i % 2 == 1) {
      $('#cell_select_' + i).css('background-color', 'white');
      $('#cell_time_' + i).css('background-color', 'white');
      $('#cell_teacher_fullname_' + i).css('background-color', 'white');
      $('#cell_teacher_username_' + i).css('background-color', 'white');
      $('#cell_teacher_mobile_' + i).css('background-color', 'white');
      $('#cell_status_' + i).css('background-color', 'white');
      $('#cell_student_' + i).css('background-color', 'white');
    } else {
      $('#cell_select_' + i).css('background-color', '#f3f3ff');
      $('#cell_time_' + i).css('background-color', '#f3f3ff');
      $('#cell_teacher_fullname_' + i).css('background-color', '#f3f3ff');
      $('#cell_teacher_username_' + i).css('background-color', '#f3f3ff');
      $('#cell_teacher_mobile_' + i).css('background-color', '#f3f3ff');
      $('#cell_status_' + i).css('background-color', '#f3f3ff');
      $('#cell_student_' + i).css('background-color', '#f3f3ff');
    }

    if (reservations[i].student_crisis_level && reservations[i].student_crisis_level !== 0) {
      $('#cell_student_' + i).css('background-color', 'rgba(255, 0, 0, ' + parseInt(reservations[i].student_crisis_level) / 1 +')');
    }
  }
  $('#cell_select_add').height(28);
  $('#cell_time_add').height(28);
  $('#cell_teacher_fullname_add').height(28);
  $('#cell_teacher_username_add').height(28);
  $('#cell_teacher_mobile_add').height(28);
  $('#cell_status_add').height(28);
  $('#cell_student_add').height(28);

  $('.table_head').height($('#head_select').height());
  $(t).css('left', (width - $(t).width()) / 2 - 11 + 'px');
  $(t).css('top', (height - $(t).height()) / 2 - 11 + 'px');
  $('#page_maintable').css('margin-left', 0.5 * ($(window).width()
    - (40 + 405 + 120 + 85 + 85 + 320)) + 'px');
}

function addReservation() {
  $('#cell_time_add').off();
  $('#cell_time_add').first().html('<input type="text" id="input_date_add" style="width: 80px"/>日，'
    + '<input style="width:20px" id="start_hour_add"/>时<input style="width:20px" id="start_minute_add"/>分'
    + '至<input style="width:20px" id="end_hour_add"/>时<input style="width:20px" id="end_minute_add"/>分');
  $('#cell_teacher_fullname_add').first().html('<input id="teacher_fullname_add" style="width:60px"/>'
    + '<button type="button" onclick="searchTeacher();">搜索</button>');
  $('#cell_teacher_username_add').first().html('<input id="teacher_username_add" style="width:120px"/>');
  $('#cell_teacher_mobile_add').first().html('<input id="teacher_mobile_add" style="width:120px"/>');
  $('#cell_status_add').first().html('<button type="button" onclick="addReservationConfirm();">确认</button>');
  $('#cell_student_add').first().html('<button type="button" onclick="window.location.reload();">取消</button>');
  $('#input_date_add').datepicker({
    showOtherMonths: true,
    selectOtherMonths: true,
    showButtonPanel: true,
    dateFormat: 'yy-mm-dd',
    showWeek: true,
    firstDay: 1
  });
  optimize();
}

function addReservationConfirm() {
  var startHour = $('#start_hour_add').val();
  var startMinute = $('#start_minute_add').val();
  var endHour = $('#end_hour_add').val();
  var endMinute = $('#end_minute_add').val();
  var startTime = $('#input_date_add').val() + ' ' + (startHour.length < 2 ? '0' : '') + startHour + ':';
  if (startMinute.length == 0) {
    startTime += '00';
  } else if (startMinute.length == 1) {
    startTime += '0' + startMinute;
  } else {
    startTime += startMinute;
  }
  var endTime = $('#input_date_add').val() + ' ' + (endHour.length < 2 ? '0' : '') + endHour + ':';
  if (endMinute.length == 0) {
    endTime += '00';
  } else if (endMinute.length == 1) {
    endTime += '0' + endMinute;
  } else {
    endTime += endMinute;
  }
  var payload = {
    start_time: startTime,
    end_time: endTime,
    teacher_username: $('#teacher_username_add').val(),
    teacher_fullname: $('#teacher_fullname_add').val(),
    teacher_mobile: $('#teacher_mobile_add').val(),
  };
  $.post('/api/admin/reservation/add', payload, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      viewReservations();
    } else if (data.status === '%CHECK%') {
      addReservationCheck(payload);
    } else {
      alert(data.err_msg);
    }
  });
}

function addReservationCheck(payload) {
  $('body').append('\
    <div id="pop_add_reservation_check" class="pop_window" style="width: 50%">\
      咨询师信息有变更，是否更新？\
      <br>\
      <button type="button" name="confirm">确认</button>\
      <button type="button" onclick="$(\'#pop_add_reservation_check\').remove();">取消</button>\
    </div>\
  ');
  $(function() {
    $('#pop_add_reservation_check [name=confirm]').click(function() {
      $('#pop_add_reservation_check').remove();
      addReservationCheckConfirm(payload);
    });
  });
  optimize('#pop_add_reservation_check');
}

function addReservationCheckConfirm(payload) {
  payload['force'] = true;
  $.post('/api/admin/reservation/add', payload, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      viewReservations();
    } else if (data.status === '%CHECK%') {
      addReservationCheck(payload);
    } else {
      alert(data.err_msg);
    }
  });
}

function editReservation(index) {
  $('#cell_time_' + index).off();
  $('#cell_time_' + index).first().html('<input type="text" id="input_date_' + index + '" style="width: 80px"/>日，'
    + '<input style="width:20px" id="start_hour_' + index + '"/>时<input style="width:20px" id="start_minute_' + index + '"/>分'
    + '至<input style="width:20px" id="end_hour_' + index + '"/>时<input style="width:20px" id="end_minute_' + index + '"/>分');
  $('#cell_teacher_fullname_' + index).first().html('<input id="teacher_fullname_' + index + '" style="width:60px" '
    + 'value="' + reservations[index].teacher_fullname + '""></input>'
    + '<button type="button" onclick="searchTeacher(' + index + ');">搜索</button>');
  $('#cell_teacher_username_' + index).first().html('<input id="teacher_username_' + index + '" style="width:120px" '
    + 'value="' + reservations[index].teacher_username + '"/>');
  $('#cell_teacher_mobile_' + index).first().html('<input id="teacher_mobile_' + index + '" style="width:120px" '
    + 'value="' + reservations[index].teacher_mobile + '"/>');
  $('#cell_status_' + index).first().html('<button type="button" onclick="editReservationConfirm(' + index + ');">确认</button>');
  $('#cell_student_' + index).first().html('<button type="button" onclick="window.location.reload();">取消</button>');
  $('#input_date_' + index).datepicker({
    showOtherMonths: true,
    selectOtherMonths: true,
    showButtonPanel: true,
    dateFormat: 'yy-mm-dd',
    showWeek: true,
    firstDay: 1
  });
  optimize();
}

function editReservationConfirm(index) {
  var startHour = $('#start_hour_' + index).val();
  var startMinute = $('#start_minute_' + index).val();
  var endHour = $('#end_hour_' + index).val();
  var endMinute = $('#end_minute_' + index).val();
  var startTime = $('#input_date_' + index).val() + ' ' + (startHour.length < 2 ? '0' : '') + startHour + ':';
  if (startMinute.length == 0) {
    startTime += '00';
  } else if (startMinute.length == 1) {
    startTime += '0' + startMinute;
  } else {
    startTime += startMinute;
  }
  var endTime = $('#input_date_' + index).val() + ' ' + (endHour.length < 2 ? '0' : '') + endHour + ':';
  if (endMinute.length == 0) {
    endTime += '00';
  } else if (endMinute.length == 1) {
    endTime += '0' + endMinute;
  } else {
    endTime += endMinute;
  }
  var payload = {
    reservation_id: reservations[index].id,
    start_time: startTime,
    end_time: endTime,
    teacher_username: $('#teacher_username_' + index).val(),
    teacher_fullname: $('#teacher_fullname_' + index).val(),
    teacher_mobile: $('#teacher_mobile_' + index).val(),
  };
  $.post('/api/admin/reservation/edit', payload, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      viewReservations();
    } else if (data.status === '%CHECK%') {
      editReservationCheck(payload);
    } else {
      alert(data.err_msg);
    }
  });
}

function editReservationCheck(payload) {
  $('body').append('\
    <div id="pop_edit_reservation_check" class="pop_window" style="width: 50%">\
      咨询师信息有变更，是否更新？\
      <br>\
      <button type="button" name="confirm">确认</button>\
      <button type="button" onclick="$(\'#pop_edit_reservation_check\').remove();">取消</button>\
    </div>\
  ');
  $(function() {
    $('#pop_edit_reservation_check [name=confirm]').click(function() {
      $('#pop_edit_reservation_check').remove();
      editReservationCheckConfirm(payload);
    });
  });
  optimize('#pop_edit_reservation_check');
}

function editReservationCheckConfirm(payload) {
  payload['force'] = true;
  $.post('/api/admin/reservation/edit', payload, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      viewReservations();
    } else if (data.status === '%CHECK%') {
      editReservationCheck(payload);
    } else {
      alert(data.err_msg);
    }
  });
}

function searchTeacher(index) {
  $.post('/api/admin/teacher/search', {
    teacher_username: $('#teacher_username_' + (index === undefined ? 'add' : index)).val(),
    teacher_fullname: $('#teacher_fullname_' + (index === undefined ? 'add' : index)).val(),
    teacher_mobile: $('#teacher_mobile_' + (index === undefined ? 'add' : index)).val()
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      $('#teacher_username_' + (index === undefined ? 'add' : index)).val(data.payload.teacher.username);
      $('#teacher_fullname_' + (index === undefined ? 'add' : index)).val(data.payload.teacher.fullname);
      $('#teacher_mobile_' + (index === undefined ? 'add' : index)).val(data.payload.teacher.mobile);
    }
  });
}

function removeReservations() {
  $('body').append('\
    <div id="pop_remove_reservations" class="pop_window" style="width: 50%">\
      确认删除选中的咨询记录？<br>\
      <div id="pop_remove_reservations_div"></div>\
      <button type="button" onclick="$(\'#pop_remove_reservations\').remove();removeReservationsConfirm();">确认</button>\
      <button type="button" onclick="$(\'#pop_remove_reservations\').remove();">取消</button>\
    </div>\
  ');
  for (var i = 0; i < reservations.length; i++) {
    if ($('#cell_checkbox_' + i)[0].checked) {
      $('#pop_remove_reservations_div').append('\
        <span>' + reservations[i].teacher_fullname + '　　　　' + reservations[i].start_time + '至' + reservations[i].end_time + '</span><br>\
      ');
    }
  }
  optimize('#pop_remove_reservations');
}

function removeReservationsConfirm() {
  var reservationIds = [];
  var sourceIds = [];
  var startTimes = [];
  for (var i = 0; i < reservations.length; ++i) {
    if ($('#cell_checkbox_' + i)[0].checked) {
      reservationIds.push(reservations[i].id);
      sourceIds.push(reservations[i].source_id);
      startTimes.push(reservations[i].start_time)
    }
  }
  var payload = {
    reservation_ids: reservationIds,
    source_ids: sourceIds,
    start_times: startTimes,
  };
  $.ajax({
    url: '/api/admin/reservation/remove',
    type: "POST",
    dataType: 'json',
    data: payload,
    traditional: true,
  })
  .done(function(data) {
    if (data.status === 'OK') {
      viewReservations();
    } else {
      alert(data.err_msg);
    }
  });
}

function cancelReservations() {
  $('body').append('\
    <div id="pop_cancel_reservations" class="pop_window" style="width: 50%">\
      确认取消选中的预约？<br>\
      <div id="pop_cancel_reservations_div"></div>\
      <button type="button" onclick="$(\'#pop_cancel_reservations\').remove();cancelReservationsConfirm();">确认</button>\
      <button type="button" onclick="$(\'#pop_cancel_reservations\').remove();">取消</button>\
    </div>\
  ');
  for (var i = 0; i < reservations.length; i++) {
    if ($('#cell_checkbox_' + i)[0].checked) {
      $('#pop_cancel_reservations_div').append('\
        <span>' + reservations[i].teacher_fullname + '　　　　' + reservations[i].start_time + '至' + reservations[i].end_time + '</span><br>\
      ');
    }
  }
  optimize('#pop_cancel_reservations');
}

function cancelReservationsConfirm() {
  var reservationIds = [];
  var sourceIds = [];
  var startTimes = [];
  for (var i = 0; i < reservations.length; ++i) {
    if ($('#cell_checkbox_' + i)[0].checked) {
      reservationIds.push(reservations[i].id);
      sourceIds.push(reservations[i].source_id);
      startTimes.push(reservations[i].start_time)
    }
  }
  var payload = {
    reservation_ids: reservationIds,
    source_ids: sourceIds,
    start_times: startTimes,
  };
  $.ajax({
    url: '/api/admin/reservation/cancel',
    type: "POST",
    dataType: 'json',
    data: payload,
    traditional: true,
  })
  .done(function(data) {
    if (data.status === 'OK') {
      viewReservations();
    } else {
      alert(data.err_msg);
    }
  });
}

function getFeedback(index) {
  $.post('/api/admin/reservation/feedback/get', {
    reservation_id: reservations[index].id,
    source_id: reservations[index].source_id,
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      showFeedback(index, data.payload.feedback);
    } else {
      alert(data.err_msg);
    }
  });
}

function showFeedback(index, feedback) {
  $('body').append('\
    <div class="pop_window" id="feedback_table_' + index + '" style="text-align: left; width: 50%">\
      咨询师反馈表<br>\
      评估分类：<br>\
      <select id="category_first_' + index + '" onchange="showSecondCategory(' + index + ')"><option value="">请选择</option></select><br>\
      <select id="category_second_' + index + '"></select><br>\
      出席人员：<br>\
      <input id="participant_student_' + index + '" type="checkbox">学生</input><input id="participant_parents_' + index + '" type="checkbox">家长</input>\
      <input id="participant_teacher_' + index + '" type="checkbox">教师</input><input id="participant_instructor_' + index + '" type="checkbox">辅导员</input>\
      <input id="participant_other_' + index + '" type="checkbox">其他</input><br>\
      重点明细：<select id="emphasis_'+ index + '"><option value="0">否</option><option value="1">是</option></select><br>\
      <div id="div_emphasis_' + index + '" style="display: none">\
        <b>严重程度：</b>\
        <input id="severity_' + index + '_0" type="checkbox">缓考</input>\
        <input id="severity_' + index + '_1" type="checkbox">休学复学</input>\
        <input id="severity_' + index + '_2" type="checkbox">家属陪读</input>\
        <input id="severity_' + index + '_3" type="checkbox">家属不知情</input>\
        <input id="severity_' + index + '_4" type="checkbox">任何其他需要知会院系关注的原因</input>\
        <br>\
        <b>疑似或明确的医疗诊断：</b>\
        <input id="medical_diagnosis_' + index + '_0" type="checkbox">服药</input>\
        <input id="medical_diagnosis_' + index + '_1" type="checkbox">精神分裂</input>\
        <input id="medical_diagnosis_' + index + '_2" type="checkbox">双相情感障碍</input>\
        <input id="medical_diagnosis_' + index + '_3" type="checkbox">焦虑症（状态）</input>\
        <input id="medical_diagnosis_' + index + '_4" type="checkbox">抑郁症（状态）</input>\
        <br>　　　　　　　　　\
        <input id="medical_diagnosis_' + index + '_5" type="checkbox">强迫症</input>\
        <input id="medical_diagnosis_' + index + '_6" type="checkbox">进食障碍</input>\
        <input id="medical_diagnosis_' + index + '_7" type="checkbox">失眠</input>\
        <input id="medical_diagnosis_' + index + '_8" type="checkbox">其他精神症状</input>\
        <input id="medical_diagnosis_' + index + '_9" type="checkbox">躯体疾病</input>\
        <input id="medical_diagnosis_' + index + '_10" type="checkbox">不遵医嘱</input>\
        <br>\
        <b>危急情况：</b>\
        <input id="crisis_' + index + '_0" type="checkbox">自伤</input>\
        <input id="crisis_' + index + '_1" type="checkbox">伤害他人</input>\
        <input id="crisis_' + index + '_2" type="checkbox">自杀念头</input>\
        <input id="crisis_' + index + '_3" type="checkbox">自杀未遂</input>\
      </div>\
      咨询记录：<br>\
      <textarea id="record_' + index + '" style="width: 100%; height:80px"></textarea><br>\
      是否危机个案：<select id="crisis_level_'+ index + '"><option value="0">否</option><option value="3">三星</option><option value="4">四星</option><option value="5">五星</option></select><br>\
      <button type="button" onclick="submitFeedback(' + index + ');">提交</button>\
      <button type="button" onclick="$(\'#feedback_table_' + index + '\').remove();">取消</button>\
    </div>\
  ');
  $(function() {
    showFirstCategory(index);
    if (feedback.category.length > 0) {
      $('#category_first_' + index).val(feedback.category.charAt(0));
      $('#category_first_' + index).change();
      $('#category_second_' + index).val(feedback.category);
    }
    if (feedback.participants.length > 0) {
      $('#participant_student_' + index).first().attr('checked', feedback.participants[0] > 0);
      $('#participant_parents_' + index).first().attr('checked', feedback.participants[1] > 0);
      $('#participant_teacher_' + index).first().attr('checked', feedback.participants[2] > 0);
      $('#participant_instructor_' + index).first().attr('checked', feedback.participants[3] > 0);
      $('#participant_other_' + index).first().attr('checked', feedback.participants[4] > 0);
    }
    var i = 1;
    for (i = 0; i < 5; i++) {
      $('#severity_' + index + '_' + i).first().attr('checked', feedback.severity[i] > 0);
    }
    for (i = 0; i < 11; i++) {
      $('#medical_diagnosis_' + index + '_' + i).first().attr('checked', feedback.medical_diagnosis[i] > 0);
    }
    for (i = 0; i < 4; i++) {
      $('#crisis_' + index + "_" + i).first().attr('checked', feedback.crisis[i] > 0);
    }
    $('#record_' + index).val(feedback.record);
    $('#emphasis_' + index).change(function() {
      if ($('#emphasis_' + index).val() === "0") {
        $('#div_emphasis_' + index).hide();
      } else {
        $('#div_emphasis_' + index).show();
      }
    });
    $('#emphasis_' + index).val(feedback.emphasis);
    $('#emphasis_' + index).change();
    $('#crisis_level_' + index).val(feedback.crisis_level);
    optimize('#feedback_table_' + index);
  });
}

function showFirstCategory(index) {
  for (var name in firstCategory) {
    if (firstCategory.hasOwnProperty(name)) {
      $('#category_first_' + index).append($("<option>", {
        value: name,
        text: firstCategory[name],
      }));
    }
  }
}

function showSecondCategory(index) {
  var first = $('#category_first_' + index).val();
  $('#category_second_' + index).find("option").remove().end().append('<option value="">请选择</option>').val('');
  if ($('#category_first_' + index).selectedIndex === 0) {
    return;
  }
  if (secondCategory.hasOwnProperty(first)) {
    for (var name in secondCategory[first]) {
      if (secondCategory[first].hasOwnProperty(name)) {
        var option = new Option(name, secondCategory[first][name]);
        $('#category_second_' + index).append($("<option>", {
          value: name,
          text: secondCategory[first][name],
        }));
      }
    }
  }
}

function submitFeedback(index) {
  var participants = [];
  participants.push($('#participant_student_' + index).first().is(':checked') ? 1 : 0);
  participants.push($('#participant_parents_' + index).first().is(':checked') ? 1 : 0);
  participants.push($('#participant_teacher_' + index).first().is(':checked') ? 1 : 0);
  participants.push($('#participant_instructor_' + index).first().is(':checked') ? 1 : 0);
  participants.push($('#participant_other_' + index).first().is(':checked') ? 1 : 0);
  var isEmphasis = $('#emphasis_' + index).val() !== "0"
  var i = 1;
  var severity = [];
  for (i = 0; i < 5; i++) {
    severity.push(isEmphasis ? ($('#severity_' + index + '_' + i).first().is(':checked') ? 1 : 0) : 0);
  }
  var medicalDiagnosis = [];
  for (i = 0; i < 11; i++) {
    medicalDiagnosis.push(isEmphasis ? ($('#medical_diagnosis_' + index + '_' + i).first().is(':checked') ? 1 : 0) : 0);
  }
  var crisis = [];
  for (i = 0; i < 4; i++) {
    crisis.push(isEmphasis ? ($('#crisis_' + index + "_" + i).first().is(':checked') ? 1 : 0) : 0);
  }
  var payload = {
    reservation_id: reservations[index].id,
    category: $('#category_second_' + index).val(),
    participants: participants,
    emphasis: $('#emphasis_' + index).val(),
    severity: severity,
    medical_diagnosis: medicalDiagnosis,
    crisis: crisis,
    record: $('#record_' + index).val(),
    crisis_level: $('#crisis_level_' + index).val(),
  };
  $.ajax({
    url: '/api/admin/reservation/feedback/submit',
    type: 'POST',
    dataType: 'json',
    data: payload,
    traditional: true,
  })
  .done(function(data) {
    if (data.status === 'OK') {
      successFeedback(index);
      viewReservations();
    } else {
      alert(data.err_msg);
    }
  });
}

function successFeedback(index) {
  $('#feedback_table_' + index).remove();
  $('body').append('\
    <div id="pop_success_feedback" class="pop_window" style="width: 50%;">\
      您已成功提交反馈！<br>\
      <button type="button" onclick="$(\'#pop_success_feedback\').remove();">确定</button>\
    </div>\
  ');
  optimize('#pop_success_feedback');
}

function setStudent(index) {
  $('body').append('\
    <div id="pop_set_student" class="pop_window" style="width: 50%;">\
      请输入您要指定的学生学号（必须为已注册学生）：<br>\
      <input id="student_username_' + index + '" oninput="searchStudent(' + index + ');"/><br>\
      <span id="search_student_' + index + '_span" style="color: red;"></span>\
      <input class="checkbox" type="checkbox" id="student_sms_' + index + '"/>是否发送提醒短信给学生<br>\
      <br>\
      <div style="text-align:left; margin: 0 20% 0 20%">\
      <div style="text-align:center;font-size:23px">学生信息登记表</div><br>\
      姓　　名：<input id="student_fullname_' + index + '"><br>\
      性　　别：<select id="student_gender_' + index + '"><option value="">请选择</option><option value="男">男</option><option value="女">女</option></select><br>\
      出生日期：<input id="student_birthday_' + index + '"><br>\
      系　　别：<input id="student_school_' + index + '"><br>\
      年　　级：<input id="student_grade_' + index + '"><br>\
      现在住址：<input id="student_current_address_' + index + '"><br>\
      家庭住址：<input id="student_family_address_' + index + '"><br>\
      联系电话：<input id="student_mobile_' + index + '"><br>\
      邮　　箱：<input id="student_email_' + index + '"><br>\
      咨询经历：<br>\
      时间：<input id="student_experience_time_' + index + '"><br>\
      地点：<input id="student_experience_location_' + index + '" style="width:60px">\
      咨询师：<input id="student_experience_teacher_' + index + '" style="width:60px"><br>\
      父　　亲<br>\
      年龄：<input id="student_father_age_' + index + '" style="width:20px"> 职业：<input id="student_father_job_' + index + '" style="width:40px"> 学历：<input id="student_father_edu_' + index + '" style="width:40px"><br>\
      母　　亲<br>\
      年龄：<input id="student_mother_age_' + index + '" style="width:20px"> 职业：<input id="student_mother_job_' + index + '" style="width:40px"> 学历：<input id="student_mother_edu_' + index + '" style="width:40px"><br>\
      父母婚姻状况：<select id="student_parent_marriage_' + index + '"><option value="">请选择</option><option value="良好">良好</option><option value="一般">一般</option><option value="离婚">离婚</option><option value="再婚">再婚</option></select><br>\
      在近三个月里，是否发生了对你有重大意义的事（如亲友的死亡、法律诉讼、失恋等）？<br>\
      <textarea id="student_significant_' + index + '"></textarea><br>\
      你现在需要接受帮助的主要问题是什么？<br>\
      <textarea id="student_problem_' + index + '"></textarea><br>\
      </div>\
      <button type="button" onclick="setStudentConfirm(' + index + ');">确认</button>\
      <button type="button" style="margin-left:20px" onclick="$(\'#pop_set_student\').remove();">取消</button>\
    </div>\
  ');
  optimize('#pop_set_student');
}

function searchStudent(index) {
  $.post('/api/admin/student/search', {
    student_username: $('#student_username_' + index).val(),
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      $('#student_fullname_' + index).val(data.payload.student.fullname);
      $('#student_gender_' + index).val(data.payload.student.gender);
      $('#student_birthday_' + index).val(data.payload.student.birthday);
      $('#student_school_' + index).val(data.payload.student.school);
      $('#student_grade_' + index).val(data.payload.student.grade);
      $('#student_current_address_' + index).val(data.payload.student.current_address);
      $('#student_family_address_' + index).val(data.payload.student.family_address);
      $('#student_mobile_' + index).val(data.payload.student.mobile);
      $('#student_email_' + index).val(data.payload.student.email);
      $('#student_experience_time_' + index).val(data.payload.student.experience_time);
      $('#student_experience_location_' + index).val(data.payload.student.experience_location);
      $('#student_experience_teacher_' + index).val(data.payload.student.experience_teacher);
      $('#student_father_age_' + index).val(data.payload.student.father_age);
      $('#student_father_job_' + index).val(data.payload.student.father_job);
      $('#student_father_edu_' + index).val(data.payload.student.father_edu);
      $('#student_mother_age_' + index).val(data.payload.student.mother_age);
      $('#student_mother_job_' + index).val(data.payload.student.mother_job);
      $('#student_mother_edu_' + index).val(data.payload.student.mother_edu);
      $('#student_parent_marriage_' + index).val(data.payload.student.parent_marriage);
      $('#student_significant_' + index).val(data.payload.student.significant);
      $('#student_problem_' + index).val(data.payload.student.problem);
      $('#search_student_' + index + "_span").text('载入成功');
    } else {
      $('#search_student_' + index + "_span").text(data.err_msg);
    }
  });
}

function setStudentConfirm(index) {
  $.post('/api/admin/reservation/student/set', {
    reservation_id: reservations[index].id,
    source_id: reservations[index].source_id,
    start_time: reservations[index].start_time,
    student_username: $('#student_username_' + index).val(),
    student_fullname: $('#student_fullname_' + index).val(),
    student_gender: $('#student_gender_' + index).val(),
    student_birthday: $('#student_birthday_' + index).val(),
    student_school: $('#student_school_' + index).val(),
    student_grade: $('#student_grade_' + index).val(),
    student_current_address: $('#student_current_address_' + index).val(),
    student_family_address: $('#student_family_address_' + index).val(),
    student_mobile: $('#student_mobile_' + index).val(),
    student_email: $('#student_email_' + index).val(),
    student_experience_time: $('#student_experience_time_' + index).val(),
    student_experience_location: $('#student_experience_location_' + index).val(),
    student_experience_teacher: $('#student_experience_teacher_' + index).val(),
    student_father_age: $('#student_father_age_' + index).val(),
    student_father_job: $('#student_father_job_' + index).val(),
    student_father_edu: $('#student_father_edu_' + index).val(),
    student_mother_age: $('#student_mother_age_' + index).val(),
    student_mother_job: $('#student_mother_job_' + index).val(),
    student_mother_edu: $('#student_mother_edu_' + index).val(),
    student_parent_marriage: $('#student_parent_marriage_' + index).val(),
    student_significant: $('#student_significant_' + index).val(),
    student_problem: $('#student_problem_' + index).val(),
    student_sms: $('#student_sms_' + index)[0].checked,
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      successSetStudent();
    } else {
      alert(data.err_msg);
    }
  });
}

function successSetStudent() {
  $('#pop_set_student').remove();
  $('body').append('\
    <div id="pop_success_set_student" class="pop_window" style="width: 50%;">\
      成功指定学生！<br>\
      <button type="button" onclick="$(\'#pop_success_set_student\').remove();viewReservations();">确定</button>\
    </div>\
  ');
  optimize('#pop_success_set_student');
}

function getStudent(index) {
  $.post('/api/admin/student/get', {
    student_id: reservations[index].student_id
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      showStudent(data.payload.student, data.payload.reservations);
    } else {
      alert(data.err_msg);
    }
  });
}

function showStudent(student, reservations) {
  $('body').append('\
    <div id="pop_show_student_' + student.id + '" class="pop_window" style="text-align: left; height: 80%; overflow:auto;">\
      <div style="width: 60%; float: left;">\
        学号：' + student.username + '<br>\
        姓名：' + student.fullname + '<br>\
        性别：' + student.gender + '<br>\
        出生日期：' + student.birthday + '<br>\
        系别：' + student.school + '<br>\
        年级：' + student.grade + '<br>\
        现住址：' + student.current_address + '<br>\
        家庭住址：' + student.family_address + '<br>\
        联系电话：' + student.mobile + '<br>\
        Email：' + student.email + '<br>\
        咨询经历：' + (student.experience_time && student.experience_time !== '' ? '时间：' + student.experience_time + ' 地点：' + student.experience_location + ' 咨询师：' + student.experience_teacher : '无') + '<br>\
        父亲年龄：' + student.father_age + ' 职业：' + student.father_job + ' 学历：' + student.father_edu + '<br>\
        母亲年龄：' + student.mother_age + ' 职业：' + student.mother_job + ' 学历：' + student.mother_edu + '<br>\
        父母婚姻状况：' + student.parent_marriage + '<br>\
        近三个月里发生的有重大意义的事：' + student.significant + '<br>\
        需要接受帮助的主要问题：' + student.problem + '<br>\
        <br>\
        是否危机个案：<select id="crisis_level_'+ student.id + '"><option value="0">否</option><option value="3">三星</option><option value="4">四星</option><option value="5">五星</option></select>\
        <button type="button" onclick="updateCrisisLevel(\'' + student.id + '\');">更新</button>\
        <span id="crisis_level_tip_' + student.id + '" style="color: red;"></span><br>\
        档案分类：<input id="archive_category_' + student.id + '" type="text" value="' + student.archive_category + '" style="width: 100px"/>\
        档案编号：<input id="archive_number_' + student.id + '" type="text" value="' + student.archive_number + '" style="width: 100px"/>\
        <button type="button" onclick="updateArchiveNumber(\'' + student.id + '\');">更新</button>\
        <span id="archive_number_tip_' + student.id + '" style="color: red;"></span><br>\
        已绑定的咨询师：<span id="binded_teacher_username_' + student.id + '">' + student.binded_teacher_username + '</span>&nbsp;\
          <span id="binded_teacher_fullname_' + student.id + '">' + student.binded_teacher_fullname + '</span>\
          <button type="button" onclick="unbindStudent(\'' + student.id + '\');">解绑</button><br>\
        请输入匹配咨询师工号：<input id="teacher_username_' + student.id + '" type="text"/>\
        <button type="button" onclick="bindStudent(\'' + student.id + '\');">绑定</button><br>\
        <div style="margin: 10px 0">\
          <button type="button" onclick="exportStudent(\'' + student.id + '\');">导出</button>\
          <button type="button" onclick="$(\'#pop_show_student_' + student.id + '\').remove();">关闭</button>\
        </div>\
        <div id="student_reservations_' + student.id + '" style="width: 600px">\
        </div>\
      </div>\
      <div style="width: 35%; float: right; border: 2px solid red; padding: 2px;">\
        <p style="margin-top: 0; background-color: red">账户相关，谨慎操作</p>\
        新密码：<input id="password_' + student.id + '" type="password"/><br>\
        确认密码：<input id="password_check_' + student.id + '" type="password"/><br>\
        <button type="button" onclick="resetStudentPassword(\'' + student.id + '\');">重置密码</button>\
        <p>删除账户</p>\
        <button type="button" onclick="deleteStudentAccount(\'' + student.id + '\');" class="btn btn-danger">删除账户</button>\
      </div>\
    </div>\
  ');
  for (var i = 0; i < reservations.length; i++) {
    $('#student_reservations_' + student.id).append('\
      <div class="has_children" style="background: ' + (reservations[i].status === 3 ? '#555' : '#F00') + '">\
        <span>' + reservations[i].start_time + ' 至 ' + reservations[i].end_time + '  ' + reservations[i].teacher_fullname + '</span>\
        <p class="children">学生反馈：' + reservations[i].student_feedback.scores + '</p>\
        <p class="children">评估分类：' + reservations[i].teacher_feedback.category + '</p>\
        <p class="children">出席人员：' + reservations[i].teacher_feedback.participants + '</p>\
        <p class="children">重点明细：' + (reservations[i].teacher_feedback.emphasis === '0' ? '否' : '是') + '</p>'
        + (reservations[i].teacher_feedback.severity === '' ? '' : '<p class="children">严重程度：' + reservations[i].teacher_feedback.severity + '</p>')
        + (reservations[i].teacher_feedback.medical_diagnosis === '' ? '' : '<p class="children">疑似或明确的医疗诊断：' + reservations[i].teacher_feedback.medical_diagnosis + '</p>')
        + (reservations[i].teacher_feedback.crisis === '' ? '' : '<p class="children">危急情况：' + reservations[i].teacher_feedback.crisis + '</p>')
        + '<p class="children">咨询记录：' + reservations[i].teacher_feedback.record + '</p>\
      </div>\
    ');
  }
  $(function() {
    $('.has_children').click(function() {
      $(this).addClass('highlight').children('p').show().end()
          .siblings().removeClass('highlight').children('p').hide();
    });
    $('#crisis_level_' + student.id).val(student.crisis_level);
  });
  optimize('#pop_show_student_' + student.id);
}

function updateCrisisLevel(studentId) {
  var payload = {
    student_id: studentId,
    crisis_level: $('#crisis_level_' + studentId).val(),
  }
  $.ajax({
    url: '/api/admin/student/crisis/update',
    type: 'POST',
    dataType: 'json',
    data: payload,
    traditional: true,
  })
  .done(function(data) {
    if (data.status === 'OK') {
      $('#crisis_level_tip_' + studentId).text('更新成功！');
      viewReservations();
    } else {
      alert(data.err_msg);
    }
  });
}

function updateArchiveNumber(studentId) {
  $.post('/api/admin/student/archive/update', {
    student_id: studentId,
    archive_category: $('#archive_category_' + studentId).val(),
    archive_number: $('#archive_number_' + studentId).val(),
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      $('#archive_number_tip_' + studentId).text('更新成功！');
    } else {
      alert(data.err_msg);
    }
  });
}

function resetStudentPassword(studentId) {
  $('body').append('\
    <div id="pop_reset" class="pop_window" style="width: 50%;">\
      您确定要重置该生的密码？<br>\
      <button type="button" onclick="$(\'#pop_reset\').remove();resetStudentPasswordConfirm(\'' + studentId + '\');">确认</button>\
      <button type="button" style="margin-left:20px" onclick="$(\'#pop_reset\').remove();">取消</button>\
    </div>\
  ');
  optimize('#pop_reset');
}

function resetStudentPasswordConfirm(studentId) {
  var password = $('#password_' + studentId).val();
  var passwordConfirm = $('#password_check_' + studentId).val();
  if (password !== passwordConfirm) {
    alert('两次密码不一致，请重新输入');
    $('#password_' + studentId).val('');
    $('#password_check_' + studentId).val('');
    return;
  }
  $.post('/api/admin/student/password/reset', {
    student_id: studentId,
    password: password,
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      resetStudentPasswordSuccess();
    } else {
      alert(data.err_msg);
    }
  });
}

function resetStudentPasswordSuccess() {
  $('body').append('\
    <div id="pop_reset_success" class="pop_window" style="width: 50%;">\
      密码重置成功！<br>\
      <button type="button" onclick="$(\'#pop_reset_success\').remove();">确定</button>\
    </div>\
  ');
  optimize('#pop_reset_success');
}

function deleteStudentAccount(studentId) {
  $('body').append('\
    <div id="pop_delete_student_account" class="pop_window" style="width: 50%;">\
      <span style="color: red"><b>删除学生账户不可恢复，请确认操作</b></span><br>\
      <button type="button" onclick="$(\'#pop_delete_student_account\').remove();deleteStudentAccountConfirm(\'' + studentId + '\');">确认</button>\
      <button type="button" style="margin-left:20px" onclick="$(\'#pop_delete_student_account\').remove();">取消</button>\
    </div>\
  ');
  optimize('#pop_delete_student_account');
}

function deleteStudentAccountConfirm(studentId) {
  $('body').append('\
    <div id="pop_delete_student_account_confirm" class="pop_window" style="width: 50%;">\
      <span style="color: red;"><b>请输入管理员密码删除学生账户</b></span><br>\
      <input id="delete_student_account_confirm_' + studentId + '" type="password"><br>\
      <button type="button" onclick="deleteStudentAccountConfirmCheck(\'' + studentId + '\');">确认</button>\
      <button type="button" style="margin-left:20px" onclick="$(\'#pop_delete_student_account_confirm\').remove();">取消</button>\
    </div>\
  ');
  optimize('#pop_delete_student_account_confirm');
}

function deleteStudentAccountConfirmCheck(studentId) {
  var adminPassword = $('#delete_student_account_confirm_' + studentId).val();
  $('#pop_delete_student_account_confirm').remove();
  $.post('/api/admin/student/account/delete', {
    student_id: studentId,
    password: adminPassword,
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      deleteStudentAccountSuccess();
    } else {
      alert(data.err_msg);
    }
  });
}

function deleteStudentAccountSuccess(studentId) {
  $('#pop_show_student_' + studentId).remove();
  $('body').append('\
    <div id="pop_delete_student_account_success" class="pop_window" style="width: 50%;">\
      学生账户已删除！<br>\
      <button type="button" onclick="$(\'#pop_delete_student_account_success\').remove();">确定</button>\
    </div>\
  ');
  optimize('#pop_delete_student_account_success');
}

function exportStudent(studentId) {
  $.post('/api/admin/student/export', {
    student_id: studentId,
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      window.open(data.payload.url);
    } else {
      alert(data.err_msg);
    }
  });
}

function unbindStudent(studentId) {
  $.post('/api/admin/student/unbind', {
    student_id: studentId,
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      $('#binded_teacher_username_' + studentId).text(data.payload.student.binded_teacher_username);
      $('#binded_teacher_fullname_' + studentId).text(data.payload.student.binded_teacher_fullname);
    } else {
      alert(data.err_msg);
    }
  });
}

function bindStudent(studentId) {
  $.post('/api/admin/student/bind', {
    student_id: studentId,
    teacher_username: $('#teacher_username_' + studentId).val(),
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      $('#binded_teacher_username_' + studentId).text(data.payload.student.binded_teacher_username);
      $('#binded_teacher_fullname_' + studentId).text(data.payload.student.binded_teacher_fullname);
    } else {
      alert(data.err_msg);
    }
  });
}

function getWorkload() {
  $.post('/api/admin/teacher/workload', {
    from_date: $('#workload_from').val(),
    to_date: $('#workload_to').val(),
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      showWorkload(data.payload.workload);
    } else {
      alert(data.err_msg);
    }
  });
}

function showWorkload(workload) {
  $('body').append('\
    <div id="pop_show_workload" class="pop_window" style="text-align: left; width: 50%; height: 70%; overflow: auto;">\
      咨询师工作量统计\
      <div id="teacher_workload" style="width: 600px; margin-top: 10px;">\
        <div class="table_col" id="col_workload_username">\
          <div class="table_head table_cell">咨询师工号</div>\
        </div>\
        <div class="table_col" id="col_workload_fullname">\
          <div class="table_head table_cell">咨询师姓名</div>\
        </div>\
        <div class="table_col" id="col_workload_student">\
          <div class="table_head table_cell">咨询人数</div>\
        </div>\
        <div class="table_col" id="col_workload_reservation">\
          <div class="table_head table_cell">咨询人次</div>\
        </div>\
        <div class="clearfix"></div>\
      </div>\
      <div style="margin: 10px 0">\
        <button type="button" onclick="$(\'#pop_show_workload\').remove();">关闭</button>\
      </div>\
    </div>\
  ');
  $('#col_workload_username').width(100);
  $('#col_workload_fullname').width(100);
  $('#col_workload_student').width(80);
  $('#col_workload_reservation').width(80);
  for (var i in workload) {
    if (workload.hasOwnProperty(i)) {
      $('#col_workload_username').append('<div class="table_cell" id="cell_workload_username_'
        + i + '">' + workload[i].teacher_username + '</div>');
      $('#col_workload_fullname').append('<div class="table_cell" id="cell_workload_fullname_'
        + i + '">' + workload[i].teacher_fullname + '</div>');
      $('#col_workload_student').append('<div class="table_cell" id="cell_workload_student_'
        + i + '">' + Object.size(workload[i].students) + '</div>');
      $('#col_workload_reservation').append('<div class="table_cell" id="cell_workload_reservation_'
        + i + '">' + Object.size(workload[i].reservations) + '</div>');
    }
  }
  optimize('#pop_show_workload');
}

Object.size = function(obj) {
  var size = 0, key;
  for (key in obj) {
    if (obj.hasOwnProperty(key)) size++;
  }
  return size;
}

function exportReportMonthly() {
  $.post('/api/admin/reservation/export/report/monthly', {
    monthly_date: $('#monthly_report_date').val(),
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      window.open(data.payload.report_url);
      // window.open(data.payload.key_case_url);
    } else {
      alert(data.err_msg);
    }
  });
}

function changeAdminPassword() {
  $('body').append('\
    <div id="change_admin_password" class="pop_window" style="text-align: left; width: 30%">\
      <p style="color: red;">更改管理员密码</p>\
      输入用户名：<input id="admin_username"><br>\
      输入旧密码：<input id="admin_old_password" type="password"><br>\
      输入新密码：<input id="admin_new_password" type="password"><br>\
      确认新密码：<input id="admin_new_password_confirm" type="password"><br>\
      <span id="change_admin_password_tip" style="color: red;"></span>\
      <br>\
      <button onclick="changeAdminPasswordConfirm();">确定</button><button style="margin-left: 20px;" onclick="$(\'#change_admin_password\').remove();">取消</button>\
    </div>\
  ');
  optimize('#change_admin_password');
}

function changeAdminPasswordConfirm() {
  var username = $('#admin_username').val();
  var oldPassword = $('#admin_old_password').val();
  var newPassword = $('#admin_new_password').val();
  var newPasswordConfirm = $('#admin_new_password_confirm').val();
  if (username === '' || oldPassword === '' || newPassword === '' || newPasswordConfirm === '') {
    alert('请填写完整');
    return;
  } else if (newPassword !== newPasswordConfirm) {
    alert('两次密码输入不一致');
    $('#admin_new_password').val('');
	  $('#admin_new_password_confirm').val('');
	  return;
  }
  $.post('/api/user/admin/password/change', {
    username: username,
    old_password: oldPassword,
    new_password: newPassword,
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      $('#change_admin_password_tip').text('更改成功');
	    $('#admin_username').val('');
	    $('#admin_old_password').val('');
	    $('#admin_new_password').val('');
	    $('#admin_new_password_confirm').val('');
    } else {
      alert(data.err_msg);
    }
  });
}

function resetTeacherPassword() {
  $('body').append('\
    <div id="reset_teacher_password" class="pop_window" style="text-align: left; width: 30%">\
      <p style="color: red;">重置咨询师密码</p>\
      咨询师工号：<input id="reset_teacher_username"><br>\
      咨询师姓名：<input id="reset_teacher_fullname"><br>\
      咨询师手机：<input id="reset_teacher_mobile"><br>\
      输入新密码：<input id="reset_teacher_new_password" type="password"><br>\
      确认新密码：<input id="reset_teacher_new_password_confirm" type="password"><br>\
      <span id="reset_teacher_password_tip" style="color: red;"></span>\
      <br>\
      <button onclick="resetTeacherPasswordConfirm();">确定</button><button style="margin-left: 20px;" onclick="$(\'#reset_teacher_password\').remove();">取消</button>\
    </div>\
  ');
  optimize('#reset_teacher_password');
}

function resetTeacherPasswordConfirm() {
  var username = $('#reset_teacher_username').val();
  var fullname = $('#reset_teacher_fullname').val();
  var mobile = $('#reset_teacher_mobile').val();
  var password = $('#reset_teacher_new_password').val();
  var passwordConfirm = $('#reset_teacher_new_password_confirm').val();
  console.log(username, fullname, mobile, password, passwordConfirm)
  if (username === '' || fullname === '' || mobile === '' || password === '' || passwordConfirm === '') {
    alert('请填写完整');
    return;
  } else if (password !== passwordConfirm) {
	  alert('两次密码输入不一致');
	  $('#reset_teacher_new_password').val('');
	  $('#reset_teacher_new_password_confirm').val('');
	  return;
  } 
  $.post('/api/admin/teacher/password/reset', {
    "teacher_username": username,
    "teacher_fullname": fullname,
    "teacher_mobile": mobile,
    "password": password,
  }, function(data, textStatus, xhr) {
    if (data.status === 'OK') {
      $('#reset_teacher_password_tip').text('重置成功');
	    $('#reset_teacher_username').val('');
      $('#reset_teacher_fullname').val('');
      $('#reset_teacher_mobile').val('');
      $('#reset_teacher_new_password').val('');
      $('#reset_teacher_new_password_confirm').val('');
    } else {
      alert(data.err_msg);
    }
  });
}

function clearAllStudentsBindedTeacher() {
	$('body').append('\
    <div id="clear_all_students_binded_teacher" class="pop_window" style="text-align: left; width: 30%">\
      <p style="color: red;">清除所有学生绑定咨询师</p>\
      <p>操作不可恢复，请谨慎！！！</p>\
      <span style="color: red;"><b>请输入管理员密码</b></span><br>\
      <input id="clear_all_students_binded_teacher_confirm" type="password"><br>\
      <button onclick="clearAllStudentsBindedTeacherConfirm();">确定</button><button style="margin-left: 20px;" onclick="$(\'#clear_all_students_binded_teacher\').remove();">取消</button>\
    </div>\
  ');
	optimize('#clear_all_students_binded_teacher');
}

function clearAllStudentsBindedTeacherConfirm() {
	$.post('/api/admin/student/unbind/all', {
	  password: $('#clear_all_students_binded_teacher_confirm').val(),
  }, function(data, textStatus, xhr) {
	  if (data.status === 'OK') {
		  clearAllStudentsBindedTeacherSuccess();
    } else {
	    alert(data.err_msg);
    }
  });
}

function clearAllStudentsBindedTeacherSuccess() {
	$('#clear_all_students_binded_teacher').remove();
  $('body').append('\
    <div id="pop_clear_all_students_binded_teacher_success" class="pop_window" style="width: 50%;">\
      所有学生的绑定咨询师已清除！<br>\
      <button type="button" onclick="$(\'#pop_clear_all_students_binded_teacher_success\').remove();">确定</button>\
    </div>\
  ');
  optimize('#pop_clear_all_students_binded_teacher_success');
}
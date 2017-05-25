var width = $(window).width();
var height = $(window).height();
var consultationCrisisList = [];

function getReservationConsultationCrisisWithTimeRange() {
  $.getJSON('/api/admin/reservation/consultation/crisis', {
    from_date: $('#from_date').val(),
    to_date: $('#to_date').val(),
  }, function(json, textStatus) {
    console.log(json + "aaa");
    if (json.status === 'OK') {
      console.log(json.payload);
      consultationCrisisList = json.payload.consultation_crisis;
	    refreshDataTable(json.payload.consultation_crisis);
      optimize();
    } else {
      alert(json.err_msg);
    }
  });
}

function getReservationConsultationCrisisWithStudentFullname() {
	$.getJSON('/api/admin/reservation/consultation/crisis', {
		student_username: $('#student_username').val(),
	}, function(json, textStatus) {
		console.log(json + "aaa");
		if (json.status === 'OK') {
			console.log(json.payload);
			consultationCrisisList = json.payload.consultation_crisis;
			refreshDataTable(json.payload.consultation_crisis);
			optimize();
		} else {
			alert(json.err_msg);
		}
	});
}

function getReservationConsultationCrisisWithSchoolContact() {
	$.getJSON('/api/admin/reservation/consultation/crisis', {
		student_username: $('#student_username').val(),
	}, function(json, textStatus) {
		console.log(json + "aaa");
		if (json.status === 'OK') {
			console.log(json.payload);
			consultationCrisisList = json.payload.consultation_crisis;
			refreshDataTable(json.payload.consultation_crisis);
			optimize();
		} else {
			alert(json.err_msg);
		}
	});
}

function initDataTable() {
    $('#page_maintable').first().html('\
        <div class="table_col" id="col_date">\
            <div class="table_head table_cell">日期</div>\
        </div>\
        <div class="table_col" id="col_fullname">\
            <div class="table_head table_cell">姓名</div>\
        </div>\
        <div class="table_col" id="col_username">\
            <div class="table_head table_cell">学号</div>\
        </div>\
        <div class="table_col" id="col_gender">\
            <div class="table_head table_cell">性别</div>\
        </div>\
        <div class="table_col" id="col_academic">\
            <div class="table_head table_cell">学历</div>\
        </div>\
        <div class="table_col" id="col_school">\
            <div class="table_head table_cell">院系</div>\
        </div>\
        <div class="table_col" id="col_teacher_fullname">\
            <div class="table_head table_cell">接待咨询师</div>\
        </div>\
        <div class="table_col" id="col_school_contact">\
            <div class="table_head table_cell">院系联系人</div>\
        </div>\
        <div class="table_col" id="col_consultation_or_crisis">\
            <div class="table_head table_cell">会商or危机处理</div>\
        </div>\
        <div class="table_col" id="col_reservated_status">\
            <div class="table_head table_cell">来访情况</div>\
        </div>\
        <div class="table_col" id="col_category">\
            <div class="table_head table_cell">评估分类</div>\
        </div>\
        <div class="table_col" id="col_emphasis_str">\
            <div class="table_head table_cell">重点标记</div>\
        </div>\
        <div class="table_col" id="col_crisis_level">\
            <div class="table_head table_cell">星级</div>\
        </div>\
        <div class="table_col" id="col_is_send_notify">\
            <div class="table_head table_cell">是否发危机通报</div>\
        </div>\
        <div class="clearfix"></div>\
    ');
    optimize();
}

function refreshDataTable(consultationCrisisList) {
  initDataTable();
  for (var i = 0; i < consultationCrisisList.length; ++i) {
    var cc = consultationCrisisList[i];
    $('#col_date').append('<div class="table_cell" id="cell_date'
      + i + '">' + cc.date + '</div>');
	  $('#col_fullname').append('<div class="table_cell" id="cell_fullname'
		  + i + '">' + cc.fullname + '</div>');
	  $('#col_username').append('<div class="table_cell" id="cell_username'
		  + i + '">' + cc.username + '</div>');
	  $('#col_gender').append('<div class="table_cell" id="cell_gender'
		  + i + '">' + cc.gender + '</div>');
	  $('#col_academic').append('<div class="table_cell" id="cell_academic'
		  + i + '">' + cc.academic + '</div>');
	  $('#col_school').append('<div class="table_cell" id="cell_school'
		  + i + '">' + cc.school + '</div>');
	  $('#col_teacher_fullname').append('<div class="table_cell" id="cell_teacher_fullname'
		  + i + '">' + cc.teacher_fullname + '</div>');
	  $('#col_school_contact').append('<div class="table_cell" id="cell_school_contact'
		  + i + '">' + cc.school_contact + '</div>');
	  $('#col_consultation_or_crisis').append('<div class="table_cell" id="cell_consultation_or_crisis'
		  + i + '">' + cc.consultation_or_crisis + '</div>');
	  $('#col_reservated_status').append('<div class="table_cell" id="cell_reservated_status'
		  + i + '">' + cc.reservated_status + '</div>');
	  $('#col_category').append('<div class="table_cell" id="cell_category'
		  + i + '">' + cc.category + '</div>');
	  $('#col_emphasis_str').append('<div class="table_cell" id="cell_emphasis_str'
		  + i + '">' + cc.emphasis_str + '</div>');
	  $('#col_crisis_level').append('<div class="table_cell" id="cell_crisis_level'
		  + i + '">' + cc.crisis_level + '</div>');
	  $('#col_is_send_notify').append('<div class="table_cell" id="cell_is_send_notify'
		  + i + '">' + cc.is_send_notify + '</div>');
  }
}

function optimize(t) {
  for (var i = 0; i < consultationCrisisList.length; ++i) {
    if (i % 2 == 1) {
      $('#col_date' + i).css('background-color', 'white');
      $('#col_fullname' + i).css('background-color', 'white');
      $('#col_username' + i).css('background-color', 'white');
      $('#col_gender' + i).css('background-color', 'white');
      $('#col_academic' + i).css('background-color', 'white');
      $('#col_school' + i).css('background-color', 'white');
	    $('#col_teacher_fullname' + i).css('background-color', 'white');
	    $('#col_school_contact' + i).css('background-color', 'white');
	    $('#col_consultation_or_crisis' + i).css('background-color', 'white');
	    $('#col_reservated_status' + i).css('background-color', 'white');
	    $('#col_category' + i).css('background-color', 'white');
	    $('#col_emphasis_str' + i).css('background-color', 'white');
	    $('#col_crisis_level' + i).css('background-color', 'white');
	    $('#col_is_send_notify' + i).css('background-color', 'white');
    } else {
	    $('#col_date' + i).css('background-color', '#f3f3ff');
	    $('#col_fullname' + i).css('background-color', '#f3f3ff');
	    $('#col_username' + i).css('background-color', '#f3f3ff');
	    $('#col_gender' + i).css('background-color', '#f3f3ff');
	    $('#col_academic' + i).css('background-color', '#f3f3ff');
	    $('#col_school' + i).css('background-color', '#f3f3ff');
	    $('#col_teacher_fullname' + i).css('background-color', '#f3f3ff');
	    $('#col_school_contact' + i).css('background-color', '#f3f3ff');
	    $('#col_consultation_or_crisis' + i).css('background-color', '#f3f3ff');
	    $('#col_reservated_status' + i).css('background-color', '#f3f3ff');
	    $('#col_category' + i).css('background-color', '#f3f3ff');
	    $('#col_emphasis_str' + i).css('background-color', '#f3f3ff');
	    $('#col_crisis_level' + i).css('background-color', '#f3f3ff');
	    $('#col_is_send_notify' + i).css('background-color', '#f3f3ff');
    }
  }
  $('.table_head').height($('#head_select').height());
  $(t).css('left', (width - $(t).width()) / 2 - 11 + 'px');
  $(t).css('top', (height - $(t).height()) / 2 - 11 + 'px');
  var tableWidth = $('#col_date').width() + $('#col_fullname').width() + $('#col_username').width() + $('#col_gender').width() + $('#col_academic').width() + $('#col_school').width() + $('#col_teacher_fullname').width() + $('#col_school_contact').width() + $('#col_consultation_or_crisis').width() + $('#col_reservated_status').width() + $('#col_category').width() + $('#col_emphasis_str').width() + $('#col_crisis_level').width() + $('#col_is_send_notify').width();
  $('#page_maintable').css('margin-left', 0.5 * ($(window).width()
    - tableWidth) + 'px');
}

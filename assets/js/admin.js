var width = $(window).width();
var height = $(window).height();
var reservations;
var firstCategory;
var secondCategory;

function viewReservations() {
	if ($("#query_date").val() !== "") {
		queryReservations();
		return;
	}
	$.ajax({
		type: "GET",
		async: false,
		url: "/admin/reservation/view",
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				console.log(data);
				reservations = data.reservations;
				refreshDataTable(reservations);
				optimize();
			} else {
				alert(data.message);
			}
		}
	});
}

function queryReservations() {
	var payload = {
		from_date: $("#query_date").val(),
	};
	$.ajax({
		type: "GET",
		async: false,
		url: "/admin/reservation/view/daily",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				console.log(data);
				reservations = data.reservations;
				refreshDataTable(reservations);
				optimize();
			}
		},
	});
}

function exportTodayReservations() {
	$.ajax({
		type: "GET",
		async: false,
		url: "/admin/reservation/export/today",
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				window.open(data.url);
			} else {
				alert(data.message);
			}
		}
	});
}

function queryStudent() {
	var payload = {
		student_username: $("#query_student").val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/student/query",
		data: payload, 
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				showStudent(data.student_info, data.reservations);
			} else {
				alert(data.message);
			}
		}
	});
}

function refreshDataTable(reservations) {
	$("#page_maintable")[0].innerHTML = "\
		<div class='table_col' id='col_select'>\
			<div class='table_head table_cell' id='head_select'>\
				<button onclick='$(\".checkbox\").click();'>全选</button>\
			</div>\
		</div>\
		<div class='table_col' id='col_time'>\
			<div class='table_head table_cell'>时间</div>\
		</div>\
		<div class='table_col' id='col_teacher_fullname'>\
			<div class='table_head table_cell'>咨询师</div>\
		</div>\
		<div class='table_col' id='col_teacher_username'>\
			<div class='table_head table_cell'>咨询师编号</div>\
		</div>\
		<div class='table_col' id='col_teacher_mobile'>\
			<div class='table_head table_cell'>咨询师手机</div>\
		</div>\
		<div class='table_col' id='col_status'>\
			<div class='table_head table_cell'>状态</div>\
		</div>\
		<div class='table_col' id='col_student'>\
			<div class='table_head table_cell'>学生</div>\
		</div>\
		<div class='clearfix'></div>\
	";

	for (var i = 0; i < reservations.length; ++i) {
		$("#col_select").append("<div class='table_cell' id='cell_select_" + i + "'>"
			+ "<input class='checkbox' type='checkbox' id='cell_checkbox_" + i + "'></div>");
		$("#col_time").append("<div class='table_cell' id='cell_time_" + i + "' onclick='editReservation("
			+ i + ")'>" + reservations[i].start_time + "至" + reservations[i].end_time + "</div>");
		$("#col_teacher_fullname").append("<div class='table_cell' id='cell_teacher_fullname_"
			+ i + "'>" + reservations[i].teacher_fullname + "</div>");
		$("#col_teacher_username").append("<div class='table_cell' id='cell_teacher_username_"
			+ i + "'>" + reservations[i].teacher_username + "</div>");
		$("#col_teacher_mobile").append("<div class='table_cell' id='cell_teacher_mobile_"
			+ i + "'>" + reservations[i].teacher_mobile + "</div>");
		if (reservations[i].status === "AVAILABLE") {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i + "'>未预约</div>");
			$("#col_student").append("<div class='table_cell' id='cell_student_" + i + "'>" 
				+ "<button type='button' id='cell_student_view_" + i + "' onclick='setStudent(" + i + ")'>指定"
				+ "</button></div>");
		} else if (reservations[i].status === "RESERVATED") {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i + "'>已预约</div>");
			$("#col_student").append("<div class='table_cell' id='cell_student_" + i + "'>" 
				+ "<button type='button' id='cell_student_view_" + i + "' onclick='getStudent(" + i + ");'>查看"
				+ "</button></div>");
		} else if (reservations[i].status === "FEEDBACK") {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i + "'>"
				+ "<button type='button' id='cell_status_feedback_" + i + "' onclick='getFeedback(" + i + ");'>"
				+ "反馈</button></div>");
			$("#col_student").append("<div class='table_cell' id='cell_student_" + i + "'>" 
				+ "<button type='button' id='cell_student_view_" + i + "' onclick='getStudent(" + i + ");'>查看"
				+ "</button></div>");
		}
	}
	$("#col_select").append("<div class='table_cell' id='cell_select_add'><input type='checkbox'></div>");
	$("#col_time").append("<div class='table_cell' id='cell_time_add' onclick='addReservation();'>点击新增</div>");
	$("#col_teacher_fullname").append("<div class='table_cell' id='cell_teacher_fullname_add'></div>");
	$("#col_teacher_username").append("<div class='table_cell' id='cell_teacher_username_add'></div>");
	$("#col_teacher_mobile").append("<div class='table_cell' id='cell_teacher_mobile_add'></div>");
	$("#col_status").append("<div class='table_cell' id='cell_status_add'></div>");
	$("#col_student").append("<div class='table_cell' id='cell_student_add'></div>");
}

function optimize(t) {
	$("#col_select").width(40);
	$("#col_time").width(405);
	$("#col_teacher_fullname").width(120);
	$("#col_teacher_username").width(160);
	$("#col_teacher_mobile").width(160);
	$("#col_status").width(85);
	$("#col_student").width(85);
	// $('#col0').css('margin-left',width*0.02+'px')
	for (var i = 0; i < reservations.length; ++i) {
		var maxHeight = Math.max(
				$("#cell_select_" + i).height(),
				$("#cell_time_" + i).height(),
				$("#cell_teacher_fullname_" + i).height(),
				$("#cell_teacher_username_" + i).height(),
				$("#cell_teacher_mobile_" + i).height(),
				$("#cell_status_" + i).height(),
				$("#cell_student_" + i).height()
			);
		$("#cell_select_" + i).height(maxHeight);
		$("#cell_time_" + i).height(maxHeight);
		$("#cell_teacher_fullname_" + i).height(maxHeight);
		$("#cell_teacher_username_" + i).height(maxHeight);
		$("#cell_teacher_mobile_" + i).height(maxHeight);
		$("#cell_status_" + i).height(maxHeight);
		$("#cell_student_" + i).height(maxHeight);

		if (i % 2 == 1) {
			$("#cell_select_" + i).css("background-color", "white");
			$("#cell_time_" + i).css("background-color", "white");
			$("#cell_teacher_fullname_" + i).css("background-color", "white");
			$("#cell_teacher_username_" + i).css("background-color", "white");
			$("#cell_teacher_mobile_" + i).css("background-color", "white");
			$("#cell_status_" + i).css("background-color", "white");
			$("#cell_student_" + i).css("background-color", "white");
		} else {
			$("#cell_select_" + i).css("background-color", "#f3f3ff");
			$("#cell_time_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_fullname_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_username_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_mobile_" + i).css("background-color", "#f3f3ff");
			$("#cell_status_" + i).css("background-color", "#f3f3ff");
			$("#cell_student_" + i).css("background-color", "#f3f3ff");
		}
	}
	$("#cell_select_add").height(28);
	$("#cell_time_add").height(28);
	$("#cell_teacher_fullname_add").height(28);
	$("#cell_teacher_username_add").height(28);
	$("#cell_teacher_mobile_add").height(28);
	$("#cell_status_add").height(28);
	$("#cell_student_add").height(28);

	$(".table_head").height($("#head_select").height());
	$(t).css("left", (width - $(t).width()) / 2 - 11 + "px");
	$(t).css("top", (height - $(t).height()) / 2 - 11 + "px");
	$("#page_maintable").css("margin-left", 0.5 * ($(window).width()
		- (40 + 405 + 120 + 85 + 85 + 320)) + "px");
}

function addReservation() {
	$("#cell_time_add")[0].onclick = "";
	$("#cell_time_add")[0].innerHTML = "<input type='text' id='input_date_add' style='width: 80px'/>日，"
		+ "<input style='width:20px' id='start_hour_add'/>时<input style='width:20px' id='start_minute_add'/>分"
		+ "至<input style='width:20px' id='end_hour_add'/>时<input style='width:20px' id='end_minute_add'/>分";
	$("#cell_teacher_fullname_add")[0].innerHTML = "<input id='teacher_fullname_add' style='width:60px'/>"
		+ "<button type='button' onclick='searchTeacher();'>搜索</button>";
	$("#cell_teacher_username_add")[0].innerHTML = "<input id='teacher_username_add' style='width:120px'/>";
	$("#cell_teacher_mobile_add")[0].innerHTML = "<input id='teacher_mobile_add' style='width:120px'/>";
	$("#cell_status_add")[0].innerHTML = "<button type='button' onclick='addReservationConfirm();'>确认</button>";
	$("#cell_student_add")[0].innerHTML = "<button type='button' onclick='window.location.reload();'>取消</button>";
	$("#input_date_add").DatePicker({
		format: "YY-m-dd",
		date: $("#input_date_add").val(),
		current: $("#input_date_add").val(),
		starts: 1,
		position: "r",
		onBeforeShow: function() {
			$("#input_date_add").DatePickerSetDate($("#input_date_add").val(), true);
		},
		onChange: function(formated, dates) {
			$("#input_date_add").val(formated);
			$("#input_date_add").val($("#input_date_add").val().substr(4, 10));
			$("#input_date_add").DatePickerHide();
		},
	});
	optimize();
}

function addReservationConfirm() {
	var startHour = $("#start_hour_add").val();
	var startMinute = $("#start_minute_add").val();
	var endHour = $("#end_hour_add").val();
	var endMinute = $("#end_minute_add").val();
	var startTime = $("#input_date_add").val() + " " + (startHour.length < 2 ? "0" : "") + startHour + ":";
	if (startMinute.length == 0) {
		startTime += "00";
	} else if (startMinute.length == 1) {
		startTime += "0" + startMinute;
	} else {
		startTime += startMinute;
	}
	var endTime = $("#input_date_add").val() + " " + (endHour.length < 2 ? "0" : "") + endHour + ":";
	if (endMinute.length == 0) {
		endTime += "00";
	} else if (endMinute.length == 1) {
		endTime += "0" + endMinute;
	} else {
		endTime += endMinute;
	}
	var payload = {
		start_time: startTime,
		end_time: endTime,
		teacher_username: $("#teacher_username_add").val(),
		teacher_fullname: $("#teacher_fullname_add").val(),
		teacher_mobile: $("#teacher_mobile_add").val(),
	};
	console.log(payload);
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/reservation/add",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				viewReservations();
			} else {
				alert(data.message);
			}
		}
	});
}

function editReservation(index) {
	$("#cell_time_" + index)[0].onclick = "";
	$("#cell_time_" + index)[0].innerHTML = "<input type='text' id='input_date_" + index + "' style='width: 80px'/>日，"
		+ "<input style='width:20px' id='start_hour_" + index + "'/>时<input style='width:20px' id='start_minute_" + index + "'/>分"
		+ "至<input style='width:20px' id='end_hour_" + index + "'/>时<input style='width:20px' id='end_minute_" + index + "'/>分";
	$("#cell_teacher_fullname_" + index)[0].innerHTML = "<input id='teacher_fullname_" + index + "' style='width:60px' "
		+ "value='" + reservations[index].teacher_fullname + "''></input>"
		+ "<button type='button' onclick='searchTeacher();'>搜索</button>";
	$("#cell_teacher_username_" + index)[0].innerHTML = "<input id='teacher_username_" + index + "' style='width:120px' "
		+ "value='" + reservations[index].teacher_username + "'/>";
	$("#cell_teacher_mobile_" + index)[0].innerHTML = "<input id='teacher_mobile_" + index + "' style='width:120px' "
		+ "value='" + reservations[index].teacher_mobile + "'/>";
	$("#cell_status_" + index)[0].innerHTML = "<button type='button' onclick='editReservationConfirm(" + index + ");'>确认</button>";
	$("#cell_student_" + index)[0].innerHTML = "<button type='button' onclick='window.location.reload();'>取消</button>";
	$("#input_date_" + index).DatePicker({
		format: "YY-m-dd",
		date: $("#input_date_" + index).val(),
		current: $("#input_date_" + index).val(),
		starts: 1,
		position: "r",
		onBeforeShow: function() {
			$("#input_date_" + index).DatePickerSetDate($("#input_date_" + index).val(), true);
		},
		onChange: function(formated, dates) {
			$("#input_date_" + index).val(formated);
			$("#input_date_" + index).val($("#input_date_" + index).val().substr(4, 10));
			$("#input_date_" + index).DatePickerHide();
		},
	});
	optimize();
}

function editReservationConfirm(index) {
	var startHour = $("#start_hour_" + index).val();
	var startMinute = $("#start_minute_" + index).val();
	var endHour = $("#end_hour_" + index).val();
	var endMinute = $("#end_minute_" + index).val();
	var startTime = $("#input_date_" + index).val() + " " + (startHour.length < 2 ? "0" : "") + startHour + ":";
	if (startMinute.length == 0) {
		startTime += "00";
	} else if (startMinute.length == 1) {
		startTime += "0" + startMinute;
	} else {
		startTime += startMinute;
	}
	var endTime = $("#input_date_" + index).val() + " " + (endHour.length < 2 ? "0" : "") + endHour + ":";
	if (endMinute.length == 0) {
		endTime += "00";
	} else if (endMinute.length == 1) {
		endTime += "0" + endMinute;
	} else {
		endTime += endMinute;
	}
	var payload = {
		reservation_id: reservations[index].reservation_id,
		start_time: startTime,
		end_time: endTime,
		teacher_username: $("#teacher_username_" + index).val(),
		teacher_fullname: $("#teacher_fullname_" + index).val(),
		teacher_mobile: $("#teacher_mobile_" + index).val(),
	};
	console.log(payload);
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/reservation/edit",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				viewReservations();
			} else {
				alert(data.message);
			}
		}
	});
}

function searchTeacher(index) {
	var payload = {
		teacher_username: $("#teacher_username_" + (index === undefined ? "add" : index)).val(),
		teacher_fullname: $("#teacher_fullname_" + (index === undefined ? "add" : index)).val(),
		teacher_mobile: $("#teacher_mobile_" + (index === undefined ? "add" : index)).val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/teacher/search",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				$("#teacher_username_" + (index === undefined ? "add" : index)).val(data.teacher.teacher_username);
				$("#teacher_fullname_" + (index === undefined ? "add" : index)).val(data.teacher.teacher_fullname);
				$("#teacher_mobile_" + (index === undefined ? "add" : index)).val(data.teacher.teacher_mobile);
			}
		}
	});
}

function removeReservations() {
	$("body").append("\
		<div class='pop_window' style='width: 50%'>\
			确认删除选中的咨询记录？\
			<br>\
			<button type='button' onclick='$(\".pop_window\").remove();removeReservationsConfirm();'>确认</button>\
			<button type='button' onclick='$(\".pop_window\").remove();'>取消</button>\
		</div>\
	");
	optimize(".pop_window");
}

function removeReservationsConfirm() {
	var reservationIds = [];
	var sourceIds = [];
	var startTimes = [];
	for (var i = 0; i < reservations.length; ++i) {
		if ($("#cell_checkbox_" + i)[0].checked) {
			reservationIds.push(reservations[i].reservation_id);
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
		type: "POST",
		async: false,
		url: "/admin/reservation/remove",
		data: payload,
		traditional: true,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				viewReservations();
			} else {
				alert(data.message);
			}
		}
	});
}

function cancelReservations() {
	$("body").append("\
		<div class='pop_window' style='width: 50%'>\
			确认取消选中的预约？\
			<br>\
			<button type='button' onclick='$(\".pop_window\").remove();cancelReservationsConfirm();'>确认</button>\
			<button type='button' onclick='$(\".pop_window\").remove();'>取消</button>\
		</div>\
	");
	optimize(".pop_window");
}

function cancelReservationsConfirm() {
	var reservationIds = [];
	var sourceIds = [];
	var startTimes = [];
	for (var i = 0; i < reservations.length; ++i) {
		if ($("#cell_checkbox_" + i)[0].checked) {
			reservationIds.push(reservations[i].reservation_id);
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
		type: "POST",
		async: false,
		url: "/admin/reservation/cancel",
		data: payload,
		traditional: true,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				viewReservations();
			} else {
				alert(data.message);
			}
		}
	});
}

function getFeedback(index) {
	var payload = {
		reservation_id: reservations[index].reservation_id,
		source_id: reservations[index].source_id,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/reservation/feedback/get",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				showFeedback(index, data.feedback);
			} else {
				alert(data.message);
			}
		},
	});
}

function showFeedback(index, feedback) {
	$("body").append("\
		<div class='pop_window' id='feedback_table_" + index + "' style='text-align: left; width: 50%'>\
			咨询师反馈表<br>\
			评估分类：<br>\
			<select id='category_first_" + index + "' onchange='showSecondCategory(" + index + ")'><option value=''>请选择</option></select><br>\
			<select id='category_second_" + index + "'></select><br>\
			出席人员：<br>\
			<input id='participant_student_" + index + "' type='checkbox'>学生</input><input id='participant_parents_" + index + "' type='checkbox'>家长</input>\
			<input id='participant_teacher_" + index + "' type='checkbox'>教师</input><input id='participant_instructor_" + index + "' type='checkbox'>辅导员</input><br>\
			问题评估：<br>\
			<textarea id='problem_" + index + "' style='width: 100%; height:80px'></textarea><br>\
			咨询记录：<br>\
			<textarea id='record_" + index + "' style='width: 100%; height:80px'></textarea><br>\
			<button type='button' onclick='submitFeedback(" + index + ");'>提交</button>\
			<button type='button' onclick='$(\".pop_window\").remove();'>取消</button>\
		</div>\
	");
	getFeedbackCategories();
	showFirstCategory(index);
	if (feedback.category.length > 0) {
		$("#category_first_" + index).val(feedback.category.charAt(0));
		showSecondCategory(index);
		$("#category_second_" + index).val(feedback.category);
	}
	if (feedback.participants.length > 0) {
		$("#participant_student_" + index)[0].checked = feedback.participants[0] > 0;
		$("#participant_parents_" + index)[0].checked = feedback.participants[1] > 0;
		$("#participant_teacher_" + index)[0].checked = feedback.participants[2] > 0;
		$("#participant_instructor_" + index)[0].checked = feedback.participants[3] > 0;
	}
	$("#problem_" + index).val(feedback.problem);
	$("#record_" + index).val(feedback.record);
	optimize(".pop_window");
}

function getFeedbackCategories() {
	$.ajax({
		type: "GET",
		async: false,
		url: "/category/feedback",
		dateType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				firstCategory = data.first_category;
				secondCategory = data.second_category;
			}
		}
	});
}

function showFirstCategory(index) {
	for (var name in firstCategory) {
		if (firstCategory.hasOwnProperty(name)) {
			$("#category_first_" + index).append($('<option>', {
				value: name,
				text: firstCategory[name],
			}));
		}
	}
}

function showSecondCategory(index) {
	var first = $("#category_first_" + index).val();
	$("#category_second_" + index).find('option').remove().end().append("<option value=''>请选择</option>").val("");
	if ($("#category_first_" + index).selectedIndex === 0) {
		return;
	}
	if (secondCategory.hasOwnProperty(first)) {
		for (var name in secondCategory[first]) {
			if (secondCategory[first].hasOwnProperty(name)) {
				var option = new Option(name, secondCategory[first][name]);
				$("#category_second_" + index).append($('<option>', {
					value: name,
					text: secondCategory[first][name],
				}));
			}
		}
	}
}

function submitFeedback(index) {
	var participants = [];
	participants.push($("#participant_student_" + index)[0].checked ? 1 : 0);
	participants.push($("#participant_parents_" + index)[0].checked ? 1 : 0);
	participants.push($("#participant_teacher_" + index)[0].checked ? 1 : 0);
	participants.push($("#participant_instructor_" + index)[0].checked ? 1 : 0);
	var payload = {
		reservation_id: reservations[index].reservation_id,
		category: $("#category_second_" + index).val(),
		participants: participants,
		problem: $("#problem_" + index).val(),
		record: $("#record_" + index).val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/reservation/feedback/submit",
		data: payload,
		dataType: "json",
		traditional: true,
		success: function(data) {
			if (data.state === "SUCCESS") {
				successFeedback();
			} else {
				alert(data.message);
			}
		},
	});
}

function successFeedback() {
	$(".pop_window").remove();
	$("body").append("\
		<div class='pop_window' style='width: 50%;'>\
			您已成功提交反馈！<br>\
			<button type='button' onclick='$(\".pop_window\").remove();'>确定</button>\
		</div>\
	");
	optimize(".pop_window");
}

function setStudent(index) {
	$("body").append("\
		<div class='pop_window' style='width: 50%;'>\
			请输入您要制定的学生学号（必须为已注册学生）：<br>\
			<input id='student_username_" + index + "'/><br>\
			<button type='button' onclick='setStudentConfirm(" + index + ");'>确认</button>\
			<button type='button' style='margin-left:20px' onclick='$(\".pop_window\").remove();'>取消</button>\
		</div>\
	");
	optimize(".pop_window");
}

function setStudentConfirm(index) {
	var payload = {
		reservation_id: reservations[index].reservation_id,
		source_id: reservations[index].source_id,
		start_time: reservations[index].start_time,
		student_username: $("#student_username_" + index).val(),
	}
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/reservation/student/set",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state == "SUCCESS") {
				successSetStudent();
			} else {
				alert(data.message);
			}
		}
	});
}

function successSetStudent() {
	$(".pop_window").remove();
	$("body").append("\
		<div class='pop_window' style='width: 50%;'>\
			成功指定学生！<br>\
			<button type='button' onclick='$(\".pop_window\").remove();viewReservations();'>确定</button>\
		</div>\
	");
	optimize(".pop_window");
}

function getStudent(index) {
	var payload = {
		student_id: reservations[index].student_id,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/student/get",
		data: payload, 
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				showStudent(data.student_info, data.reservations);
			} else {
				alert(data.message);
			}
		}
	});
}

function showStudent(student, reservations) {
	$("body").append("\
		<div class='pop_window' style='text-align: left; height: 70%; overflow:auto;'>\
			学号：" + student.student_username + "<br>\
			姓名：" + student.student_fullname + "<br>\
			性别：" + student.student_gender + "<br>\
			出生日期：" + student.student_birthday + "<br>\
			系别：" + student.student_school + "<br>\
			年级：" + student.student_grade + "<br>\
			现住址：" + student.student_current_address + "<br>\
			家庭住址：" + student.student_family_address + "<br>\
			联系电话：" + student.student_mobile + "<br>\
			Email：" + student.student_email + "<br>\
			咨询经历：" + (student.student_experience_time ? "时间：" + student.student_experience_time + " 地点：" + student.student_experience_location + " 咨询师：" + student.student_experience_teacher : "无") + "<br>\
			父亲年龄：" + student.student_father_age + " 职业：" + student.student_father_job + " 学历：" + student.student_father_edu + "<br>\
			母亲年龄：" + student.student_mother_age + " 职业：" + student.student_mother_job + " 学历：" + student.student_mother_edu + "<br>\
			父母婚姻状况：" + student.student_parent_marriage + "<br>\
			近三个月里发生的有重大意义的事：" + student.student_significant + "<br>\
			需要接受帮助的主要问题：" + student.student_problem + "<br>\
			<br>\
			已绑定的咨询师：<span id='binded_teacher_username'>" + student.student_binded_teacher_username + "</span>&nbsp;\
				<span id='binded_teacher_fullname'>" + student.student_binded_teacher_fullname + "</span>\
				<button type='button' onclick='unbindStudent(\"" + student.student_id + "\");'>解绑</button><br>\
			请输入匹配咨询师工号：<input id='teacher_username' type='text'/>\
			<button type='button' onclick='bindStudent(\"" + student.student_id + "\");'>绑定</button><br>\
			<div style='margin: 10px 0'>\
				<button type='button' onclick='exportStudent(\"" + student.student_id + "\");'>导出</button>\
				<button type='button' onclick='$(\".pop_window\").remove();'>关闭</button>\
			</div>\
			<div id='student_reservations_" + student.student_id + "' style='width: 600px'>\
			</div>\
		</div>\
	");
	for (var i = 0; i < reservations.length; i++) {
		$("#student_reservations_" + student.student_id).append("\
			<div class='has_children' style='background: " + (reservations[i].status === "FEEDBACK" ? "#555" : "#F00") + "'>\
				<span>" + reservations[i].start_time + " 至 " + reservations[i].end_time + "  " + reservations[i].teacher_fullname + "</span>\
				<p class='children'>学生反馈：" + reservations[i].student_feedback.scores + "</p>\
				<p class='children'>评估分类：" + reservations[i].teacher_feedback.category + "</p>\
				<p class='children'>出席人员：" + reservations[i].teacher_feedback.participants + "</p>\
				<p class='children'>问题评估：" + reservations[i].teacher_feedback.problem + "</p>\
				<p class='children'>咨询记录：" + reservations[i].teacher_feedback.record + "</p>\
			</div>\
		");
	}
	$(document).ready(function() {
		$(".has_children").click(function() {
			$(this).addClass("highlight").children("p").show().end()
					.siblings().removeClass("highlight").children("p").hide();
		});
	});
	optimize(".pop_window");
}

function exportStudent(studentId) {
	var payload = {
		student_id: studentId,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/student/export",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				window.open(data.url);
			} else {
				alert(data.message);
			}
		},
	});
}

function unbindStudent(studentId) {
	var payload = {
		student_id: studentId,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/student/unbind",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				$("#binded_teacher_username").text(data.student_info.student_binded_teacher_username);
				$("#binded_teacher_fullname").text(data.student_info.student_binded_teacher_fullname);
			} else {
				alert(data.message);
			}
		},
	});
}

function bindStudent(studentId) {
	var payload = {
		student_id: studentId,
		teacher_username: $("#teacher_username").val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/student/bind",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				$("#binded_teacher_username").text(data.student_info.student_binded_teacher_username);
				$("#binded_teacher_fullname").text(data.student_info.student_binded_teacher_fullname);
			} else {
				alert(data.message);
			}
		},
	});
}

function getWorkload() {
	var payload = {
		from_date: $("#workload_from").val(),
		to_date: $("#workload_to").val(),
	}
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/teacher/workload",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				showWorkload(data.workload);
			} else {
				alert(data.message);
			}
		}
	});
}

function showWorkload(workload) {
	$("body").append("\
		<div class='pop_window' style='text-align: left; width: 50%; height: 70%; overflow: auto;'>\
			咨询师工作量统计\
			<div id='teacher_workload' style='width: 600px; margin-top: 10px;'>\
				<div class='table_col' id='col_workload_username'>\
					<div class='table_head table_cell'>咨询师工号</div>\
				</div>\
				<div class='table_col' id='col_workload_fullname'>\
					<div class='table_head table_cell'>咨询师姓名</div>\
				</div>\
				<div class='table_col' id='col_workload_student'>\
					<div class='table_head table_cell'>咨询人数</div>\
				</div>\
				<div class='table_col' id='col_workload_reservation'>\
					<div class='table_head table_cell'>咨询人次</div>\
				</div>\
				<div class='clearfix'></div>\
			</div>\
			<div style='margin: 10px 0'>\
				<button type='button' onclick='$(\".pop_window\").remove();'>关闭</button>\
			</div>\
		</div>\
	");
	$("#col_workload_username").width(100);
	$("#col_workload_fullname").width(100);
	$("#col_workload_student").width(80);
	$("#col_workload_reservation").width(80);
	for (var i in workload) {
		if (workload.hasOwnProperty(i)) {
			$("#col_workload_username").append("<div class='table_cell' id='cell_workload_username_"
				+ i + "'>" + workload[i].teacher_username + "</div>");
			$("#col_workload_fullname").append("<div class='table_cell' id='cell_workload_fullname_"
				+ i + "'>" + workload[i].teacher_fullname + "</div>");
			$("#col_workload_student").append("<div class='table_cell' id='cell_workload_student_"
				+ i + "'>" + Object.size(workload[i].students) + "</div>");
			$("#col_workload_reservation").append("<div class='table_cell' id='cell_workload_reservation_"
				+ i + "'>" + Object.size(workload[i].reservations) + "</div>");
		}
	}
	optimize(".pop_window");
}

Object.size = function(obj) {
	var size = 0, key;
	for (key in obj) {
		if (obj.hasOwnProperty(key)) size++;
	}
	return size;
}

function exportMonthlyReport() {
	var payload = {
		monthly_date: $("#monthly_report_date").val(),
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/admin/reservation/export/report/monthly",
		data: payload,
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				window.open(data.url);
			} else {
				alert(data.message);
			}
		}
	});
}

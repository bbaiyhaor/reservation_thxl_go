var width=$(window).width();
var height=$(window).height();
var teacher;
var reservations;
var firstCategory;
var secondCategory;

function viewReservations() {
	$.ajax({
		type: "GET",
		async: false,
		url: "/teacher/reservation/view",
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				console.log(data);
				reservations = data.reservations;
				teacher = data.teacher_info;
				refreshDataTable(reservations);
				optimize();
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
				<button onclick='$(\".checkbox\").click();' style='padding:0px;'>全选</button>\
			</div>\
		</div>\
		<div class='table_col' id='col_time'>\
			<div class='table_head table_cell'>时间</div>\
		</div>\
		<div class='table_col' id='col_teacher_fullname'>\
			<div class='table_head table_cell'>咨询师</div>\
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
			+ i + ")'>" + reservations[i].start_time.split(" ")[0].substr(2) + "<br>" 
			+ reservations[i].start_time.split(" ")[1] + "-" + reservations[i].end_time.split(" ")[1] + "</div>");
		$("#col_teacher_fullname").append("<div class='table_cell' id='cell_teacher_fullname_"
			+ i + "'>" + reservations[i].teacher_fullname + "</div>");
		$("#col_teacher_mobile").append("<div class='table_cell' id='cell_teacher_mobile_"
			+ i + "'>" + reservations[i].teacher_mobile + "</div>");
		if (reservations[i].status === "AVAILABLE") {
			$("#col_status").append("<div class='table_cell' id='cell_status_" + i + "'>未预约</div>");
			$("#col_student").append("<div class='table_cell' id='cell_student_" + i + "'>" 
				+ "<button type='button' id='cell_student_view_" + i + "' disabled='true'>查看"
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
}

function optimize(t) {
	$("#col_select").width(20);
	$("#col_time").width(80);
	$("#col_teacher_fullname").width(44);
	$("#col_teacher_mobile").width(76);
	$("#col_status").width(40);
	$("#col_student").width(40);
	$("#col_select").css("margin-left", (width - 300) / 2 + "px");
	for (var i = 0; i < reservations.length; ++i) {
		var maxHeight = Math.max(
				$("#cell_select_" + i).height(),
				$("#cell_time_" + i).height(),
				$("#cell_teacher_fullname_" + i).height(),
				$("#cell_teacher_mobile_" + i).height(),
				$("#cell_status_" + i).height(),
				$("#cell_student_" + i).height()
			);
		$("#cell_select_" + i).height(maxHeight);
		$("#cell_time_" + i).height(maxHeight);
		$("#cell_teacher_fullname_" + i).height(maxHeight);
		$("#cell_teacher_mobile_" + i).height(maxHeight);
		$("#cell_status_" + i).height(maxHeight);
		$("#cell_student_" + i).height(maxHeight);

		if (i % 2 == 1) {
			$("#cell_select_" + i).css("background-color", "white");
			$("#cell_time_" + i).css("background-color", "white");
			$("#cell_teacher_fullname_" + i).css("background-color", "white");
			$("#cell_teacher_mobile_" + i).css("background-color", "white");
			$("#cell_status_" + i).css("background-color", "white");
			$("#cell_student_" + i).css("background-color", "white");
		} else {
			$("#cell_select_" + i).css("background-color", "#f3f3ff");
			$("#cell_time_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_fullname_" + i).css("background-color", "#f3f3ff");
			$("#cell_teacher_mobile_" + i).css("background-color", "#f3f3ff");
			$("#cell_status_" + i).css("background-color", "#f3f3ff");
			$("#cell_student_" + i).css("background-color", "#f3f3ff");
		}
	}
	$(t).css("left", (width - $(t).width()) / 2 - 11 + "px");
	$(t).css("top", (height - $(t).height()) / 2 - 11 + "px");
	$(".table_head").height($("#head_select").height());
}

function getFeedback(index) {
	var payload = {
		reservation_id: reservations[index].reservation_id,
		source_id: reservations[index].source_id,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/teacher/reservation/feedback/get",
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
		<div class='fankui_tch' id='feedback_table_" + index + "' style='font-size:11px;text-align:left;top:100px;height:400;left:5px'>\
			咨询师反馈表<br>\
			评估分类：<br>\
			<select id='category_first_" + index + "' onchange='showSecondCategory(" + index + ")'><option value=''>请选择</option></select><br>\
			<select id='category_second_" + index + "'></select><br>\
			出席人员：<br>\
			<input id='participant_student_" + index + "' type='checkbox'>学生</input><input id='participant_parents_" + index + "' type='checkbox'>家长</input>\
			<input id='participant_teacher_" + index + "' type='checkbox'>教师</input><input id='participant_instructor_" + index + "' type='checkbox'>辅导员</input><br>\
			问题评估：<br>\
			<textarea id='problem_" + index + "' style='width:180px;height:80px'></textarea><br>\
			咨询记录：<br>\
			<textarea id='record_" + index + "' style='width:180px;height:80px'></textarea><br>\
			<button type='button' onclick='submitFeedback(" + index + ");'>提交</button>\
			<button type='button' onclick='$(\".fankui_tch\").remove();'>取消</button>\
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
	optimize(".fankui_tch");
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
		url: "/teacher/reservation/feedback/submit",
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
	$(".fankui_tch").remove();
	$("body").append("\
		<div class='fankui_tch_success'>\
			您已成功提交反馈！<br>\
			<button type='button' onclick='$(\".fankui_tch_success\").remove();'>确定</button>\
		</div>\
	");
	optimize(".fankui_tch_success");
}

function getStudent(index) {
	var payload = {
		reservation_id: reservations[index].reservation_id,
	};
	$.ajax({
		type: "POST",
		async: false,
		url: "/teacher/student/get",
		data: payload, 
		dataType: "json",
		success: function(data) {
			if (data.state === "SUCCESS") {
				showStudent(data.student_info);
			} else {
				alert(data.message);
			}
		}
	});
}

function showStudent(student) {
	$("body").append("\
		<div class='admin_chakan' style='text-align: left'>\
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
				<span id='binded_teacher_fullname'>" + student.student_binded_teacher_fullname + "</span><br>\
			<br>\
			<button type='button' onclick='$(\".admin_chakan\").remove();'>关闭</button>\
		</div>\
	");
	optimize(".admin_chakan");
}
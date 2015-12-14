var width = $(window).width();
var height = $(window).height();
var timedReservations;

function viewTimedReservations() {
    $.ajax({
        type: "GET",
        async: false,
        url: "/admin/timetable/view",
        dataType: "json",
        success: function(data) {
            if (data.state === "SUCCESS") {
                console.log(data);
                timedReservations = data.timed_reservations;
                for (var weekday in timedReservations) {
                    if (timedReservations.hasOwnProperty(weekday)) {
                        refreshDataTable(weekday, timedReservations[weekday])
                        optimize(weekday);
                    }
                }
            } else {
                alert(data.message);
            }
        }
    });
}

function refreshDataTable(weekday, timedReservations) {
    $("#page_maintable_" + weekday)[0].innerHTML = "\
		<div class='table_col' id='col_select_" + weekday + "'>\
			<div class='table_head table_cell' id='head_select_" + weekday + "'>\
				<button onclick='$(\".checkbox\").click();'>全选</button>\
			</div>\
		</div>\
		<div class='table_col' id='col_time_" + weekday + "'>\
			<div class='table_head table_cell'>时间</div>\
		</div>\
		<div class='table_col' id='col_teacher_fullname_" + weekday + "'>\
			<div class='table_head table_cell'>咨询师</div>\
		</div>\
		<div class='table_col' id='col_teacher_username_" + weekday + "'>\
			<div class='table_head table_cell'>咨询师编号</div>\
		</div>\
		<div class='table_col' id='col_teacher_mobile_" + weekday + "'>\
			<div class='table_head table_cell'>咨询师手机</div>\
		</div>\
		<div class='table_col' id='col_status_" + weekday + "'>\
			<div class='table_head table_cell'>状态</div>\
		</div>\
		<div class='table_col' id='col_operation_" + weekday + "'>\
			<div class='table_head table_cell'>操作</div>\
		</div>\
		<div class='clearfix'></div>\
	";

    for (var i = 0; i < timedReservations.length; ++i) {
        $("#col_select_" + weekday).append("<div class='table_cell' id='cell_select_" + weekday + "_" + i + "'>"
            + "<input class='checkbox' type='checkbox' id='cell_checkbox_" + weekday + "_" + i + "'></div>");
        $("#col_time_" + weekday).append("<div class='table_cell' id='cell_time_" + weekday + "_" + i + "' onclick='editTimedReservation(\""
            + weekday + "\"," + i + ")'>" + timedReservations[i].start_clock + "　-　" + timedReservations[i].end_clock + "</div>");
        $("#col_teacher_fullname_" + weekday).append("<div class='table_cell' id='cell_teacher_fullname_" + weekday + "_"
            + i + "'>" + timedReservations[i].teacher_fullname + "</div>");
        $("#col_teacher_username_" + weekday).append("<div class='table_cell' id='cell_teacher_username_" + weekday + "_"
            + i + "'>" + timedReservations[i].teacher_username + "</div>");
        $("#col_teacher_mobile_" + weekday).append("<div class='table_cell' id='cell_teacher_mobile_" + weekday + "_"
            + i + "'>" + timedReservations[i].teacher_mobile + "</div>");
        if (timedReservations[i].status === "AVAILABLE") {
            $("#col_status_" + weekday).append("<div class='table_cell' id='cell_status_" + weekday + "_" + i + "'>生效</div>");
            $("#col_operation_" + weekday).append("<div class='table_cell' id='cell_operation_" + weekday + "_" + + i + "'>"
                + "<button type='button' id='cell_operation_oper_" + weekday + "_" + + i + "' closeTimedReservation='(" + i + ");' disabled='true'>关闭"
                + "</button></div>");
        } else {
            $("#col_status_" + weekday).append("<div class='table_cell' id='cell_status_" + weekday + "_" + i + "'>关闭</div>");
            $("#col_operation_" + weekday).append("<div class='table_cell' id='cell_operation_" + weekday + "_" + + i + "'>"
                + "<button type='button' id='cell_operation_oper_" + weekday + "_" + + i + "' openTimedReservation='(" + i + ");' disabled='true'>打开"
                + "</button></div>");
        }
    }
    $("#col_select_" + weekday).append("<div class='table_cell' id='cell_select_" + weekday + "_add'><input type='checkbox'></div>");
    $("#col_time_" + weekday).append("<div class='table_cell' id='cell_time_" + weekday + "_add' onclick='addTimedReservation(\"" + weekday + "\");'>点击新增</div>");
    $("#col_teacher_fullname_" + weekday).append("<div class='table_cell' id='cell_teacher_fullname_" + weekday + "_add'></div>");
    $("#col_teacher_username_" + weekday).append("<div class='table_cell' id='cell_teacher_username_" + weekday + "_add'></div>");
    $("#col_teacher_mobile_" + weekday).append("<div class='table_cell' id='cell_teacher_mobile_" + weekday + "_add'></div>");
    $("#col_status_" + weekday).append("<div class='table_cell' id='cell_status_" + weekday + "_add'></div>");
    $("#col_operation_" + weekday).append("<div class='table_cell' id='cell_operation_" + weekday + "_add'></div>");
}

function optimize(weekday, t) {
    $("#col_select_" + weekday).width(40);
    $("#col_time_" + weekday).width(300);
    $("#col_teacher_fullname_" + weekday).width(120);
    $("#col_teacher_username_" + weekday).width(160);
    $("#col_teacher_mobile_" + weekday).width(160);
    $("#col_status_" + weekday).width(85);
    $("#col_operation_" + weekday).width(85);
    // $('#col0').css('margin-left',width*0.02+'px')
    for (var i = 0; i < timedReservations[weekday].length; ++i) {
        var maxHeight = Math.max(
            $("#cell_select_" + weekday + "_" + i).height(),
            $("#cell_time_" + weekday + "_" + i).height(),
            $("#cell_teacher_fullname_" + weekday + "_" + i).height(),
            $("#cell_teacher_username_" + weekday + "_" + i).height(),
            $("#cell_teacher_mobile_" + weekday + "_" + i).height(),
            $("#cell_status_" + weekday + "_" + i).height(),
            $("#cell_operation_" + weekday + "_" + i).height()
        );
        $("#cell_select_" + weekday + "_" + i).height(maxHeight);
        $("#cell_time_" + weekday + "_" + i).height(maxHeight);
        $("#cell_teacher_fullname_" + weekday + "_" + i).height(maxHeight);
        $("#cell_teacher_username_" + weekday + "_" + i).height(maxHeight);
        $("#cell_teacher_mobile_" + weekday + "_" + i).height(maxHeight);
        $("#cell_status_" + weekday + "_" + i).height(maxHeight);
        $("#cell_operation_" + weekday + "_" + i).height(maxHeight);

        if (i % 2 == 1) {
            $("#cell_select_" + weekday + "_" + i).css("background-color", "white");
            $("#cell_time_" + weekday + "_" + i).css("background-color", "white");
            $("#cell_teacher_fullname_" + weekday + "_" + i).css("background-color", "white");
            $("#cell_teacher_username_" + weekday + "_" + i).css("background-color", "white");
            $("#cell_teacher_mobile_" + weekday + "_" + i).css("background-color", "white");
            $("#cell_status_" + weekday + "_" + i).css("background-color", "white");
            $("#cell_operation_" + weekday + "_" + i).css("background-color", "white");
        } else {
            $("#cell_select_" + weekday + "_" + i).css("background-color", "#f3f3ff");
            $("#cell_time_" + weekday + "_" + i).css("background-color", "#f3f3ff");
            $("#cell_teacher_fullname_" + weekday + "_" + i).css("background-color", "#f3f3ff");
            $("#cell_teacher_username_" + weekday + "_" + i).css("background-color", "#f3f3ff");
            $("#cell_teacher_mobile_" + weekday + "_" + i).css("background-color", "#f3f3ff");
            $("#cell_status_" + weekday + "_" + i).css("background-color", "#f3f3ff");
            $("#cell_operation_" + weekday + "_" + i).css("background-color", "#f3f3ff");
        }
    }
    $("#cell_select_" + weekday + "_add").height(28);
    $("#cell_time_" + weekday + "_add").height(28);
    $("#cell_teacher_fullname_" + weekday + "_add").height(28);
    $("#cell_teacher_username_" + weekday + "_add").height(28);
    $("#cell_teacher_mobile_" + weekday + "_add").height(28);
    $("#cell_status_" + weekday + "_add").height(28);
    $("#cell_operation_" + weekday + "_add").height(28);

    $(".table_head").height($("#head_select_" + weekday).height());
    $(t).css("left", (width - $(t).width()) / 2 - 11 + "px");
    $(t).css("top", (height - $(t).height()) / 2 - 11 + "px");
    $("#page_maintable_" + weekday).css("margin-left", 0.5 * (width - (40 + 300 + 120 + 85 + 85 + 320)) + "px");
}

function addTimedReservation(weekday) {
    $("#cell_time_" + weekday + "_add")[0].onclick = "";
    $("#cell_time_" + weekday + "_add")[0].innerHTML = "<input style='width:20px' id='start_hour_" + weekday + "_add'/>时" +
        "<input style='width:20px' id='start_minute_" + weekday + "_add'/>分　-　" +
        "<input style='width:20px' id='end_hour_" + weekday + "_add'/>时" +
        "<input style='width:20px' id='end_minute_" + weekday + "_add'/>分";
    $("#cell_teacher_fullname_" + weekday + "_add")[0].innerHTML = "<input id='teacher_fullname_" + weekday + "_add' style='width:60px'/>"
        + "<button type='button' onclick='searchTeacher(\"" + weekday + "\");'>搜索</button>";
    $("#cell_teacher_username_" + weekday + "_add")[0].innerHTML = "<input id='teacher_username_" + weekday + "_add' style='width:120px'/>";
    $("#cell_teacher_mobile_" + weekday + "_add")[0].innerHTML = "<input id='teacher_mobile_" + weekday + "_add' style='width:120px'/>";
    $("#cell_status_" + weekday + "_add")[0].innerHTML = "<button type='button' onclick='addTimedReservationConfirm(\"" + weekday + "\");'>确认</button>";
    $("#cell_operation_" + weekday + "_add")[0].innerHTML = "<button type='button' onclick='window.location.reload();'>取消</button>";
    optimize(weekday);
}

function addTimedReservationConfirm(weekday) {
    var startHour = $("#start_hour_" + weekday + "_add").val();
    var startMinute = $("#start_minute_" + weekday + "_add").val();
    var startClock = (startHour.length < 2 ? "0" : "") + startHour + ":";
    if (startMinute.length == 0) {
        startClock += "00";
    } else if (startMinute.length == 1) {
        startClock += "0" + startMinute;
    } else {
        startClock += startMinute;
    }
    var endHour = $("#end_hour_" + weekday + "_add").val();
    var endMinute = $("#end_minute_" + weekday + "_add").val();
    var endClock = (endHour.length < 2 ? "0" : "") + endHour + ":";
    if (endMinute.length == 0) {
        endClock += "00";
    } else if (endMinute.length == 1) {
        endClock += "0" + endMinute;
    } else {
        endClock += endMinute;
    }
    var payload = {
        weekday: weekday,
        start_clock: startClock,
        end_clock: endClock,
        teacher_username: $("#teacher_username_" + weekday + "_add").val(),
        teacher_fullname: $("#teacher_fullname_" + weekday + "_add").val(),
        teacher_mobile: $("#teacher_mobile_" + weekday + "_add").val(),
    };
    $.ajax({
        type: "POST",
        async: false,
        url: "/admin/timetable/add",
        data: payload,
        dataType: "json",
        success: function(data) {
            if (data.state === "SUCCESS") {
                viewTimedReservations();
            } else {
                alert(data.message);
            }
        }
    });
}

function editTimedReservation(weekday, index) {
    $("#cell_time_" + weekday + "_" + index)[0].onclick = "";
    $("#cell_time_" + weekday + "_" + index)[0].innerHTML = "<input style='width:20px' id='start_hour_" + weekday + "_" + index + "'/>时" +
        "<input style='width:20px' id='start_minute_" + weekday + "_" + index + "'/>分　-　" +
        "<input style='width:20px' id='end_hour_" + weekday + "_" + index + "'/>时" +
        "<input style='width:20px' id='end_minute_" + weekday + "_" + index + "'/>分";
    $("#cell_teacher_fullname_" + weekday + "_" + index)[0].innerHTML = "<input id='teacher_fullname_" + weekday + "_" + index + "' style='width:60px' "
        + "value='" + timedReservations[weekday][index].teacher_fullname + "''></input>"
        + "<button type='button' onclick='searchTeacher(\"" + weekday + "\"," + index + ");'>搜索</button>";
    $("#cell_teacher_username_" + weekday + "_" + index)[0].innerHTML = "<input id='teacher_username_" + weekday + "_" + index + "' style='width:120px' "
        + "value='" + timedReservations[weekday][index].teacher_username + "'/>";
    $("#cell_teacher_mobile_" + weekday + "_" + index)[0].innerHTML = "<input id='teacher_mobile_" + weekday + "_" + index + "' style='width:120px' "
        + "value='" + timedReservations[weekday][index].teacher_mobile + "'/>";
    $("#cell_status_" + weekday + "_" + index)[0].innerHTML = "<button type='button' onclick='editTimedReservationConfirm(\"" + weekday + "\"," + index + ");'>确认</button>";
    $("#cell_operation_" + weekday + "_" + index)[0].innerHTML = "<button type='button' onclick='window.location.reload();'>取消</button>";
    optimize(weekday);
}

function editTimedReservationConfirm(weekday, index) {
    var startHour = $("#start_hour_" + weekday + "_" + index).val();
    var startMinute = $("#start_minute_" + weekday + "_" + index).val();
    var startClock = (startHour.length < 2 ? "0" : "") + startHour + ":";
    if (startMinute.length == 0) {
        startClock += "00";
    } else if (startMinute.length == 1) {
        startClock += "0" + startMinute;
    } else {
        startClock += startMinute;
    }
    var endHour = $("#end_hour_" + weekday + "_" + index).val();
    var endMinute = $("#end_minute_" + weekday + "_" + index).val();
    var endClock = (endHour.length < 2 ? "0" : "") + endHour + ":";
    if (endMinute.length == 0) {
        endClock += "00";
    } else if (endMinute.length == 1) {
        endClock += "0" + endMinute;
    } else {
        endClock += endMinute;
    }
    var payload = {
        timed_reservation_id: timedReservations[weekday][index].timed_reservation_id,
        weekday: weekday,
        start_clock: startClock,
        end_clock: endClock,
        teacher_username: $("#teacher_username_" + weekday + "_" + index).val(),
        teacher_fullname: $("#teacher_fullname_" + weekday + "_" + index).val(),
        teacher_mobile: $("#teacher_mobile_" + weekday + "_" + index).val(),
    };
    $.ajax({
        type: "POST",
        async: false,
        url: "/admin/timetable/edit",
        data: payload,
        dataType: "json",
        success: function(data) {
            if (data.state === "SUCCESS") {
                viewTimedReservations();
            } else {
                alert(data.message);
            }
        }
    });
}

function searchTeacher(weekday, index) {
    var payload = {
        teacher_username: $("#teacher_username_" + weekday + "_" + (index === undefined ? "add" : index)).val(),
        teacher_fullname: $("#teacher_fullname_" + weekday + "_" + (index === undefined ? "add" : index)).val(),
        teacher_mobile: $("#teacher_mobile_" + weekday + "_" + (index === undefined ? "add" : index)).val(),
    };
    $.ajax({
        type: "POST",
        async: false,
        url: "/admin/teacher/search",
        data: payload,
        dataType: "json",
        success: function(data) {
            if (data.state === "SUCCESS") {
                $("#teacher_username_" + weekday + "_" + (index === undefined ? "add" : index)).val(data.teacher.teacher_username);
                $("#teacher_fullname_" + weekday + "_" + (index === undefined ? "add" : index)).val(data.teacher.teacher_fullname);
                $("#teacher_mobile_" + weekday + "_" + (index === undefined ? "add" : index)).val(data.teacher.teacher_mobile);
            }
        }
    });
}

function removeTimedReservations() {
    $("body").append("\
		<div class='delete_admin_pre'>\
			确认删除选中的预设咨询？\
			<br>\
			<button type='button' onclick='$(\".delete_admin_pre\").remove();removeReservationsConfirm();'>确认</button>\
			<button type='button' onclick='$(\".delete_admin_pre\").remove();'>取消</button>\
		</div>\
	");
    optimize("Monday", ".delete_admin_pre");
}

function removeReservationsConfirm() {
    var timedReservationIds = [];
    for (var weekday in timedReservations) {
        if (timedReservations.hasOwnProperty(weekday)) {
            for (var i = 0; i < timedReservations[weekday].length; ++i) {
                if ($("#cell_checkbox_" + weekday + "_" + i)[0].checked) {
                    timedReservationIds.push(timedReservations[weekday][i].timed_reservation_id);
                }
            }
        }
    }
    var payload = {
        timed_reservation_ids: timedReservationIds,
    };
    $.ajax({
        type: "POST",
        async: false,
        url: "/admin/timetable/remove",
        data: payload,
        traditional: true,
        dataType: "json",
        success: function(data) {
            if (data.state === "SUCCESS") {
                viewTimedReservations();
            } else {
                alert(data.message);
            }
        }
    });
}
package rerror

import "fmt"

const (
	OK    = 0
	CHECK = 1
	// 请求类错误
	ERROR_MISSING_PARAM = 10
	ERROR_INVALID_PARAM = 11
	// 格式类错误
	ERROR_FORMAT_MOBILE    = 51
	ERROR_FORMAT_EMAIL     = 52
	ERROR_FORMAT_WEEKDAY   = 53
	ERROR_FORMAT_STUDENTID = 54
	// 账户类错误
	ERROR_EXPIRE_SESSION                      = 100
	ERROR_NO_LOGIN                            = 101
	ERROR_NOT_AUTHORIZED                      = 102
	ERROR_LOGIN_PASSWORD_WRONG                = 103
	ERROR_LOGIN_PWDCHANGE_OLDPWD_MISMATCH     = 104
	ERROR_NO_USER                             = 105
	ERROR_EXIST_USERNAME                      = 106
	ERROR_LOGIN_PWDCHANGE_OLDPWD_EQUAL_NEWPED = 107
	ERROR_LOGIN_PWDCHANGE_USERNAME_MISMATCH   = 108
	ERROR_LOGIN_PWDCHANGE_MOBILE_MISMATCh     = 109
	ERROR_LOGIN_PWDCHANGE_VERIFY_CODE_WRONG   = 110
	ERROR_NO_STUDENT                          = 111
	// 通用逻辑类错误
	ERROR_FEEDBACK_AVAILABLE_RESERVATION = 201
	ERROR_FEEDBACK_FUTURE_RESERVATION    = 202
	ERROR_FEEDBACK_OTHER_RESERVATION     = 203
	ERROR_START_TIME_MISMATCH            = 204
	ERROR_SEND_SMS                       = 205
	// 管理员错误
	ERROR_ADMIN_EDIT_RESERVATION_END_TIME_BEFORE_START_TIME = 301
	ERROR_ADMIN_EDIT_RESERVATION_TEACHER_TIME_CONFLICT      = 302
	ERROR_ADMIN_EDIT_TIMETABLE_IN_RESERVATION               = 303
	ERROR_ADMIN_EDIT_RESERVATED_RESERVATION                 = 304
	ERROR_ADMIN_EDIT_RESERVATION_OUTDATED                   = 305
	ERROR_ADMIN_SET_RESERVATED_RESERVATION                  = 306
	ERROR_ADMIN_ARCHIVE_NUMBER_ALREADY_EXIST                = 308
	ERROR_ADMIN_EXPORT_STUDENT_NO_ARCHIVE_NUMBER            = 309
	ERROR_ADMIN_NO_RESERVATIONS_TODAY                       = 310
	// 学生错误
	ERROR_STUDENT_ALREADY_HAVE_RESERVATION            = 401
	ERROR_STUDENT_MAKE_OUTDATED_RESERVATION           = 402
	ERROR_STUDENT_MAKE_RESERVATED_RESERVATION         = 403
	ERROR_STUDENT_MAKE_NOT_BINDED_TEACHER_RESERVATION = 404
	// 咨询师错误
	ERROR_TEACHER_VIEW_OTHER_STUDENT = 501
	// 数据库类错误
	ERROR_DATABASE   = 1000
	ERROR_ID_INVALID = 1001

	ERROR_UNKNOWN = -1
)

func ReturnMessage(code int, args ...interface{}) string {
	switch code {
	case OK:
		return "OK"
	case CHECK:
		return "CHECK"
	case ERROR_MISSING_PARAM:
		return fmt.Sprintf("参数缺失：%s", args...)
	case ERROR_INVALID_PARAM:
		return fmt.Sprintf("参数格式错误：%s", args...)
	case ERROR_EXPIRE_SESSION:
		return "会话过期，请重新登录"
	case ERROR_NO_LOGIN:
		return "请先登录"
	case ERROR_NOT_AUTHORIZED:
		return "权限不足"
	case ERROR_LOGIN_PASSWORD_WRONG:
		return "用户名或密码不正确"
	case ERROR_LOGIN_PWDCHANGE_OLDPWD_MISMATCH:
		return "旧登录密码不正确"
	case ERROR_NO_USER:
		return "未找到用户"
	case ERROR_EXIST_USERNAME:
		return "用户名已被注册"
	case ERROR_LOGIN_PWDCHANGE_OLDPWD_EQUAL_NEWPED:
		return "新密码不能与原有密码一样"
	case ERROR_LOGIN_PWDCHANGE_VERIFY_CODE_WRONG:
		return "验证码错误或已过期"
	case ERROR_LOGIN_PWDCHANGE_USERNAME_MISMATCH:
		return "用户名不匹配"
	case ERROR_LOGIN_PWDCHANGE_MOBILE_MISMATCh:
		return "手机号不匹配"
	case ERROR_NO_STUDENT:
		return "学生未注册"
	case ERROR_FORMAT_MOBILE:
		return "手机号格式不正确"
	case ERROR_FORMAT_EMAIL:
		return "邮箱格式不正确"
	case ERROR_FORMAT_WEEKDAY:
		return "星期格式不正确"
	case ERROR_FORMAT_STUDENTID:
		return "学号格式不正确"
	case ERROR_SEND_SMS:
		return "发送短信失败，请稍后重试"
	case ERROR_ADMIN_EDIT_RESERVATION_END_TIME_BEFORE_START_TIME:
		return "开始时间不能晚于结束时间"
	case ERROR_ADMIN_EDIT_RESERVATION_TEACHER_TIME_CONFLICT:
		return "咨询师时间有冲突"
	case ERROR_ADMIN_EDIT_TIMETABLE_IN_RESERVATION:
		return "请在预设安排表中编辑预设咨询"
	case ERROR_ADMIN_EDIT_RESERVATED_RESERVATION:
		return "不能编辑已预约的咨询"
	case ERROR_ADMIN_EDIT_RESERVATION_OUTDATED:
		return "不能编辑已过期咨询"
	case ERROR_FEEDBACK_AVAILABLE_RESERVATION:
		return "不能反馈未被预约的咨询"
	case ERROR_FEEDBACK_FUTURE_RESERVATION:
		return "不能反馈还未开始的咨询"
	case ERROR_ADMIN_SET_RESERVATED_RESERVATION:
		return "不能设定已被预约的咨询"
	case ERROR_START_TIME_MISMATCH:
		return "开始时间不匹配"
	case ERROR_ADMIN_ARCHIVE_NUMBER_ALREADY_EXIST:
		return "档案号已存在，请重新分配"
	case ERROR_ADMIN_EXPORT_STUDENT_NO_ARCHIVE_NUMBER:
		return "请先分配档案号"
	case ERROR_ADMIN_NO_RESERVATIONS_TODAY:
		return "今日无咨询"
	case ERROR_STUDENT_ALREADY_HAVE_RESERVATION:
		return "你好！你已有一个咨询预约，请完成这次咨询后再预约下一次，或致电62782007取消已有预约。"
	case ERROR_STUDENT_MAKE_OUTDATED_RESERVATION:
		return "不能预约已过期咨询"
	case ERROR_STUDENT_MAKE_RESERVATED_RESERVATION:
		return "不能预约已被预约的咨询"
	case ERROR_STUDENT_MAKE_NOT_BINDED_TEACHER_RESERVATION:
		return "只能预约匹配咨询师的咨询"
	case ERROR_FEEDBACK_OTHER_RESERVATION:
		return "只能反馈自己预约的咨询"
	case ERROR_TEACHER_VIEW_OTHER_STUDENT:
		return "只能查看本人绑定的学生"
	case ERROR_DATABASE:
		return "获取数据失败"
	case ERROR_ID_INVALID:
		return "ID值不合法"
	case ERROR_UNKNOWN:
		return "未知错误"
	default:
		return "服务器开小差了，请稍后重试~"
	}
}

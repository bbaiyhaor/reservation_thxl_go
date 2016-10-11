package workflow

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxl_go/config"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	SMS_SUCCESS_STUDENT  = "%s你好，你已成功预约星期%s（%d月%d日）%s-%s咨询，地点：紫荆C楼409室。电话：62782007。"
	SMS_SUCCESS_TEACHER  = "%s您好，%s已预约您星期%s（%d月%d日）%s-%s咨询，地点：紫荆C楼409室。电话：62782007。"
	SMS_CANCEL_STUDENT   = "%s你好，因特殊原因你预约的星期%s（%d月%d日）%s-%s咨询已被取消，详询62782007。"
	SMS_CANCEL_TEACHER   = "%s您好，%s预约的星期%s（%d月%d日）%s-%s咨询已被取消，详询62782007。"
	SMS_REMINDER_STUDENT = "温馨提示：%s你好，你已成功预约明天%s-%s咨询，地点：紫荆C楼409室。电话：62782007。"
	SMS_REMINDER_TEACHER = "温馨提示：%s您好，%s已预约您明天%s-%s咨询，地点：紫荆C楼409室。电话：62782007。"
	SMS_FEEDBACK_STUDENT = "温馨提示：%s你好，感谢使用我们的一对一咨询服务，请再次登录预约界面，为咨询师反馈评分，帮助我们成长。"
)

var (
	SMS_ERROR_MSG = map[string]string{
		"-1":  "没有该用户账户",
		"-2":  "接口密钥不正确",
		"-21": "MD5接口密钥加密不正确",
		"-3":  "短信数量不足",
		"-11": "该用户被禁用",
		"-14": "短信内容出现非法字符",
		"-4":  "手机号格式不正确",
		"-41": "手机号码为空",
		"-42": "短信内容为空",
		"-51": "短信签名格式不正确",
		"-6":  "IP限制",
	}
)

func SendSuccessSMS(reservation *models.Reservation) error {
	startTime := reservation.StartTime
	endTime := reservation.EndTime
	student, err := models.GetStudentById(reservation.StudentId)
	if err != nil {
		return errors.New("学生未注册")
	}
	studentSMS := fmt.Sprintf(SMS_SUCCESS_STUDENT, student.Fullname, utils.Weekdays[startTime.Weekday()],
		startTime.Month(), startTime.Day(), startTime.Format("15:04"), endTime.Format("15:04"))
	if err := sendSMS(student.Mobile, studentSMS); err != nil {
		return err
	}
	//teacher, err := models.GetTeacherById(reservation.TeacherId)
	//if err != nil {
	//	return errors.New("咨询师未注册")
	//}
	//teacherSMS := fmt.Sprintf(SMS_SUCCESS_TEACHER, teacher.Fullname, student.Fullname,
	//	utils.Weekdays[startTime.Weekday()], startTime.Month(), startTime.Day(),
	//	startTime.Format("15:04"), endTime.Format("15:04"))
	//if err := sendSMS(teacher.Mobile, teacherSMS); err != nil {
	//	return err
	//}
	return nil
}

func SendCancelSMS(reservation *models.Reservation) error {
	startTime := reservation.StartTime
	endTime := reservation.EndTime
	student, err := models.GetStudentById(reservation.StudentId)
	if err != nil {
		return errors.New("学生未注册")
	}
	studentSMS := fmt.Sprintf(SMS_CANCEL_STUDENT, student.Fullname, utils.Weekdays[startTime.Weekday()],
		startTime.Month(), startTime.Day(), startTime.Format("15:04"), endTime.Format("15:04"))
	if err := sendSMS(student.Mobile, studentSMS); err != nil {
		return err
	}
	//teacher, err := models.GetTeacherById(reservation.TeacherId)
	//if err != nil {
	//	return errors.New("咨询师未注册")
	//}
	//teacherSMS := fmt.Sprintf(SMS_CANCEL_TEACHER, teacher.Fullname, student.Fullname,
	//	utils.Weekdays[startTime.Weekday()], startTime.Month(), startTime.Day(),
	//	startTime.Format("15:04"), endTime.Format("15:04"))
	//if err := sendSMS(teacher.Mobile, teacherSMS); err != nil {
	//	return err
	//}
	return nil
}

func SendReminderSMS(reservation *models.Reservation) error {
	startTime := reservation.StartTime
	endTime := reservation.EndTime
	student, err := models.GetStudentById(reservation.StudentId)
	if err != nil {
		return errors.New("学生未注册")
	}
	studentSMS := fmt.Sprintf(SMS_REMINDER_STUDENT, student.Fullname, startTime.Format("15:04"), endTime.Format("15:04"))
	if err := sendSMS(student.Mobile, studentSMS); err != nil {
		return err
	}
	//teacher, err := models.GetTeacherById(reservation.TeacherId)
	//if err != nil {
	//	return errors.New("咨询师未注册")
	//}
	//teacherSMS := fmt.Sprintf(SMS_REMINDER_TEACHER, teacher.Fullname, student.Fullname,
	//	startTime.Format("15:04"), endTime.Format("15:04"))
	//if err := sendSMS(teacher.Mobile, teacherSMS); err != nil {
	//	return err
	//}
	return nil
}

func SendFeedbackSMS(reservation *models.Reservation) error {
	student, err := models.GetStudentById(reservation.StudentId)
	if err != nil {
		return errors.New("学生未注册")
	}
	studentSMS := fmt.Sprintf(SMS_FEEDBACK_STUDENT, student.Fullname)
	if err := sendSMS(student.Mobile, studentSMS); err != nil {
		return err
	}
	return nil
}

func sendSMS(mobile string, content string) error {
	if m := utils.IsMobile(mobile); !m {
		return errors.New("手机号格式不正确")
	}
	if config.Instance().IsSmockServer() {
		log.Printf("SMOCK Send SMS: \"%s\" to %s", content, mobile)
		return nil
	}
	requestUrl := "http://utf8.sms.webchinese.cn"
	payload := url.Values{
		"Uid":     {config.Instance().SMSUid},
		"Key":     {config.Instance().SMSKey},
		"smsMob":  {mobile},
		"smsText": {content},
	}
	requestBody := bytes.NewBufferString(payload.Encode())
	response, err := http.Post(requestUrl, "application/x-www-form-urlencoded;charset=utf8", requestBody)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	errCode := string(responseBody)
	if errMsg, ok := SMS_ERROR_MSG[errCode]; ok {
		log.Printf("Fail to send SMS \"%s\" to %s: %s", content, mobile, errMsg)
		EmailWarn("thxlfzzx报警：短信发送失败", fmt.Sprintf("Fail to send SMS \"%s\" to %s: %s", content, mobile, errMsg))
		return errors.New(fmt.Sprintf("短信发送失败：%s", errMsg))
	}
	log.Printf("Send SMS \"%s\" to %s: return %s", content, mobile, errCode)
	return nil
}

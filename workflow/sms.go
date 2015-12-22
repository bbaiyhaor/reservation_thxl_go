package workflow

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	SMS_SUCCESS_STUDENT  = "%s你好，你已成功预约星期%s（%d月%d日）%s-%s咨询，地点：紫荆C楼409室。电话：62782007。"
	SMS_SUCCESS_TEACHER  = "%s您好，%s已预约您星期%s（%d月%d日）%s-%s咨询，地点：紫荆C楼409室。电话：62782007。"
	SMS_CANCEL_STUDENT = ""
	SMS_CANCEL_TEACHER = ""
	SMS_REMINDER_STUDENT = "温馨提示：%s你好，你已成功预约明天%s-%s咨询，地点：紫荆C楼409室。电话：62782007。"
	SMS_REMINDER_TEACHER = "温馨提示：%s您好，%s已预约您明天%s-%s咨询，地点：紫荆C楼409室。电话：62782007。"
	SMS_FEEDBACK_STUDENT = "温馨提示：%s你好，感谢使用我们的一对一咨询服务，请再次登录预约界面，为咨询师反馈评分，帮助我们成长。"
)

func SendSuccessSMS(reservation *models.Reservation) error {
	startTime := reservation.StartTime.In(utils.Location)
	endTime := reservation.EndTime.In(utils.Location)
	student, err := models.GetStudentById(reservation.StudentId)
	if err != nil {
		return errors.New("学生未注册")
	}
	studentSMS := fmt.Sprintf(SMS_SUCCESS_STUDENT, student.Fullname, utils.Weekdays[startTime.Weekday()],
		startTime.Month(), startTime.Day(), startTime.Format("15:04"),
		endTime.Format("15:04"))
	if err := sendSMS(student.Mobile, studentSMS); err != nil {
		return err
	}
	teacher, err := models.GetTeacherById(reservation.TeacherId)
	if err != nil {
		return errors.New("咨询师未注册")
	}
	teacherSMS := fmt.Sprintf(SMS_SUCCESS_TEACHER, teacher.Fullname, student.Fullname,
		utils.Weekdays[startTime.Weekday()], startTime.Month(), startTime.Day(),
		startTime.Format("15:04"), endTime.Format("15:04"))
	if err := sendSMS(teacher.Mobile, teacherSMS); err != nil {
		return err
	}
	return nil
}

func SendCancelSMS(reservation *models.Reservation) error {
//	startTime := reservation.StartTime.In(utils.Location)
//	endTime := reservation.EndTime.In(utils.Location)
	student, err := models.GetStudentById(reservation.StudentId)
	if err != nil {
		return errors.New("学生未注册")
	}
	studentSMS := fmt.Sprintf(SMS_CANCEL_STUDENT)
	if err := sendSMS(student.Mobile, studentSMS); err != nil {
		return err
	}
	teacher, err := models.GetTeacherById(reservation.TeacherId)
	if err != nil {
		return errors.New("咨询师未注册")
	}
	teacherSMS := fmt.Sprintf(SMS_CANCEL_TEACHER)
	if err := sendSMS(teacher.Mobile, teacherSMS); err != nil {
		return err
	}
	return nil
}

func SendReminderSMS(reservation *models.Reservation) error {
	startTime := reservation.StartTime.In(utils.Location)
	endTime := reservation.EndTime.In(utils.Location)
	student, err := models.GetStudentById(reservation.StudentId)
	if err != nil {
		return errors.New("学生未注册")
	}
	studentSMS := fmt.Sprintf(SMS_REMINDER_STUDENT, student.Fullname, startTime.Format("15:04"), endTime.Format("15:04"))
	if err := sendSMS(student.Mobile, studentSMS); err != nil {
		return err
	}
	teacher, err := models.GetTeacherById(reservation.TeacherId)
	if err != nil {
		return errors.New("咨询师未注册")
	}
	teacherSMS := fmt.Sprintf(SMS_REMINDER_TEACHER, teacher.Fullname, student.Fullname,
		startTime.Format("15:04"), endTime.Format("15:04"))
	if err := sendSMS(teacher.Mobile, teacherSMS); err != nil {
		return err
	}
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
	appEnv := os.Getenv("RESERVATION_THXL_ENV")
	if !strings.EqualFold(appEnv, "ONLINE") {
		fmt.Printf("Send SMS: \"%s\" to %s.\n", content, mobile)
		return nil
	}
	requestUrl := "http://utf8.sms.webchinese.cn"
	payload := url.Values{
		"Uid":     {utils.SMS_UID},
		"Key":     {utils.SMS_KEY},
		"smsMob":  {mobile},
		"smsText": {content},
	}
	requestBody := bytes.NewBufferString(payload.Encode())
	response, err := http.Post(requestUrl, "application/x-www-form-urlencoded;charset=utf8", requestBody)
	if err != nil {
		return errors.New("短信发送失败")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("短信发送失败")
	}
	fmt.Println(string(responseBody))
	if code, err := strconv.Atoi(string(responseBody)); err != nil || code < 0 {
		return errors.New(fmt.Sprintf("短信发送失败,ErrCode:%d", code))
	}
	return nil
}

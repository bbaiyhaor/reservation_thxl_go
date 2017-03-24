package buslogic

import (
	"github.com/shudiwsh2009/reservation_thxl_go/config"
	"github.com/shudiwsh2009/reservation_thxl_go/model"
	re "github.com/shudiwsh2009/reservation_thxl_go/rerror"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/mijia/sweb/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	SMS_SUCCESS_STUDENT        = "%s你好，你已成功预约星期%s（%d月%d日）%s-%s咨询，地点：紫荆C楼409室。电话：62782007。请记得携带学生卡，如需取消预约请提前一天致电62782007。"
	SMS_SUCCESS_TEACHER        = "%s您好，%s已预约您星期%s（%d月%d日）%s-%s咨询，地点：紫荆C楼409室。电话：62782007。"
	SMS_CANCEL_STUDENT         = "%s你好，因特殊原因你预约的星期%s（%d月%d日）%s-%s咨询已被取消，详询62782007。"
	SMS_CANCEL_TEACHER         = "%s您好，%s预约的星期%s（%d月%d日）%s-%s咨询已被取消，详询62782007。"
	SMS_REMINDER_STUDENT       = "温馨提示：%s你好，你已成功预约明天%s-%s咨询，地点：紫荆C楼409室。电话：62782007。"
	SMS_REMINDER_TEACHER       = "温馨提示：%s您好，%s已预约您明天%s-%s咨询，地点：紫荆C楼409室。电话：62782007。"
	SMS_FEEDBACK_STUDENT       = "温馨提示：%s你好，感谢使用我们的一对一咨询服务，请再次登录预约界面，为咨询师反馈评分，帮助我们成长。"
	SMS_TEACHER_RESET_PASSWORD = "%s您好，您正在申请重置咨询中心登录密码，验证码：%s，请在10分钟内输入。"
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

func (w *Workflow) SendSuccessSMS(reservation *model.Reservation) error {
	startTime := reservation.StartTime
	endTime := reservation.EndTime
	student, err := w.mongoClient.GetStudentById(reservation.StudentId)
	if err != nil {
		return re.NewRErrorCode("学生未注册", nil, re.ERROR_DATABASE)
	}
	studentSMS := fmt.Sprintf(SMS_SUCCESS_STUDENT, student.Fullname, utils.Weekdays[startTime.Weekday()],
		startTime.Month(), startTime.Day(), startTime.Format("15:04"), endTime.Format("15:04"))
	if err := w.sendSMS(student.Mobile, studentSMS); err != nil {
		return err
	}
	/*
		teacher, err := models.GetTeacherById(reservation.TeacherId)
		if err != nil {
			return errors.New("咨询师未注册")
		}
		teacherSMS := fmt.Sprintf(SMS_SUCCESS_TEACHER, teacher.Fullname, student.Fullname,
			utils.Weekdays[startTime.Weekday()], startTime.Month(), startTime.Day(),
			startTime.Format("15:04"), endTime.Format("15:04"))
		if err := w.sendSMS(teacher.Mobile, teacherSMS); err != nil {
			return err
		}
	*/
	return nil
}

func (w *Workflow) SendCancelSMS(reservation *model.Reservation) error {
	startTime := reservation.StartTime
	endTime := reservation.EndTime
	student, err := w.mongoClient.GetStudentById(reservation.StudentId)
	if err != nil {
		return re.NewRErrorCode("学生未注册", nil, re.ERROR_DATABASE)
	}
	studentSMS := fmt.Sprintf(SMS_CANCEL_STUDENT, student.Fullname, utils.Weekdays[startTime.Weekday()],
		startTime.Month(), startTime.Day(), startTime.Format("15:04"), endTime.Format("15:04"))
	if err := w.sendSMS(student.Mobile, studentSMS); err != nil {
		return err
	}
	/*
		teacher, err := models.GetTeacherById(reservation.TeacherId)
		if err != nil {
			return errors.New("咨询师未注册")
		}
		teacherSMS := fmt.Sprintf(SMS_CANCEL_TEACHER, teacher.Fullname, student.Fullname,
			utils.Weekdays[startTime.Weekday()], startTime.Month(), startTime.Day(),
			startTime.Format("15:04"), endTime.Format("15:04"))
		if err := w.sendSMS(teacher.Mobile, teacherSMS); err != nil {
			return err
		}
	*/
	return nil
}

func (w *Workflow) SendReminderSMS(reservation *model.Reservation) error {
	startTime := reservation.StartTime
	endTime := reservation.EndTime
	student, err := w.mongoClient.GetStudentById(reservation.StudentId)
	if err != nil {
		return re.NewRErrorCode("学生未注册", nil, re.ERROR_DATABASE)
	}
	studentSMS := fmt.Sprintf(SMS_REMINDER_STUDENT, student.Fullname, startTime.Format("15:04"), endTime.Format("15:04"))
	if err := w.sendSMS(student.Mobile, studentSMS); err != nil {
		return err
	}
	/*
		teacher, err := models.GetTeacherById(reservation.TeacherId)
		if err != nil {
			return errors.New("咨询师未注册")
		}
		teacherSMS := fmt.Sprintf(SMS_REMINDER_TEACHER, teacher.Fullname, student.Fullname,
			startTime.Format("15:04"), endTime.Format("15:04"))
		if err := w.sendSMS(teacher.Mobile, teacherSMS); err != nil {
			return err
		}
	*/
	return nil
}

func (w *Workflow) SendFeedbackSMS(reservation *model.Reservation) error {
	student, err := w.mongoClient.GetStudentById(reservation.StudentId)
	if err != nil {
		return re.NewRErrorCode("学生未注册", nil, re.ERROR_DATABASE)
	}
	studentSMS := fmt.Sprintf(SMS_FEEDBACK_STUDENT, student.Fullname)
	if err := w.sendSMS(student.Mobile, studentSMS); err != nil {
		return err
	}
	return nil
}

func (w *Workflow) SendTeacherResetPasswordSMS(teacher *model.Teacher, verifyCode string) error {
	verifySMS := fmt.Sprintf(SMS_TEACHER_RESET_PASSWORD, teacher.Fullname, verifyCode)
	return w.sendSMS(teacher.Mobile, verifySMS)
}

func (w *Workflow) sendSMS(mobile string, content string) error {
	if !utils.IsMobile(mobile) {
		return re.NewRErrorCode("手机号格式不正确", nil, re.ERROR_FORMAT_MOBILE)
	}
	if config.Instance().IsSmockServer() {
		log.Infof("SMOCK Send SMS: \"%s\" to %s", content, mobile)
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
		log.Errorf("Fail to send SMS \"%s\" to %s: %s", content, mobile, errMsg)
		EmailWarn("thxlfzzx报警：短信发送失败", fmt.Sprintf("Fail to send SMS \"%s\" to %s: %s", content, mobile, errMsg))
		return re.NewRError(fmt.Sprintf("短信发送失败：%s", errMsg), nil)
	}
	log.Infof("Send SMS \"%s\" to %s: return %s", content, mobile, errCode)
	return nil
}

// 生成验证码
var verifyCodeTable = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateVerifyCode(length int) (string, error) {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if err != nil || n != length {
		return "", re.NewRError("生成验证码出错", err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = verifyCodeTable[int(b[i])%len(verifyCodeTable)]
	}
	return string(b), nil
}

// 每天20:00发送第二天预约咨询的提醒短信
func (w *Workflow) SendTomorrowReservationReminderSMS() {
	today := utils.BeginOfDay(time.Now())
	from := today.AddDate(0, 0, 1)
	to := today.AddDate(0, 0, 2)
	reservations, err := w.mongoClient.GetReservationsBetweenTime(from, to)
	if err != nil {
		log.Errorf("获取咨询列表失败：%v", err)
		return
	}
	succCnt, failCnt := 0, 0
	for _, reservation := range reservations {
		if reservation.Status == model.RESERVATION_STATUS_RESERVATED {
			if err = w.SendReminderSMS(reservation); err == nil {
				succCnt++
			} else {
				log.Errorf("发送短信失败：%+v %+v", reservation, err)
				failCnt++
			}
		}
	}
	log.Infof("发送%d个预约记录的提醒短信，成功%d个，失败%d个", succCnt+failCnt, succCnt, failCnt)
}

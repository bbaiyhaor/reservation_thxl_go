package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	re "bitbucket.org/shudiwsh2009/reservation_thxl_go/rerror"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"fmt"
	"gopkg.in/redis.v5"
	"time"
)

const (
	ADMIN_DEFAULT_PASSWORD   = "THXLFZZX"
	TEACHER_DEFAULT_PASSWORD = "thxlfzzx"
)

// 学生登录
func (w *Workflow) StudentLogin(username string, password string) (*model.Student, error) {
	if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if password == "" {
		return nil, re.NewRErrorCodeContext("password is empty", nil, re.ERROR_MISSING_PARAM, "password")
	}
	student, err := w.mongoClient.GetStudentByUsername(username)
	if err == nil && student.Password == model.EncodePassword(student.Salt, password) {
		return student, nil
	}
	if err == nil {
		if student.Salt == "" && utils.ValidatePassword(password, student.EncryptedPassword) {
			student.Password = password
			student.PreInsert()
			w.mongoClient.UpdateStudent(student)
			return student, nil
		} else if student.Password == model.EncodePassword(student.Salt, password) {
			return student, nil
		}
	}
	return nil, re.NewRErrorCode("wrong password", nil, re.ERROR_LOGIN_PASSWORD_WRONG)
}

// 咨询师登录
func (w *Workflow) TeacherLogin(username string, password string) (*model.Teacher, error) {
	if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if password == "" {
		return nil, re.NewRErrorCodeContext("password is empty", nil, re.ERROR_MISSING_PARAM, "password")
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(username)
	if err == nil {
		if teacher.Salt == "" && utils.ValidatePassword(password, teacher.EncryptedPassword) {
			teacher.Password = password
			teacher.PreInsert()
			w.mongoClient.UpdateTeacher(teacher)
			return teacher, nil
		} else if teacher.Password == model.EncodePassword(teacher.Salt, password) {
			return teacher, nil
		}
	}
	return nil, re.NewRErrorCode("wrong password", nil, re.ERROR_LOGIN_PASSWORD_WRONG)
}

// 咨询师更改密码
func (w *Workflow) TeacherChangePassword(username, oldPassword, newPassword string, userId string, userType int) (*model.Teacher, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_TEACHER {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ERROR_NOT_AUTHORIZED)
	} else if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if oldPassword == "" {
		return nil, re.NewRErrorCodeContext("old password is empty", nil, re.ERROR_MISSING_PARAM, "old_password")
	} else if newPassword == "" {
		return nil, re.NewRErrorCodeContext("new password is empty", nil, re.ERROR_MISSING_PARAM, "new_password")
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	} else if username != teacher.Username {
		return nil, re.NewRErrorCode("username not match", nil, re.ERROR_LOGIN_PWDCHANGE_USERNAME_MISMATCH)
	}
	if teacher.Password != model.EncodePassword(teacher.Salt, oldPassword) {
		return nil, re.NewRErrorCode("old password mismatch", nil, re.ERROR_LOGIN_PWDCHANGE_OLDPWD_MISMATCH)
	}
	if oldPassword == newPassword {
		return nil, re.NewRErrorCode("new password should be different", nil, re.ERROR_LOGIN_PWDCHANGE_OLDPWD_EQUAL_NEWPED)
	}
	teacher.Password = newPassword
	teacher.PreInsert()
	if err = w.mongoClient.UpdateTeacher(teacher); err != nil {
		return nil, re.NewRErrorCode("fail to update teacher", err, re.ERROR_DATABASE)
	}
	return teacher, nil
}

// 咨询师重置密码发送短信
func (w *Workflow) TeacherResetPasswordSms(username, mobile string) error {
	if username == "" {
		return re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if mobile == "" {
		return re.NewRErrorCodeContext("mobile is empty", nil, re.ERROR_MISSING_PARAM, "mobile")
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(username)
	if err != nil || teacher.UserType != model.USER_TYPE_TEACHER {
		return re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	} else if mobile != teacher.Mobile {
		return re.NewRErrorCode("mobile not match", nil, re.ERROR_LOGIN_PWDCHANGE_MOBILE_MISMATCh)
	}
	verifyCode, err := GenerateVerifyCode(6)
	if err != nil {
		return re.NewRErrorCode("fail to generate verify code", err, re.ERROR_SEND_SMS)
	}
	err = w.redisClient.Set(fmt.Sprintf(model.REDIS_KEY_TEACHER_RESET_PASSWORD_VERIFY_CODE, teacher.Id.Hex()), verifyCode, 10*time.Minute).Err()
	if err != nil {
		return re.NewRErrorCode("fail to set verify code to redis", err, re.ERROR_SEND_SMS)
	}
	if err = w.SendResetPasswordSMS(teacher, verifyCode); err != nil {
		return re.NewRErrorCode("fail to send sms", err, re.ERROR_SEND_SMS)
	}
	return nil
}

func (w *Workflow) TeacherRestPasswordVerify(username, newPassword, verifyCode string) error {
	if username == "" {
		return re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if newPassword == "" {
		return re.NewRErrorCodeContext("new password is empty", nil, re.ERROR_MISSING_PARAM, "new_password")
	} else if verifyCode == "" {
		return re.NewRErrorCodeContext("verify code is empty", nil, re.ERROR_MISSING_PARAM, "verify_code")
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(username)
	if err != nil || teacher.UserType != model.USER_TYPE_TEACHER {
		return re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	}
	val, err := w.redisClient.Get(fmt.Sprintf(model.REDIS_KEY_TEACHER_RESET_PASSWORD_VERIFY_CODE, teacher.Id.Hex())).Result()
	if err != nil || err == redis.Nil || val != verifyCode {
		return re.NewRErrorCode("fail to get verify code from redis", err, re.ERROR_LOGIN_PWDCHANGE_VERIFY_CODE_WRONG)
	}
	w.redisClient.Del(fmt.Sprintf(model.REDIS_KEY_TEACHER_RESET_PASSWORD_VERIFY_CODE, teacher.Id.Hex()))
	teacher.Password = newPassword
	teacher.PreInsert()
	if err = w.mongoClient.UpdateTeacher(teacher); err != nil {
		return re.NewRErrorCode("fail to update teacher", err, re.ERROR_DATABASE)
	}
	return nil
}

// 管理员登录
func (w *Workflow) AdminLogin(username string, password string) (*model.Admin, error) {
	if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if password == "" {
		return nil, re.NewRErrorCodeContext("password is empty", nil, re.ERROR_MISSING_PARAM, "password")
	}
	admin, err := w.mongoClient.GetAdminByUsername(username)
	if err == nil {
		if admin.Salt == "" && utils.ValidatePassword(password, admin.EncryptedPassword) {
			admin.Password = password
			admin.PreInsert()
			w.mongoClient.UpdateAdmin(admin)
			return admin, nil
		} else if admin.Password == model.EncodePassword(admin.Salt, password) {
			return admin, nil
		}
	}
	return nil, re.NewRErrorCode("wrong password", nil, re.ERROR_LOGIN_PASSWORD_WRONG)
}

// 学生注册
func (w *Workflow) StudentRegister(username string, password string) (*model.Student, error) {
	if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if password == "" {
		return nil, re.NewRErrorCodeContext("password is empty", nil, re.ERROR_MISSING_PARAM, "password")
	}
	if !utils.IsStudentId(username) {
		return nil, re.NewRErrorCode("student id format is wrong", nil, re.ERROR_FORMAT_STUDENTID)
	}
	if student, _ := w.mongoClient.GetStudentByUsername(username); student != nil {
		return nil, re.NewRErrorCode("student already exists", nil, re.ERROR_EXIST_USERNAME)
	}
	student := &model.Student{
		Username: username,
		Password: password,
		UserType: model.USER_TYPE_STUDENT,
	}
	archive, err := w.mongoClient.GetArchiveByStudentUsername(username)
	if err == nil {
		student.ArchiveCategory = archive.ArchiveCategory
		student.ArchiveNumber = archive.ArchiveNumber
	}
	if err := w.mongoClient.InsertStudent(student); err != nil {
		return nil, re.NewRErrorCode("fail to insert student", err, re.ERROR_DATABASE)
	}
	return student, nil
}

// 更新session
func (w *Workflow) UpdateSession(userId string, userType int) (map[string]interface{}, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("user not login", nil, re.ERROR_NO_LOGIN)
	}
	result := make(map[string]interface{})
	switch userType {
	case model.USER_TYPE_ADMIN:
		admin, err := w.mongoClient.GetAdminById(userId)
		if err != nil || admin.UserType != userType {
			return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
		}
		result["user_id"] = admin.Id.Hex()
		result["username"] = admin.Username
		result["user_type"] = admin.UserType
		result["fullname"] = admin.Fullname
	case model.USER_TYPE_TEACHER:
		teacher, err := w.mongoClient.GetTeacherById(userId)
		if err != nil || teacher.UserType != userType {
			return nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
		}
		result["user_id"] = teacher.Id.Hex()
		result["username"] = teacher.Username
		result["user_type"] = teacher.UserType
		result["fullname"] = teacher.Fullname
	case model.USER_TYPE_STUDENT:
		student, err := w.mongoClient.GetStudentById(userId)
		if err != nil || student.UserType != userType {
			return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
		}
		result["user_id"] = student.Id.Hex()
		result["username"] = student.Username
		result["user_type"] = student.UserType
		result["fullname"] = student.Fullname
	default:
		return nil, re.NewRErrorCode("fail to get user", nil, re.ERROR_NO_USER)
	}
	return result, nil
}

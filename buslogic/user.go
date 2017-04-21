package buslogic

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxl_go/model"
	re "github.com/shudiwsh2009/reservation_thxl_go/rerror"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"gopkg.in/redis.v5"
	"time"
)

const (
	ADMIN_DEFAULT_PASSWORD   = "THXLFZZX"
	TEACHER_DEFAULT_PASSWORD = "thxlfzzx"
)

var (
	USER_LOGIN_COUNT = map[int]int{
		model.USER_TYPE_ADMIN:   1,
		model.USER_TYPE_TEACHER: 1,
		model.USER_TYPE_STUDENT: 5,
	}
	USER_LOGIN_EXPIRE = map[int]time.Duration{
		model.USER_TYPE_ADMIN:   time.Hour * 24 * 7,
		model.USER_TYPE_TEACHER: time.Hour * 24 * 7,
		model.USER_TYPE_STUDENT: time.Hour * 24 * 30,
	}
)

// 学生登录
func (w *Workflow) StudentLogin(username string, password string) (*model.Student, error) {
	if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if password == "" {
		return nil, re.NewRErrorCodeContext("password is empty", nil, re.ERROR_MISSING_PARAM, "password")
	}
	student, err := w.mongoClient.GetStudentByUsername(username)
	if err == nil && student != nil && student.UserType == model.USER_TYPE_STUDENT {
		if student.Password == model.EncodePassword(student.Salt, password) {
			return student, nil
		} else if student.Salt == "" && utils.ValidatePassword(password, student.EncryptedPassword) {
			student.Password = password
			student.PreInsert()
			w.mongoClient.UpdateStudent(student)
			return student, nil
		} else if student.Salt != "" && student.Password == model.EncodePassword(student.Salt, password) {
			return student, nil
		}
	}
	return nil, re.NewRErrorCode("wrong password", nil, re.ERROR_LOGIN_PASSWORD_WRONG)
}

// 学生注册
func (w *Workflow) StudentRegister(username string, password string, fullname string, gender string, birthday string,
	school string, grade string, currentAddress string, familyAddress string, mobile string, email string,
	experienceTime string, experienceLocation string, experienceTeacher string, fatherAge string, fatherJob string, fatherEdu string,
	motherAge string, motherJob string, motherEdu string, parentMarriage string) (*model.Student, error) {
	if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if password == "" {
		return nil, re.NewRErrorCodeContext("password is empty", nil, re.ERROR_MISSING_PARAM, "password")
	} else if fullname == "" {
		return nil, re.NewRErrorCodeContext("fullname is empty", nil, re.ERROR_MISSING_PARAM, "fullname")
	} else if gender == "" {
		return nil, re.NewRErrorCodeContext("gender is empty", nil, re.ERROR_MISSING_PARAM, "gender")
	} else if birthday == "" {
		return nil, re.NewRErrorCodeContext("birthday is empty", nil, re.ERROR_MISSING_PARAM, "birthday")
	} else if school == "" {
		return nil, re.NewRErrorCodeContext("school is empty", nil, re.ERROR_MISSING_PARAM, "school")
	} else if grade == "" {
		return nil, re.NewRErrorCodeContext("grade is empty", nil, re.ERROR_MISSING_PARAM, "grade")
	} else if currentAddress == "" {
		return nil, re.NewRErrorCodeContext("current_address is empty", nil, re.ERROR_MISSING_PARAM, "current_address")
	} else if familyAddress == "" {
		return nil, re.NewRErrorCodeContext("family_address is empty", nil, re.ERROR_MISSING_PARAM, "family_address")
	} else if mobile == "" {
		return nil, re.NewRErrorCodeContext("mobile is empty", nil, re.ERROR_MISSING_PARAM, "mobile")
	} else if email == "" {
		return nil, re.NewRErrorCodeContext("email is empty", nil, re.ERROR_MISSING_PARAM, "email")
	} else if !utils.IsMobile(mobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ERROR_FORMAT_MOBILE)
	} else if !utils.IsEmail(email) {
		return nil, re.NewRErrorCode("email format is wrong", nil, re.ERROR_FORMAT_EMAIL)
	}
	if !utils.IsStudentId(username) {
		return nil, re.NewRErrorCode("student id format is wrong", nil, re.ERROR_FORMAT_STUDENTID)
	}
	student, err := w.mongoClient.GetStudentByUsername(username)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	} else if student != nil && student.UserType == model.USER_TYPE_STUDENT {
		return nil, re.NewRErrorCode("student already exists", nil, re.ERROR_EXIST_USERNAME)
	}
	student = &model.Student{
		Username:       username,
		Password:       password,
		UserType:       model.USER_TYPE_STUDENT,
		Fullname:       fullname,
		Gender:         gender,
		Birthday:       birthday,
		School:         school,
		Grade:          grade,
		CurrentAddress: currentAddress,
		FamilyAddress:  familyAddress,
		Mobile:         mobile,
		Email:          email,
		Experience: model.Experience{
			Time:     experienceTime,
			Location: experienceLocation,
			Teacher:  experienceTeacher,
		},
		ParentInfo: model.ParentInfo{
			FatherAge:      fatherAge,
			FatherJob:      fatherJob,
			FatherEdu:      fatherEdu,
			MotherAge:      motherAge,
			MotherJob:      motherJob,
			MotherEdu:      motherEdu,
			ParentMarriage: parentMarriage,
		},
	}
	archive, err := w.mongoClient.GetArchiveByStudentUsername(username)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get archive", err, re.ERROR_DATABASE)
	} else if archive != nil {
		student.ArchiveCategory = archive.ArchiveCategory
		student.ArchiveNumber = archive.ArchiveNumber
	}
	if err := w.mongoClient.InsertStudent(student); err != nil {
		return nil, re.NewRErrorCode("fail to insert student", err, re.ERROR_DATABASE)
	}
	return student, nil
}

// 咨询师登录
func (w *Workflow) TeacherLogin(username string, password string) (*model.Teacher, error) {
	if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if password == "" {
		return nil, re.NewRErrorCodeContext("password is empty", nil, re.ERROR_MISSING_PARAM, "password")
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(username)
	if err == nil && teacher != nil && teacher.UserType == model.USER_TYPE_TEACHER {
		if teacher.Salt == "" && utils.ValidatePassword(password, teacher.EncryptedPassword) {
			teacher.Password = password
			teacher.PreInsert()
			w.mongoClient.UpdateTeacher(teacher)
			return teacher, nil
		} else if teacher.Salt != "" && teacher.Password == model.EncodePassword(teacher.Salt, password) {
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
	if err != nil || teacher == nil || teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	} else if username != teacher.Username {
		return nil, re.NewRErrorCode("username not match", nil, re.ERROR_LOGIN_PWDCHANGE_INFO_MISMATCH)
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
func (w *Workflow) TeacherResetPasswordSms(username, fullname, mobile string) error {
	if username == "" {
		return re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if fullname == "" {
		return re.NewRErrorCodeContext("fullname is empty", nil, re.ERROR_MISSING_PARAM, "fullname")
	} else if mobile == "" {
		return re.NewRErrorCodeContext("mobile is empty", nil, re.ERROR_MISSING_PARAM, "mobile")
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(username)
	if err != nil || teacher == nil || teacher.UserType != model.USER_TYPE_TEACHER {
		return re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	} else if fullname != teacher.Fullname || mobile != teacher.Mobile {
		return re.NewRErrorCode("fullname or mobile not match", nil, re.ERROR_LOGIN_PWDCHANGE_INFO_MISMATCH)
	}
	verifyCode, err := GenerateVerifyCode(6)
	if err != nil {
		return re.NewRErrorCode("fail to generate verify code", err, re.ERROR_SEND_SMS)
	}
	err = w.redisClient.Set(fmt.Sprintf(model.REDIS_KEY_TEACHER_RESET_PASSWORD_VERIFY_CODE, teacher.Id.Hex()), verifyCode, 10*time.Minute).Err()
	if err != nil {
		return re.NewRErrorCode("fail to set verify code to redis", err, re.ERROR_SEND_SMS)
	}
	if err = w.SendTeacherResetPasswordSMS(teacher, verifyCode); err != nil {
		return re.NewRErrorCode("fail to send sms", err, re.ERROR_SEND_SMS)
	}
	return nil
}

// 咨询师验证短信重置密码
func (w *Workflow) TeacherRestPasswordVerify(username, newPassword, verifyCode string) error {
	if username == "" {
		return re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if newPassword == "" {
		return re.NewRErrorCodeContext("new password is empty", nil, re.ERROR_MISSING_PARAM, "new_password")
	} else if verifyCode == "" {
		return re.NewRErrorCodeContext("verify code is empty", nil, re.ERROR_MISSING_PARAM, "verify_code")
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(username)
	if err != nil || teacher == nil || teacher.UserType != model.USER_TYPE_TEACHER {
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
	w.ClearUserLoginRedisKey(teacher.Id.Hex(), teacher.UserType)
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
	if err == nil && admin != nil && admin.UserType == model.USER_TYPE_ADMIN {
		if admin.Salt == "" && utils.ValidatePassword(password, admin.EncryptedPassword) {
			admin.Password = password
			admin.PreInsert()
			w.mongoClient.UpdateAdmin(admin)
			return admin, nil
		} else if admin.Salt != "" && admin.Password == model.EncodePassword(admin.Salt, password) {
			return admin, nil
		}
	}
	return nil, re.NewRErrorCode("wrong password", nil, re.ERROR_LOGIN_PASSWORD_WRONG)
}

// 管理员更改密码
func (w *Workflow) AdminChangePassword(username, oldPassword, newPassword string, userId string, userType int) (*model.Admin, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ERROR_NOT_AUTHORIZED)
	} else if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ERROR_MISSING_PARAM, "username")
	} else if oldPassword == "" {
		return nil, re.NewRErrorCodeContext("old password is empty", nil, re.ERROR_MISSING_PARAM, "old_password")
	} else if newPassword == "" {
		return nil, re.NewRErrorCodeContext("new password is empty", nil, re.ERROR_MISSING_PARAM, "new_password")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	} else if username != admin.Username {
		return nil, re.NewRErrorCode("username not match", nil, re.ERROR_LOGIN_PWDCHANGE_INFO_MISMATCH)
	}
	if admin.Password != model.EncodePassword(admin.Salt, oldPassword) {
		return nil, re.NewRErrorCode("old password mismatch", nil, re.ERROR_LOGIN_PWDCHANGE_OLDPWD_MISMATCH)
	}
	if oldPassword == newPassword {
		return nil, re.NewRErrorCode("new password should be different", nil, re.ERROR_LOGIN_PWDCHANGE_OLDPWD_EQUAL_NEWPED)
	}
	admin.Password = newPassword
	admin.PreInsert()
	if err = w.mongoClient.UpdateAdmin(admin); err != nil {
		return nil, re.NewRErrorCode("fail to update admin", err, re.ERROR_DATABASE)
	}
	w.ClearUserLoginRedisKey(admin.Id.Hex(), admin.UserType)
	return admin, nil
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
		if err != nil || admin == nil || admin.UserType != userType {
			return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
		}
		result["user"] = w.WrapAdmin(admin)
	case model.USER_TYPE_TEACHER:
		teacher, err := w.mongoClient.GetTeacherById(userId)
		if err != nil || teacher == nil || teacher.UserType != userType {
			return nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
		}
		result["user"] = w.WrapTeacher(teacher)
	case model.USER_TYPE_STUDENT:
		student, err := w.mongoClient.GetStudentById(userId)
		if err != nil || student == nil || student.UserType != userType {
			return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
		}
		result["user"] = w.WrapSimpleStudent(student)
	default:
		return nil, re.NewRErrorCode("fail to get user", nil, re.ERROR_NO_USER)
	}
	return result, nil
}

// external: 重置账户密码
func (w *Workflow) ResetUserPassword(username string, userType int, password string) error {
	if username == "" || password == "" {
		return re.NewRError("missing parameters", nil)
	}
	var err error
	var userId string
	switch userType {
	case model.USER_TYPE_STUDENT:
		student, err := w.mongoClient.GetStudentByUsername(username)
		if err != nil || student == nil || student.UserType != userType {
			return re.NewRError("fail to get student", err)
		}
		student.Password = password
		student.PreInsert()
		err = w.mongoClient.UpdateStudent(student)
		userId = student.Id.Hex()
	case model.USER_TYPE_TEACHER:
		teacher, err := w.mongoClient.GetTeacherByUsername(username)
		if err != nil || teacher == nil || teacher.UserType != userType {
			return re.NewRError("fail to get teacher", err)
		}
		teacher.Password = password
		teacher.PreInsert()
		err = w.mongoClient.UpdateTeacher(teacher)
		userId = teacher.Id.Hex()
	case model.USER_TYPE_ADMIN:
		admin, err := w.mongoClient.GetAdminByUsername(username)
		if err != nil || admin == nil || admin.UserType != userType {
			return re.NewRError("fail to get student", err)
		}
		admin.Password = password
		admin.PreInsert()
		err = w.mongoClient.UpdateAdmin(admin)
		userId = admin.Id.Hex()
	default:
		return re.NewRError(fmt.Sprintf("unknown user_type: %d", userType), nil)
	}
	if err != nil {
		return re.NewRError("fail to update user", err)
	}
	return w.ClearUserLoginRedisKey(userId, userType)
}

func (w *Workflow) ClearUserLoginRedisKey(userId string, userType int) error {
	redisKeys, err := w.redisClient.Keys(fmt.Sprintf(model.REDIS_KEY_USER_LOGIN, userType, userId, "*")).Result()
	if err != nil {
		return re.NewRError("fail to get user login session keys from redis", err)
	}
	for _, k := range redisKeys {
		if err := w.redisClient.Del(k).Err(); err != nil {
			return err
		}
	}
	return nil
}

// external: 添加新管理员
func (w *Workflow) AddNewAdmin(username string, password string) (*model.Admin, error) {
	if username == "" || password == "" {
		return nil, re.NewRError("missing parameters", nil)
	}
	oldAdmin, err := w.mongoClient.GetAdminByUsername(username)
	if err != nil {
		return nil, re.NewRError("fail to get database", err)
	} else if oldAdmin != nil && oldAdmin.UserType == model.USER_TYPE_ADMIN {
		return oldAdmin, re.NewRError(fmt.Sprintf("admin already exists: %+v", oldAdmin), nil)
	}
	newAdmin := &model.Admin{
		Username: username,
		Password: password,
		UserType: model.USER_TYPE_ADMIN,
	}
	err = w.mongoClient.InsertAdmin(newAdmin)
	return newAdmin, err
}

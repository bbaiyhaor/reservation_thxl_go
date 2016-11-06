package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"errors"
	"fmt"
	"gopkg.in/redis.v3"
	"time"
)

const (
	AdminDefaultPassword   = "THXLFZZX"
	TeacherDefaultPassword = "thxlfzzx"
)

// 学生登录
func (w *Workflow) StudentLogin(username string, password string) (*model.Student, error) {
	if username == "" {
		return nil, errors.New("用户名为空")
	} else if password == "" {
		return nil, errors.New("密码为空")
	}
	student, err := w.model.GetStudentByUsername(username)
	if err == nil {
		if student.EncryptedPassword != "" {
			if utils.ValidatePassword(password, student.EncryptedPassword) {
				return student, nil
			}
		} else if password == student.Password {
			if encryptedPassword, err := utils.EncryptPassword(password); err == nil {
				student.EncryptedPassword = encryptedPassword
				w.model.UpsertStudent(student)
			}
			return student, nil
		}
	}
	return nil, errors.New("用户名或密码不正确")
}

// 咨询师登录
func (w *Workflow) TeacherLogin(username string, password string) (*model.Teacher, error) {
	if username == "" {
		return nil, errors.New("用户名为空")
	} else if password == "" {
		return nil, errors.New("密码为空")
	}
	teacher, err := w.model.GetTeacherByUsername(username)
	if err == nil {
		if teacher.EncryptedPassword != "" {
			if utils.ValidatePassword(password, teacher.EncryptedPassword) {
				return teacher, nil
			}
		} else if password == teacher.Password {
			if encryptedPassword, err := utils.EncryptPassword(password); err == nil {
				teacher.EncryptedPassword = encryptedPassword
				w.model.UpsertTeacher(teacher)
			}
			return teacher, nil
		}
	}
	return nil, errors.New("用户名或密码不正确")
}

// 咨询师更改密码
func (w *Workflow) TeacherChangePassword(username, oldPassword, newPassword string, userId string, userType int) (*model.Teacher, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	} else if userType != model.USER_TYPE_TEACHER {
		return nil, errors.New("权限不足")
	} else if username == "" {
		return nil, errors.New("用户名为空")
	} else if oldPassword == "" {
		return nil, errors.New("旧密码为空")
	} else if newPassword == "" {
		return nil, errors.New("新密码为空")
	}
	teacher, err := w.model.GetTeacherById(userId)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, errors.New("权限不足")
	} else if username != teacher.Username {
		return nil, errors.New("权限不足")
	}
	if (teacher.EncryptedPassword != "" && !utils.ValidatePassword(oldPassword, teacher.EncryptedPassword)) ||
		(teacher.EncryptedPassword == "" && oldPassword != teacher.Password) {
		return nil, errors.New("旧密码不正确")
	}
	if oldPassword == newPassword {
		return nil, errors.New("新密码不能与原有密码一样")
	}
	encryptedPassword, err := utils.EncryptPassword(newPassword)
	if err != nil {
		return nil, errors.New("新密码不符合要求")
	}
	teacher.EncryptedPassword = encryptedPassword
	if err = w.model.UpsertTeacher(teacher); err != nil {
		return nil, errors.New("更改密码失败")
	}
	return teacher, nil
}

func (w *Workflow) TeacherResetPasswordSms(username, mobile string) error {
	if username == "" {
		return errors.New("用户名为空")
	} else if mobile == "" {
		return errors.New("手机号为空")
	}
	teacher, err := w.model.GetTeacherByUsername(username)
	if err != nil || teacher.Mobile != mobile {
		return errors.New("用户名与手机号不匹配")
	}
	verifyCode, err := utils.GenerateVerifyCode(6)
	if err != nil {
		return errors.New("发送短信失败，请稍后重试")
	}
	err = w.redisClient.Set(fmt.Sprintf(model.REDIS_KEY_TEACHER_RESET_PASSWORD_VERIFY_CODE, teacher.Id.Hex()), verifyCode, 10*time.Minute).Err()
	if err != nil {
		return errors.New("发送短信失败，请稍后重试")
	}
	if err = w.SendResetPasswordSMS(teacher, verifyCode); err != nil {
		return errors.New("发送短信失败，请稍后重试")
	}
	return nil
}

func (w *Workflow) TeacherRestPasswordVerify(username, newPassword, verifyCode string) error {
	if username == "" {
		return errors.New("用户名为空")
	} else if newPassword == "" {
		return errors.New("新密码为空")
	} else if verifyCode == "" {
		return errors.New("验证码为空")
	}
	teacher, err := w.model.GetTeacherByUsername(username)
	if err != nil {
		return errors.New("咨询师不存在")
	}
	val, err := w.redisClient.Get(fmt.Sprintf(model.REDIS_KEY_TEACHER_RESET_PASSWORD_VERIFY_CODE, teacher.Id.Hex())).Result()
	if err != nil || err == redis.Nil || val != verifyCode {
		return errors.New("验证码错误或已过期")
	}
	w.redisClient.Del(fmt.Sprintf(model.REDIS_KEY_TEACHER_RESET_PASSWORD_VERIFY_CODE, teacher.Id.Hex()))
	encryptedPassword, err := utils.EncryptPassword(newPassword)
	if err != nil {
		return errors.New("新密码不符合要求")
	}
	teacher.EncryptedPassword = encryptedPassword
	if err = w.model.UpsertTeacher(teacher); err != nil {
		return errors.New("更改密码失败")
	}
	return nil
}

// 管理员登录
func (w *Workflow) AdminLogin(username string, password string) (*model.Admin, error) {
	if username == "" {
		return nil, errors.New("用户名为空")
	} else if password == "" {
		return nil, errors.New("密码为空")
	}
	admin, err := w.model.GetAdminByUsername(username)
	if err == nil {
		if admin.EncryptedPassword != "" {
			if utils.ValidatePassword(password, admin.EncryptedPassword) {
				return admin, nil
			}
		} else if password == admin.Password {
			if encryptedPassword, err := utils.EncryptPassword(password); err == nil {
				admin.EncryptedPassword = encryptedPassword
				w.model.UpsertAdmin(admin)
			}
			return admin, nil
		}
	}
	return nil, errors.New("用户名或密码不正确")
}

// 学生注册
func (w *Workflow) StudentRegister(username string, password string) (*model.Student, error) {
	if username == "" {
		return nil, errors.New("用户名为空")
	} else if password == "" {
		return nil, errors.New("密码为空")
	}
	if !utils.IsStudentId(username) {
		return nil, errors.New("请用学号注册")
	}
	if student, _ := w.model.GetStudentByUsername(username); student != nil {
		return nil, errors.New("该学号已被注册")
	}
	newStudent, err := w.model.AddStudent(username, password)
	if err != nil {
		return nil, errors.New("注册失败，请联系管理员")
	}
	archive, err := w.model.GetArchiveByStudentUsername(newStudent.Username)
	if err == nil {
		newStudent.ArchiveCategory = archive.ArchiveCategory
		newStudent.ArchiveNumber = archive.ArchiveNumber
		w.model.UpsertStudent(newStudent)
	}
	return newStudent, nil
}

// 更新session
func (w *Workflow) UpdateSession(userId string, userType int) (map[string]interface{}, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	}
	result := make(map[string]interface{})
	switch userType {
	case model.USER_TYPE_ADMIN:
		admin, err := w.model.GetAdminById(userId)
		if err != nil || admin.UserType != userType {
			return nil, errors.New("请重新登录")
		}
		result["user_id"] = admin.Id.Hex()
		result["username"] = admin.Username
		result["user_type"] = admin.UserType
		result["fullname"] = admin.Fullname
	case model.USER_TYPE_TEACHER:
		teacher, err := w.model.GetTeacherById(userId)
		if err != nil || teacher.UserType != userType {
			return nil, errors.New("请重新登录")
		}
		result["user_id"] = teacher.Id.Hex()
		result["username"] = teacher.Username
		result["user_type"] = teacher.UserType
		result["fullname"] = teacher.Fullname
	case model.USER_TYPE_STUDENT:
		student, err := w.model.GetStudentById(userId)
		if err != nil || student.UserType != userType {
			return nil, errors.New("请重新登录")
		}
		result["user_id"] = student.Id.Hex()
		result["username"] = student.Username
		result["user_type"] = student.UserType
		result["fullname"] = student.Fullname
	default:
		return nil, errors.New("请重新登录")
	}
	return result, nil
}

// 获取学生
func (w *Workflow) GetStudentById(userId string) (*model.Student, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	}
	student, err := w.model.GetStudentById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	}
	return student, nil
}

func (w *Workflow) GetStudentByUsername(studentUsername string) (*model.Student, error) {
	if studentUsername == "" {
		return nil, errors.New("学号为空")
	}
	student, err := w.model.GetStudentByUsername(studentUsername)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	return student, nil
}

// 获取咨询师
func (w *Workflow) GetTeacherByUsername(username string) (*model.Teacher, error) {
	if len(username) == 0 {
		return nil, errors.New("请先登录")
	}
	teacher, err := w.model.GetTeacherByUsername(username)
	if err != nil {
		return nil, errors.New("请先登录")
	}
	return teacher, nil
}

// 获取咨询师
func (w *Workflow) GetTeacherById(userId string) (*model.Teacher, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	}
	teacher, err := w.model.GetTeacherById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	}
	return teacher, nil
}

// 获取管理员
func (w *Workflow) GetAdminById(userId string) (*model.Admin, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	}
	admin, err := w.model.GetAdminById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	}
	return admin, nil
}

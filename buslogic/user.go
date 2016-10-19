package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"errors"
)

const (
	AdminDefaultPassword   = "THXLFZZX"
	TeacherDefaultPassword = "thxlfzzx"
)

// 学生登录
func (w *Workflow) StudentLogin(username string, password string) (*model.Student, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名为空")
	} else if len(password) == 0 {
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
	if len(username) == 0 {
		return nil, errors.New("用户名为空")
	} else if len(password) == 0 {
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

// 管理员登录
func (w *Workflow) AdminLogin(username string, password string) (*model.Admin, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名为空")
	} else if len(password) == 0 {
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
	if len(username) == 0 {
		return nil, errors.New("用户名为空")
	} else if len(password) == 0 {
		return nil, errors.New("密码为空")
	}
	if !utils.IsStudentId(username) {
		return nil, errors.New("请用学号注册")
	}
	if _, err := w.model.GetStudentByUsername(username); err == nil {
		return nil, errors.New("该学号已被注册")
	}
	newStudent, err := w.model.AddStudent(username, password)
	if err != nil {
		return nil, errors.New("注册失败，请联系管理员")
	}
	archive, err := w.model.GetArchiveByStudentUsername(newStudent.Username)
	if err == nil && archive != nil {
		newStudent.ArchiveCategory = archive.ArchiveCategory
		newStudent.ArchiveNumber = archive.ArchiveNumber
		w.model.UpsertStudent(newStudent)
	}
	return newStudent, nil
}

// 获取学生
func (w *Workflow) GetStudentById(userId string) (*model.Student, error) {
	if len(userId) == 0 {
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
	if len(userId) == 0 {
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

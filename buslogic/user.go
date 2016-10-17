package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/util"
	"errors"
	"strings"
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
	if err == nil && strings.EqualFold(student.Password, password) {
		return student, nil
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
	if err == nil && (strings.EqualFold(password, teacher.Password) ||
		(teacher.UserType == model.TEACHER && strings.EqualFold(password, TeacherDefaultPassword))) {
		return teacher, nil
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
	if err == nil && (strings.EqualFold(password, admin.Password) ||
		(admin.UserType == model.ADMIN && strings.EqualFold(password, AdminDefaultPassword))) {
		return admin, nil
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
	if !util.IsStudentId(username) {
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

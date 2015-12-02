package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"strings"
)

const (
	AdminDefaultPassword   = "THXLFZZX"
	TeacherDefaultPassword = "thxlfzzx"
)

type UserLogic struct {
}

// 学生登录
func (ul *UserLogic) StudentLogin(username string, password string) (*models.Student, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名为空")
	} else if len(password) == 0 {
		return nil, errors.New("密码为空")
	}
	student, err := models.GetStudentByUsername(username)
	if err == nil && strings.EqualFold(student.Password, password) {
		return student, nil
	}
	return nil, errors.New("用户名或密码不正确")
}

// 咨询师及管理员登录
func (ul *UserLogic) TeacherLogin(username string, password string) (*models.Teacher, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名为空")
	} else if len(password) == 0 {
		return nil, errors.New("密码为空")
	}
	teacher, err := models.GetTeacherByUsername(username)
	if err == nil && (strings.EqualFold(teacher.Password, password) ||
		(teacher.UserType == models.ADMIN && strings.EqualFold(teacher.Password, AdminDefaultPassword)) ||
		(teacher.UserType == models.TEACHER && strings.EqualFold(teacher.Password, TeacherDefaultPassword))) {
		return teacher, nil
	}
	return nil, errors.New("用户名或密码不正确")
}

// 学生注册
func (ul *UserLogic) StudentRegister(username string, password string) (*models.Student, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名为空")
	} else if len(password) == 0 {
		return nil, errors.New("密码为空")
	}
	if !utils.IsStudentId(username) {
		return nil, errors.New("请用学号注册")
	}
	if _, err := models.GetStudentByUsername(username); err == nil {
		return nil, errors.New("该学号已被注册")
	}
	newStudent, err := models.AddStudent(username, password)
	if err != nil {
		return nil, errors.New("注册失败，请联系管理员")
	}
	return newStudent, nil
}

// 获取学生
func (ul *UserLogic) GetStudentByUsername(username string) (*models.Student, error) {
	if len(username) == 0 {
		return nil, errors.New("请先登录")
	}
	student, err := models.GetStudentByUsername(username)
	if err != nil {
		return nil, errors.New("请先登录")
	}
	return student, nil
}

// 获取咨询师或管理员
func (ul *UserLogic) GetTeacherByUsername(username string) (*models.Teacher, error) {
	if len(username) == 0 {
		return nil, errors.New("请先登录")
	}
	teacher, err := models.GetTeacherByUsername(username)
	if err != nil {
		return nil, errors.New("请先登录")
	}
	return teacher, nil
}

// 获取咨询师或管理员
func (ul *UserLogic) GetTeacherById(userId string) (*models.Teacher, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	}
	teacher, err := models.GetTeacherById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	}
	return teacher, nil
}
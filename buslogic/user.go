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

// 咨询师登录
func (ul *UserLogic) TeacherLogin(username string, password string) (*models.Teacher, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名为空")
	} else if len(password) == 0 {
		return nil, errors.New("密码为空")
	}
	teacher, err := models.GetTeacherByUsername(username)
	if err == nil && (strings.EqualFold(password, teacher.Password) ||
		(teacher.UserType == models.TEACHER && strings.EqualFold(password, TeacherDefaultPassword))) {
		return teacher, nil
	}
	return nil, errors.New("用户名或密码不正确")
}

// 管理员登录
func (ul *UserLogic) AdminLogin(username string, password string) (*models.Admin, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名为空")
	} else if len(password) == 0 {
		return nil, errors.New("密码为空")
	}
	admin, err := models.GetAdminByUsername(username)
	if err == nil && (strings.EqualFold(password, admin.Password) ||
		(admin.UserType == models.ADMIN && strings.EqualFold(password, AdminDefaultPassword))) {
		return admin, nil
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
	archive, err := models.GetArchiveByStudentUsername(newStudent.Username)
	if err == nil && archive != nil {
		newStudent.ArchiveCategory = archive.ArchiveCategory
		newStudent.ArchiveNumber = archive.ArchiveNumber
		models.UpsertStudent(newStudent)
	}
	return newStudent, nil
}

// 获取学生
func (ul *UserLogic) GetStudentById(userId string) (*models.Student, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	}
	student, err := models.GetStudentById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	}
	return student, nil
}

func (ul *UserLogic) GetStudentByUsername(studentUsername string) (*models.Student, error) {
	if studentUsername == "" {
		return nil, errors.New("学号为空")
	}
	student, err := models.GetStudentByUsername(studentUsername)
	if err != nil {
		return nil, errors.New("学生未注册")
	}
	return student, nil
}

// 获取咨询师
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

// 获取咨询师
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

// 获取管理员
func (ul *UserLogic) GetAdminById(userId string) (*models.Admin, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	}
	admin, err := models.GetAdminById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	}
	return admin, nil
}

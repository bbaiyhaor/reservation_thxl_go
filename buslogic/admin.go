package buslogic
import (
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"errors"
"strings"
)

type AdminLogic struct  {
}

// 查找咨询师
// 查找顺序:全名 > 工号 > 手机号
func (al *AdminLogic) SearchTeacherByAdmin(teacherFullname string, teacherUsername string, teacherMobile string, userId string, userType models.UserType) (*models.User, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := models.GetTeacherById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	if !strings.EqualFold(teacherFullname, "") {
		teacher, err := models.GetTeacherByFullname(teacherFullname)
		if err == nil {
			return teacher, nil
		}
	}
	if !strings.EqualFold(teacherUsername, "") {
		teacher, err := models.GetTeacherByUsername(teacherUsername)
		if err == nil {
			return teacher, nil
		}
	}
	if !strings.EqualFold(teacherMobile, "") {
		teacher, err := models.GetTeacherByMobile(teacherMobile)
		if err == nil {
			return teacher, nil
		}
	}
	return nil, errors.New("用户不存在")
}
package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"errors"
	"strconv"
	"time"
)

// 学生预约咨询
func (w *Workflow) MakeReservationByStudent(reservationId string, sourceId string, startTime string,
	fullname string, gender string, birthday string, school string, grade string, currentAddress string,
	familyAddress string, mobile string, email string, experienceTime string, experienceLocation string,
	experienceTeacher string, fatherAge string, fatherJob string, fatherEdu string, motherAge string, motherJob string,
	motherEdu string, parentMarriage string, siginificant string, problem string,
	userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	} else if userType != model.USER_TYPE_STUDENT {
		return nil, errors.New("请重新登录")
	} else if reservationId == "" {
		return nil, errors.New("咨询已下架")
	} else if fullname == "" {
		return nil, errors.New("姓名为空")
	} else if gender != model.USER_GENDER_MALE && gender != model.USER_GENDER_FEMALE {
		return nil, errors.New("性别为空")
	} else if birthday == "" {
		return nil, errors.New("出生日期为空")
	} else if school == "" {
		return nil, errors.New("院系为空")
	} else if grade == "" {
		return nil, errors.New("年级为空")
	} else if currentAddress == "" {
		return nil, errors.New("现在住址为空")
	} else if familyAddress == "" {
		return nil, errors.New("家庭住址为空")
	} else if mobile == "" {
		return nil, errors.New("手机号为空")
	} else if email == "" {
		return nil, errors.New("邮箱为空")
	} else if problem == "" {
		return nil, errors.New("问题为空")
	} else if !utils.IsMobile(mobile) {
		return nil, errors.New("手机号格式不正确")
	} else if !utils.IsEmail(email) {
		return nil, errors.New("邮箱格式不正确")
	}
	student, err := w.model.GetStudentById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	} else if student.UserType != model.USER_TYPE_STUDENT {
		return nil, errors.New("请重新登录")
	}
	studentReservations, err := w.model.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	for _, r := range studentReservations {
		if r.Status == model.RESERVATION_STATUS_RESERVATED && r.StartTime.After(time.Now()) {
			return nil, errors.New("你好！你已有一个咨询预约，请完成这次咨询后再预约下一次，或致电62782007取消已有预约。")
		}
	}
	var reservation *model.Reservation
	if sourceId == "" {
		// Source为ADD，无SourceId：直接预约
		reservation, err = w.model.GetReservationById(reservationId)
		if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
			return nil, errors.New("咨询已下架")
		} else if reservation.StartTime.Before(time.Now()) {
			return nil, errors.New("咨询已过期")
		} else if reservation.Status != model.RESERVATION_STATUS_AVAILABLE {
			return nil, errors.New("咨询已被预约")
		} else if student.BindedTeacherId != "" && student.BindedTeacherId != reservation.TeacherId {
			return nil, errors.New("只能预约匹配咨询师")
		}
	} else if reservationId == sourceId {
		// Source为TIMETABLE且未被预约
		timedReservation, err := w.model.GetTimedReservationById(sourceId)
		if err != nil || timedReservation.Status == model.RESERVATION_STATUS_DELETED {
			return nil, errors.New("咨询已下架")
		}
		start, err := time.ParseInLocation("2006-01-02 15:04", startTime, time.Local)
		if err != nil {
			return nil, errors.New("开始时间格式错误")
		} else if start.Before(time.Now()) {
			return nil, errors.New("咨询已过期")
		} else if start.Format("15:04") != timedReservation.StartTime.Format("15:04") {
			return nil, errors.New("开始时间不匹配")
		} else if timedReservation.Timed[start.Format("2006-01-02")] {
			return nil, errors.New("咨询已被预约")
		} else if student.BindedTeacherId != "" && student.BindedTeacherId != timedReservation.TeacherId {
			return nil, errors.New("只能预约匹配咨询师")
		}
		end := utils.ConcatTime(start, timedReservation.EndTime)
		reservation, err = w.model.AddReservation(start, end, model.RESERVATION_SOURCE_TIMETABLE, timedReservation.Id.Hex(), timedReservation.TeacherId)
		if err != nil {
			return nil, errors.New("获取数据失败")
		}
		timedReservation.Timed[start.Format("2006-01-02")] = true
		if w.model.UpsertTimedReservation(timedReservation) != nil {
			return nil, errors.New("获取数据失败")
		}
	} else {
		return nil, errors.New("咨询已被预约")
	}
	// 更新学生信息
	student.Fullname = fullname
	student.Gender = gender
	student.Birthday = birthday
	student.School = school
	student.Grade = grade
	student.CurrentAddress = currentAddress
	student.FamilyAddress = familyAddress
	student.Mobile = mobile
	student.Email = email
	student.Experience.Time = experienceTime
	student.Experience.Location = experienceLocation
	student.Experience.Teacher = experienceTeacher
	student.FatherAge = fatherAge
	student.FatherJob = fatherJob
	student.FatherEdu = fatherEdu
	student.MotherAge = motherAge
	student.MotherJob = motherJob
	student.MotherEdu = motherEdu
	student.ParentMarriage = parentMarriage
	student.Significant = siginificant
	student.Problem = problem
	student.BindedTeacherId = reservation.TeacherId
	if w.model.UpsertStudent(student) != nil {
		return nil, errors.New("获取数据失败")
	}
	// 更新咨询信息
	reservation.StudentId = student.Id.Hex()
	reservation.Status = model.RESERVATION_STATUS_RESERVATED
	if w.model.UpsertReservation(reservation) != nil {
		return nil, errors.New("获取数据失败")
	}
	// send success sms
	w.SendSuccessSMS(reservation)
	return reservation, nil
}

// 学生拉取反馈
func (w *Workflow) GetFeedbackByStudent(reservationId string, sourceId string,
	userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	} else if userType != model.USER_TYPE_STUDENT {
		return nil, errors.New("请重新登录")
	} else if reservationId == "" {
		return nil, errors.New("咨询已下架")
	} else if reservationId == sourceId {
		return nil, errors.New("咨询未被预约，不能反馈")
	}
	student, err := w.model.GetStudentById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	} else if student.UserType != model.USER_TYPE_STUDENT {
		return nil, errors.New("请重新登录")
	}
	reservation, err := w.model.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == model.RESERVATION_STATUS_AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if reservation.StudentId != student.Id.Hex() {
		return nil, errors.New("只能反馈本人预约的咨询")
	}
	return reservation, nil
}

// 学生反馈
func (w *Workflow) SubmitFeedbackByStudent(reservationId string, sourceId string, scores []int,
	userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, errors.New("请先登录")
	} else if userType != model.USER_TYPE_STUDENT {
		return nil, errors.New("请重新登录")
	} else if reservationId == "" {
		return nil, errors.New("咨询已下架")
	} else if len(scores) != model.RESERVATION_STUDENT_FEEDBACK_SCORES_LENGTH {
		return nil, errors.New("请完整填写反馈")
	} else if reservationId == sourceId {
		return nil, errors.New("咨询未被预约，不能反馈")
	}
	student, err := w.model.GetStudentById(userId)
	if err != nil {
		return nil, errors.New("请先登录")
	} else if student.UserType != model.USER_TYPE_STUDENT {
		return nil, errors.New("请重新登录")
	}
	reservation, err := w.model.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == model.RESERVATION_STATUS_AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if reservation.StudentId != student.Id.Hex() {
		return nil, errors.New("只能反馈本人预约的咨询")
	}
	reservation.StudentFeedback = model.StudentFeedback{
		Scores: scores,
	}
	if w.model.UpsertReservation(reservation) != nil {
		return nil, errors.New("获取数据失败")
	}
	return reservation, nil
}

//
func (w *Workflow) ExportStudentInfoToFile(student *model.Student, filename string) error {
	data := make([][]string, 0)
	data = append(data, []string{"档案分类", student.ArchiveCategory, "档案编号", student.ArchiveNumber})
	// 学生基本信息
	data = append(data, []string{"学号", student.Username})
	data = append(data, []string{"姓名", student.Fullname})
	data = append(data, []string{"性别", student.Gender})
	data = append(data, []string{"出生日期", student.Birthday})
	data = append(data, []string{"系别", student.School})
	data = append(data, []string{"年级", student.Grade})
	data = append(data, []string{"现住址", student.CurrentAddress})
	data = append(data, []string{"家庭住址", student.FamilyAddress})
	data = append(data, []string{"联系电话", student.Mobile})
	data = append(data, []string{"Email", student.Email})
	if !student.Experience.IsEmpty() {
		data = append(data, []string{"咨询经历", "时间", student.Experience.Time, "地点", student.Experience.Location,
			"咨询师姓名", student.Experience.Teacher})
	} else {
		data = append(data, []string{"咨询经历", "无"})
	}
	data = append(data, []string{"父亲", "年龄", student.FatherAge, "职业", student.FatherJob, "学历", student.FatherEdu})
	data = append(data, []string{"母亲", "年龄", student.MotherAge, "职业", student.MotherJob, "学历", student.MotherEdu})
	data = append(data, []string{"父母婚姻状况", student.ParentMarriage})
	data = append(data, []string{"在近三个月里，是否发生了对你有重大意义的事（如亲友的死亡、法律诉讼、失恋等）？", student.Significant})
	data = append(data, []string{"你现在需要接受帮助的主要问题是什么？", student.Problem})
	bindedTeacher, err := w.model.GetTeacherById(student.BindedTeacherId)
	if err != nil {
		data = append(data, []string{"匹配咨询师", "无"})
	} else {
		data = append(data, []string{"匹配咨询师", bindedTeacher.Username, bindedTeacher.Fullname})
	}
	data = append(data, []string{"危机等级", strconv.Itoa(student.CrisisLevel)})
	data = append(data, []string{""})
	data = append(data, []string{""})

	//咨询小结
	if reservations, err := w.model.GetReservationsByStudentId(student.Id.Hex()); err == nil {
		for i, r := range reservations {
			teacher, err := w.model.GetTeacherById(r.TeacherId)
			if err != nil {
				continue
			}
			data = append(data, []string{"咨询小结" + strconv.Itoa(i+1)})
			data = append(data, []string{"咨询师", teacher.Username, teacher.Fullname})
			data = append(data, []string{"咨询日期", r.StartTime.Format("2006-01-02")})
			if !r.TeacherFeedback.IsEmpty() {
				data = append(data, []string{"评估分类", model.FeedbackAllCategory[r.TeacherFeedback.Category]})
				participants := []string{"出席人员"}
				for j := 0; j < len(r.TeacherFeedback.Participants); j++ {
					if r.TeacherFeedback.Participants[j] > 0 {
						participants = append(participants, model.PARTICIPANTS[j])
					}
				}
				data = append(data, participants)

				if r.TeacherFeedback.Emphasis > 0 {
					data = append(data, []string{"重点明细", "是"})
				} else {
					data = append(data, []string{"重点明细", "否"})
				}
				severity := []string{"严重程度"}
				if len(r.TeacherFeedback.Severity) == len(model.SEVERITY) {
					for i := 0; i < len(r.TeacherFeedback.Severity); i++ {
						if r.TeacherFeedback.Severity[i] > 0 {
							severity = append(severity, model.SEVERITY[i])
						}
					}
				}
				data = append(data, severity)
				medicalDiagnosis := []string{"疑似或明确的医疗诊断"}
				if len(r.TeacherFeedback.MedicalDiagnosis) == len(model.MEDICAL_DIAGNOSIS) {
					for i := 0; i < len(r.TeacherFeedback.MedicalDiagnosis); i++ {
						if r.TeacherFeedback.MedicalDiagnosis[i] > 0 {
							medicalDiagnosis = append(medicalDiagnosis, model.MEDICAL_DIAGNOSIS[i])
						}
					}
				}
				data = append(data, medicalDiagnosis)
				crisis := []string{"危急情况"}
				if len(r.TeacherFeedback.Crisis) == len(model.CRISIS) {
					for i := 0; i < len(r.TeacherFeedback.Crisis); i++ {
						if r.TeacherFeedback.Crisis[i] > 0 {
							crisis = append(crisis, model.CRISIS[i])
						}
					}
				}
				data = append(data, crisis)

				data = append(data, []string{"咨询记录", r.TeacherFeedback.Record})
			}
			if !r.StudentFeedback.IsEmpty() {
				scores := []string{"来访者反馈"}
				for _, s := range r.StudentFeedback.Scores {
					scores = append(scores, strconv.Itoa(s))
				}
				data = append(data, scores)
			}
		}
		data = append(data, []string{""})
	}
	if err := utils.WriteToCSV(data, filename); err != nil {
		return err
	}
	return nil
}

func (w *Workflow) WrapSimpleStudent(student *model.Student) map[string]interface{} {
	var result = make(map[string]interface{})
	if student == nil {
		return result
	}
	result["id"] = student.Id.Hex()
	result["username"] = student.Username
	result["user_type"] = student.UserType
	result["binded_teacher_id"] = student.BindedTeacherId
	result["fullname"] = student.Fullname
	result["gender"] = student.Gender
	result["birthday"] = student.Birthday
	result["school"] = student.School
	result["grade"] = student.Grade
	result["current_address"] = student.CurrentAddress
	result["family_address"] = student.FamilyAddress
	result["mobile"] = student.Mobile
	result["email"] = student.Email
	result["experience_time"] = student.Experience.Time
	result["experience_location"] = student.Experience.Location
	result["experience_teacher"] = student.Experience.Teacher
	result["father_age"] = student.FatherAge
	result["father_job"] = student.FatherJob
	result["father_edu"] = student.FatherEdu
	result["mother_age"] = student.MotherAge
	result["mother_job"] = student.MotherJob
	result["mother_edu"] = student.MotherEdu
	result["parent_marriage"] = student.ParentMarriage
	result["significant"] = student.Significant
	result["problem"] = student.Problem
	return result
}

func (w *Workflow) WrapStudent(student *model.Student) map[string]interface{} {
	var result = make(map[string]interface{})
	if student == nil {
		return result
	}
	result["id"] = student.Id.Hex()
	result["username"] = student.Username
	result["user_type"] = student.UserType
	result["binded_teacher_id"] = student.BindedTeacherId
	if bindedTeacher, err := w.model.GetTeacherById(student.BindedTeacherId); err == nil {
		result["binded_teacher_username"] = bindedTeacher.Username
		result["binded_teacher_fullname"] = bindedTeacher.Fullname
	} else {
		result["binded_teacher_username"] = ""
		result["binded_teacher_fullname"] = ""
	}
	result["archive_category"] = student.ArchiveCategory
	result["archive_number"] = student.ArchiveNumber
	result["crisis_level"] = student.CrisisLevel
	result["fullname"] = student.Fullname
	result["gender"] = student.Gender
	result["birthday"] = student.Birthday
	result["school"] = student.School
	result["grade"] = student.Grade
	result["current_address"] = student.CurrentAddress
	result["family_address"] = student.FamilyAddress
	result["mobile"] = student.Mobile
	result["email"] = student.Email
	result["experience_time"] = student.Experience.Time
	result["experience_location"] = student.Experience.Location
	result["experience_teacher"] = student.Experience.Teacher
	result["father_age"] = student.FatherAge
	result["father_job"] = student.FatherJob
	result["father_edu"] = student.FatherEdu
	result["mother_age"] = student.MotherAge
	result["mother_job"] = student.MotherJob
	result["mother_edu"] = student.MotherEdu
	result["parent_marriage"] = student.ParentMarriage
	result["significant"] = student.Significant
	result["problem"] = student.Problem
	return result
}

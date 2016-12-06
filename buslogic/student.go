package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	re "bitbucket.org/shudiwsh2009/reservation_thxl_go/rerror"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
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
		return nil, re.NewRErrorCode("student not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_STUDENT {
		return nil, re.NewRErrorCode("user is not student", nil, re.ERROR_NOT_AUTHORIZED)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ERROR_MISSING_PARAM, "reservation_id")
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
	} else if problem == "" {
		return nil, re.NewRErrorCodeContext("problem is empty", nil, re.ERROR_MISSING_PARAM, "problem")
	} else if !utils.IsMobile(mobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ERROR_FORMAT_MOBILE)
	} else if !utils.IsEmail(email) {
		return nil, re.NewRErrorCode("email format is wrong", nil, re.ERROR_FORMAT_EMAIL)
	}
	student, err := w.mongoClient.GetStudentById(userId)
	if err != nil || student.UserType != model.USER_TYPE_STUDENT {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	studentReservations, err := w.mongoClient.GetReservationsByStudentId(student.Id.Hex())
	if err != nil {
		return nil, re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	for _, r := range studentReservations {
		if r.Status == model.RESERVATION_STATUS_RESERVATED && r.StartTime.After(time.Now()) {
			return nil, re.NewRErrorCode("already have reservation", nil, re.ERROR_STUDENT_ALREADY_HAVE_RESERVATION)
		}
	}
	var reservation *model.Reservation
	if sourceId == "" {
		// Source为ADD，无SourceId：直接预约
		reservation, err = w.mongoClient.GetReservationById(reservationId)
		if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
			return nil, re.NewRErrorCode("fail to get reservation", nil, re.ERROR_DATABASE)
		} else if reservation.StartTime.Before(time.Now()) {
			return nil, re.NewRErrorCode("cannot make outdated reservation", nil, re.ERROR_STUDENT_MAKE_OUTDATED_RESERVATION)
		} else if reservation.Status != model.RESERVATION_STATUS_AVAILABLE {
			return nil, re.NewRErrorCode("cannot make reservated reservation", nil, re.ERROR_STUDENT_MAKE_RESERVATED_RESERVATION)
		} else if student.BindedTeacherId != "" && student.BindedTeacherId != reservation.TeacherId {
			return nil, re.NewRErrorCode("only make binded teacher reservation", nil, re.ERROR_STUDENT_MAKE_NOT_BINDED_TEACHER_RESERVATION)
		}
	} else if reservationId == sourceId {
		// Source为TIMETABLE且未被预约
		timedReservation, err := w.mongoClient.GetTimedReservationById(sourceId)
		if err != nil || timedReservation.Status == model.RESERVATION_STATUS_DELETED {
			return nil, re.NewRErrorCode("fail to get timetable", nil, re.ERROR_DATABASE)
		}
		start, err := time.ParseInLocation("2006-01-02 15:04", startTime, time.Local)
		if err != nil {
			return nil, re.NewRErrorCodeContext("start time is not valid", err, re.ERROR_INVALID_PARAM, "start_time")
		} else if start.Before(time.Now()) {
			return nil, re.NewRErrorCode("cannot make outdated reservation", nil, re.ERROR_STUDENT_MAKE_OUTDATED_RESERVATION)
		} else if start.Format("15:04") != timedReservation.StartTime.Format("15:04") {
			return nil, re.NewRErrorCode("start time mismatch", nil, re.ERROR_START_TIME_MISMATCH)
		} else if timedReservation.Timed[start.Format("2006-01-02")] {
			return nil, re.NewRErrorCode("cannot make reservated reservation", nil, re.ERROR_STUDENT_MAKE_RESERVATED_RESERVATION)
		} else if student.BindedTeacherId != "" && student.BindedTeacherId != timedReservation.TeacherId {
			return nil, re.NewRErrorCode("only make binded teacher reservation", nil, re.ERROR_STUDENT_MAKE_NOT_BINDED_TEACHER_RESERVATION)
		}
		end := utils.ConcatTime(start, timedReservation.EndTime)
		reservation = &model.Reservation{
			StartTime:       start,
			EndTime:         end,
			Status:          model.RESERVATION_STATUS_RESERVATED,
			Source:          model.RESERVATION_SOURCE_TIMETABLE,
			SourceId:        timedReservation.Id.Hex(),
			TeacherId:       timedReservation.TeacherId,
			StudentFeedback: model.StudentFeedback{},
			TeacherFeedback: model.TeacherFeedback{},
		}
		timedReservation.Timed[start.Format("2006-01-02")] = true
		if err = w.mongoClient.InsertReservationAndUpdateTimedReservation(reservation, timedReservation); err != nil {
			return nil, re.NewRErrorCode("fail to insert reservation and update timetable", err, re.ERROR_DATABASE)
		}
	} else {
		return nil, re.NewRErrorCode("cannot make reservated reservation", nil, re.ERROR_STUDENT_MAKE_RESERVATED_RESERVATION)
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
	student.ParentInfo = model.ParentInfo{
		FatherAge:      fatherAge,
		FatherJob:      fatherJob,
		FatherEdu:      fatherEdu,
		MotherAge:      motherAge,
		MotherJob:      motherJob,
		MotherEdu:      motherEdu,
		ParentMarriage: parentMarriage,
	}
	student.Significant = siginificant
	student.Problem = problem
	student.BindedTeacherId = reservation.TeacherId
	// 更新咨询信息
	reservation.StudentId = student.Id.Hex()
	reservation.Status = model.RESERVATION_STATUS_RESERVATED
	if err = w.mongoClient.UpdateReservationAndStudent(reservation, student); err != nil {
		return nil, re.NewRErrorCode("fail to update reservation and student", err, re.ERROR_DATABASE)
	}
	// send success sms
	w.SendSuccessSMS(reservation)
	return reservation, nil
}

// 学生拉取反馈
func (w *Workflow) GetFeedbackByStudent(reservationId string, sourceId string,
	userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("student not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_STUDENT {
		return nil, re.NewRErrorCode("user is not student", nil, re.ERROR_NOT_AUTHORIZED)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ERROR_MISSING_PARAM, "reservation_id")
	} else if reservationId == sourceId {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
	}
	student, err := w.mongoClient.GetStudentById(userId)
	if err != nil || student.UserType != model.USER_TYPE_STUDENT {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ERROR_DATABASE)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ERROR_FEEDBACK_FUTURE_RESERVATION)
	} else if reservation.Status == model.RESERVATION_STATUS_AVAILABLE {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
	} else if reservation.StudentId != student.Id.Hex() {
		return nil, re.NewRErrorCode("cannot get feedback of other one's reservation", nil, re.ERROR_FEEDBACK_OTHER_RESERVATION)
	}
	return reservation, nil
}

// 学生反馈
func (w *Workflow) SubmitFeedbackByStudent(reservationId string, sourceId string, scores []int,
	userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("student not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_STUDENT {
		return nil, re.NewRErrorCode("user is not student", nil, re.ERROR_NOT_AUTHORIZED)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ERROR_MISSING_PARAM, "reservation_id")
	} else if len(scores) != model.RESERVATION_STUDENT_FEEDBACK_SCORES_LENGTH {
		return nil, re.NewRErrorCodeContext("scores is not valid", nil, re.ERROR_INVALID_PARAM, "scores")
	} else if reservationId == sourceId {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
	}
	student, err := w.mongoClient.GetStudentById(userId)
	if err != nil || student.UserType != model.USER_TYPE_STUDENT {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation.Status == model.RESERVATION_STATUS_DELETED {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ERROR_DATABASE)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ERROR_FEEDBACK_FUTURE_RESERVATION)
	} else if reservation.Status == model.RESERVATION_STATUS_AVAILABLE {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ERROR_FEEDBACK_AVAILABLE_RESERVATION)
	} else if reservation.StudentId != student.Id.Hex() {
		return nil, re.NewRErrorCode("cannot get feedback of other one's reservation", nil, re.ERROR_FEEDBACK_OTHER_RESERVATION)
	}
	reservation.StudentFeedback = model.StudentFeedback{
		Scores: scores,
	}
	if err = w.mongoClient.UpdateReservation(reservation); err != nil {
		return nil, re.NewRErrorCode("fail to update reservation", err, re.ERROR_DATABASE)
	}
	return reservation, nil
}

func (w *Workflow) ExportStudentInfoToFile(student *model.Student, path string) error {
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
	data = append(data, []string{"父亲", "年龄", student.ParentInfo.FatherAge, "职业", student.ParentInfo.FatherJob, "学历", student.ParentInfo.FatherEdu})
	data = append(data, []string{"母亲", "年龄", student.ParentInfo.MotherAge, "职业", student.ParentInfo.MotherJob, "学历", student.ParentInfo.MotherEdu})
	data = append(data, []string{"父母婚姻状况", student.ParentInfo.ParentMarriage})
	data = append(data, []string{"在近三个月里，是否发生了对你有重大意义的事（如亲友的死亡、法律诉讼、失恋等）？", student.Significant})
	data = append(data, []string{"你现在需要接受帮助的主要问题是什么？", student.Problem})
	bindedTeacher, err := w.mongoClient.GetTeacherById(student.BindedTeacherId)
	if err != nil {
		data = append(data, []string{"匹配咨询师", "无"})
	} else {
		data = append(data, []string{"匹配咨询师", bindedTeacher.Username, bindedTeacher.Fullname})
	}
	data = append(data, []string{"危机等级", strconv.Itoa(student.CrisisLevel)})
	data = append(data, []string{""})
	data = append(data, []string{""})

	//咨询小结
	if reservations, err := w.mongoClient.GetReservationsByStudentId(student.Id.Hex()); err == nil {
		for i, r := range reservations {
			teacher, err := w.mongoClient.GetTeacherById(r.TeacherId)
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
	if err := utils.WriteToCSV(data, path); err != nil {
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
	result["father_age"] = student.ParentInfo.FatherAge
	result["father_job"] = student.ParentInfo.FatherJob
	result["father_edu"] = student.ParentInfo.FatherEdu
	result["mother_age"] = student.ParentInfo.MotherAge
	result["mother_job"] = student.ParentInfo.MotherJob
	result["mother_edu"] = student.ParentInfo.MotherEdu
	result["parent_marriage"] = student.ParentInfo.ParentMarriage
	result["significant"] = student.Significant
	result["problem"] = student.Problem
	return result
}

func (w *Workflow) WrapStudent(student *model.Student) map[string]interface{} {
	result := w.WrapSimpleStudent(student)
	if student == nil {
		return result
	}
	result["archive_category"] = student.ArchiveCategory
	result["archive_number"] = student.ArchiveNumber
	result["crisis_level"] = student.CrisisLevel
	if student.BindedTeacherId != "" {
		if bindedTeacher, err := w.mongoClient.GetTeacherById(student.BindedTeacherId); err == nil {
			result["binded_teacher_username"] = bindedTeacher.Username
			result["binded_teacher_fullname"] = bindedTeacher.Fullname
		} else {
			result["binded_teacher_username"] = ""
			result["binded_teacher_fullname"] = ""
		}
	} else {
		result["binded_teacher_username"] = ""
		result["binded_teacher_fullname"] = ""
	}
	return result
}

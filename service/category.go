package service

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"net/http"
)

func (s *Service) GetFeedbackCategories(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}

	result["first_category"] = model.FeedbackFirstCategory
	result["second_category"] = model.FeedbackSecondCategory

	return result
}

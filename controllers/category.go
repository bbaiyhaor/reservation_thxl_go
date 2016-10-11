package controllers

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/models"
	"net/http"
)

func GetFeedbackCategories(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}

	result["first_category"] = models.FeedbackFirstCategory
	result["second_category"] = models.FeedbackSecondCategory

	return result
}

package controllers

import (
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"net/http"
)

func GetFeedbackCategories(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}

	var first = map[string]interface{}{
		"A": "A 学业问题",
		"B": "B 情感问题",
		"C": "C 人际问题",
		"D": "D 发展问题",
		"E": "E 情绪问题",
		"F": "F 身心与行为问题",
		"G": "G 危机干预",
		"H": "H 心理测验",
		"I": "I 其他",
		"Y": "Y 团体辅导",
		"Z": "Z 个体心理督导",
	}
	result["first_category"] = first
	var second = map[string]interface{}{
		"A": map[string]interface{}{
			"A1": "A1 学业成就困扰",
			"A2": "A2 专业认同困扰",
			"A3": "A3 缓考评估",
			"A4": "A4 休学复学评估",
		},
		"B": map[string]interface{}{
			"B1": "B1 恋爱困扰",
			"B2": "B2 性困扰",
			"B3": "B3 性取向",
		},
		"C": map[string]interface{}{
			"C1": "C1 同伴人际",
			"C2": "C2 家庭人际",
			"C3": "C3 与辅导员人际",
			"C4": "C4 与教师人际",
		},
		"D": map[string]interface{}{
			"D1": "D1 就业困扰",
			"D2": "D2 事业探索",
			"D3": "D3 价值感与意义感",
			"D4": "D4 完美情结",
		},
		"E": map[string]interface{}{
			"E1": "E1 焦虑情绪",
			"E2": "E2 抑郁情绪",
			"E3": "E3 焦虑抑郁情绪",
		},
		"F": map[string]interface{}{
			"F1": "F1 睡眠问题",
			"F2": "F2 进食问题",
			"F3": "F3 身心问题",
			"F4": "F4 电脑依赖",
			"F5": "F5 强迫问题",
			"F6": "F6 品行问题",
		},
		"G": map[string]interface{}{
			"G1": "G1 应激状态干预",
			"G2": "G2 精神障碍发作期干预",
			"G3": "G3 精神障碍康复期干预",
			"G4": "G4 创伤后干预",
		},
		"H": map[string]interface{}{
			"H1": "H1 人格测验与反馈",
			"H2": "H2 情绪测验与反馈",
			"H3": "H3 学业测验与反馈",
			"H4": "H4 职业测验与反馈",
		},
		"I": map[string]interface{}{
			"I1": "I1 躯体疾病转介",
			"I2": "I2 严重心理问题转介",
			"I3": "I3 转介至学习发展中心",
			"I4": "I4 转介至就业指导中心",
			"I5": "I5 反映学生情况",
		},
		"Y": map[string]interface{}{
			"Y1": "Y1 学习压力团体",
			"Y2": "Y2 人际关系团体",
			"Y3": "Y3 恋爱情感团体",
			"Y4": "Y4 辅导员团体",
		},
		"Z": map[string]interface{}{
			"Z1": "Z1 个体心理督导",
		},
	}
	result["second_category"] = second

	return result
}

package define

import "asset-management/app/model"

type LogInfo struct {
	Method       string           `json:"method"`
	URL          string           `json:"url"`
	Status       int              `json:"status"`
	ErrorCode    int              `json:"error_code"`
	ErrorMessage string           `json:"error_message"`
	UserID       uint             `json:"user_id"`
	Username     string           `json:"username"`
	EntityID     uint             `json:"entity_id"`
	DepartmentID uint             `json:"department_id"`
	Time         *model.ModelTime `json:"time"`
	Level        string           `json:"level"`
	Message      string           `json:"message"`
}

type LogListResponse struct {
	LogList  []*LogInfo `json:"log_list"`
	AllCount uint       `json:"all_count"`
}

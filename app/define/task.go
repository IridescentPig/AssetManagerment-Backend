package define

import "asset-management/app/model"

type CreateTaskReq struct {
	TaskType        uint             `json:"task_type" binding:"gte=0,lte=3"`
	TaskDescription string           `json:"task_description"`
	TargetID        uint             `json:"target_id" binding:"gte=0"`
	AssetList       []ExpireAssetReq `json:"asset_list" binding:"gt=0,dive,gt=0"`
}

type TaskInfo struct {
	ID              uint           `json:"task_id"`
	TaskType        uint           `json:"task_type"` // 0领用、1退库、2维保、3转移
	TaskDescription string         `json:"task_description"`
	UserID          uint           `json:"user_id"`
	UserName        string         `json:"username"`
	TargetID        uint           `json:"target_id"`
	TargetName      string         `json:"target_name"`
	DepartmentID    uint           `json:"department_id"`
	DepartmentName  string         `json:"department"`
	AssetList       []*model.Asset `json:"asset_list"`
	State           uint           `json:"state"`
}

type TaskListResponse struct {
	TaskList []TaskInfo `json:"task_list"`
}

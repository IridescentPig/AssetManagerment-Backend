package define

type CreateTaskReq struct {
	TaskType        uint             `json:"task_type" binding:"gte=0,lte=3"`
	TaskDescription string           `json:"task_description"`
	TargetID        uint             `json:"target_id" binding:"gte=0"`
	AssetList       []ExpireAssetReq `json:"asset_list" binding:"gt=0,dive,gt=0"`
}

package define

type FeishuBindOrLoginRequest struct {
	Code string `json:"code"`
}

type FeishuCallBackReq struct {
	ActionType string `json:"action_type" binding:"omitempty,oneof=APPROVE REJECT"`
	Token      string `json:"token"`
	InstanceID string `json:"instance_id"`
	UserID     string `json:"user_id"`
}

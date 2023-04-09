package define

import "asset-management/app/model"

type CreateEntityReq struct {
	Name string `json:"name"`
}

type EntityManager struct {
	ManagerName string `json:"manager_name"`
	ManagerID   string `json:"manager_id"`
}

type EntityBasicInfo struct {
	EntityID   uint   `json:"entity_id"`
	EntityName string `json:"entity_name"`
}

/*
Response format of GET /entity/{entity_id}
*/
type EntityInfoResponse struct {
	EntityID    uint            `json:"entity_id"`
	EntityName  string          `json:"entity_name"`
	Description string          `json:"description"`
	CreatedTime model.ModelTime `json:"created_time"`
	ManagerList []EntityManager `json:"manager_list"`
}

/*
Response format of GET /entity/list
*/
type EntityListResponse struct {
	EntityList []EntityBasicInfo `json:"entity_list"`
}

type UserInfo struct {
	Username   string              `json:"username"`
	UserID     uint                `json:"user_id"`
	Indentity  []bool              `json:"indentity"`
	Lock       bool                `json:"lock"`
	Department DepartmentBasicInfo `json:"department"`
}

/*
Response format of GET /entity/{entity_id}/user/list
*/
type EntityUserListResponse struct {
	UserList []UserInfo `json:"user_list"`
}

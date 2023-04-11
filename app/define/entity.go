package define

import "asset-management/app/model"

type CreateEntityReq struct {
	EntityName string `json:"entity_name" binding:"required"`
}

type ManagerReq struct {
	Username string  `json:"username" binding:"required"`
	Password *string `json:"password"`
}

type ModifyEntityInfoReq struct {
	EntityName  *string `json:"entity_name"`
	Description *string `json:"description"`
}

type EntityManager struct {
	ManagerID   uint   `json:"manager_id" copier:"ID"`
	ManagerName string `json:"manager_name" copier:"UserName"`
}

type EntityBasicInfo struct {
	ID   uint   `json:"entity_id"`
	Name string `json:"entity_name"`
}

/*
Response format of GET /entity/{entity_id}
*/
type EntityInfoResponse struct {
	EntityID    uint             `json:"entity_id"`
	EntityName  string           `json:"entity_name"`
	Description string           `json:"description"`
	CreatedAt   *model.ModelTime `json:"created_at"`
	ManagerList []EntityManager  `json:"manager_list"`
}

/*
Response format of GET /entity/list
*/
type EntityListResponse struct {
	EntityList []EntityBasicInfo `json:"entity_list"`
}

type EntityUserInfo struct {
	UserName        string            `json:"username"`
	ID              uint              `json:"user_id"`
	Ban             bool              `json:"lock"`
	SystemSuper     bool              `json:"id3"`
	EntitySuper     bool              `json:"id2"`
	DepartmentSuper bool              `json:"id1"`
	IsEmployee      bool              `json:"id0"`
	Department      *model.Department `json:"department"`
}

/*
Response format of GET /entity/{entity_id}/user/list
*/
type EntityUserListResponse struct {
	UserList []EntityUserInfo `json:"user_list"`
}

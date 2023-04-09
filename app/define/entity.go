package define

import "asset-management/app/model"

type CreateEntityReq struct {
	EntityName string `json:"entity_name"`
}

type EntityManager struct {
	ManagerName string `json:"manager_name" copier:"UserName"`
	ManagerID   string `json:"manager_id" copier:"ID"`
}

type EntityBasicInfo struct {
	ID   uint   `json:"entity_id"`
	Name string `json:"entity_name"`
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

type EntityUserInfo struct {
	UserName        string            `json:"username"`
	ID              uint              `json:"user_id"`
	Ban             bool              `json:"lock"`
	SystemSuper     bool              `json:"id3"`
	EntitySuper     bool              `json:"id2"`
	DepartmentSuper bool              `json:"id1"`
	Employee        bool              `json:"id0" default:"true"`
	Department      *model.Department `json:"department"`
}

/*
Response format of GET /entity/{entity_id}/user/list
*/
type EntityUserListResponse struct {
	UserList []EntityUserInfo `json:"user_list"`
}

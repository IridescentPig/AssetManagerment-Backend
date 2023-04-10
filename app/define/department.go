package define

import "asset-management/app/model"

type CreateDepartmentReq struct {
	EntityID     uint
	DepartmentID uint
}

type DepartmentBasicInfo struct {
	ID       uint   `json:"department_id"`
	Name     string `json:"department_name"`
	ParentID uint   `json:"parent_id"`
}

type CreateDepartmentUserReq struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	DepartmentSuper bool   `json:"department_super"`
}

/*
Response format of GET /entity/{entity_id}/department/list

	and GET /entity/{entity_id}/department/{department_id}/department/list
*/
type DepartmentListResponse struct {
	DepartmentList []DepartmentBasicInfo `json:"department_list"`
}

type DepartmentInfo struct {
	DepartmentBasicInfo
	Entity EntityBasicInfo `json:"entity"`
}

/*
Response format of GET /entity/{entity_id}/department/{department_id}
*/
type DepartmentInfoResponse struct {
	Department DepartmentInfo `json:"department"`
}

type DepartmentUserInfo struct {
	UserName        string            `json:"username"`
	ID              uint              `json:"user_id"`
	Ban             bool              `json:"lock"`
	SystemSuper     bool              `json:"id3"`
	EntitySuper     bool              `json:"id2"`
	DepartmentSuper bool              `json:"id1"`
	Employee        bool              `json:"id0" default:"true"`
	Department      *model.Department `json:"department"`
}

type DepartmentUserListResponse struct {
	UserList []DepartmentUserInfo `json:"user_list"`
}

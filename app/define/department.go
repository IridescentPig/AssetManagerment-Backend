package define

type DepartmentBasicInfo struct {
	ID       uint   `json:"department_id"`
	Name     string `json:"department_name"`
	ParentID *uint  `json:"parent_id"`
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

type DepartmentUserListResponse struct {
	UserList []UserInfo `json:"user_list"`
}

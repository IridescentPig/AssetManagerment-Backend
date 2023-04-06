package define

type CreateEntityReq struct {
	Name string `json:"name"`
}

type EntityInfoResponse struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Departments []DepartmentBasicInfo `json:"departments"`
}

package define

type DepartmentBasicInfo struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

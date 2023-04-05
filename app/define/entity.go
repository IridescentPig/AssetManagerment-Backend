package define

type EntityInfoResponse struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Departments []DepartmentBasicInfo `json:"departments"`
}

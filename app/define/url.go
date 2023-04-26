package define

type CreateOrModifyUrlReq struct {
	Name            string `json:"name"`
	Url             string `json:"url"`
	DepartmentSuper bool   `json:"department_super"`
	EntitySuper     bool   `json:"entity_super"`
	SystemSuper     bool   `json:"system_super"`
}

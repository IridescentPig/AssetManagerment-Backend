package define

type CreateUrlReq struct {
	Name            string `json:"name"`
	Url             string `json:"url"`
	DepartmentSuper bool   `json:"department_super"`
	EntitySuper     bool   `json:"entity_super"`
	SystemSuper     bool   `json:"system_super"`
}

type ModifyUrlReq struct {
	OldName         string `json:"old_name"`
	Name            string `json:"name"`
	Url             string `json:"url"`
	DepartmentSuper bool   `json:"department_super"`
	EntitySuper     bool   `json:"entity_super"`
	SystemSuper     bool   `json:"system_super"`
}

type DeleteUrlReq struct {
	Name string `json:"name"`
}

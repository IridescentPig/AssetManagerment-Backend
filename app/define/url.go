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

type GetUrlResponse struct {
	UrlList []UrlInfo `json:"url_list"`
}

type UrlInfo struct {
	Name            string `gorm:"column:name;not null" json:"name"`
	Url             string `gorm:"column:url" json:"url"`
	DepartmentSuper bool   `gorm:"column:department_super;default:false" json:"department_super"`
	EntitySuper     bool   `gorm:"column:entity_super;default:false" json:"entity_super"`
	SystemSuper     bool   `gorm:"column:system_super;default:false" json:"system_super"`
}

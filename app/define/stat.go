package define

import "asset-management/app/model"

type DepartmentStatTotalResponse struct {
	Stats []*model.Stat `json:"stats"`
}

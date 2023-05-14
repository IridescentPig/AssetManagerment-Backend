package define

import (
	"asset-management/app/model"

	"github.com/shopspring/decimal"
)

type DepartmentStatTotalResponse struct {
	Stats []*model.Stat `json:"stats"`
}

type AssetDistribution struct {
	State uint            `gorm:"column:state" json:"state"`
	Count uint            `gorm:"column:count" json:"count"`
	Total decimal.Decimal `gorm:"total" json:"total_worth"`
}

type AssetDistributionResponse struct {
	Distribution []*AssetDistribution `json:"distribution"`
}

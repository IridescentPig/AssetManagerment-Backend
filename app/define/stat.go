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
	Total decimal.Decimal `gorm:"column:total" json:"total_worth"`
}

type AssetDistributionResponse struct {
	Distribution []*AssetDistribution `json:"distribution"`
}

type DepartmentAssetDistribution struct {
	DepartmentID   uint            `gorm:"column:department_id" json:"department_id"`
	Count          uint            `gorm:"column:count" json:"count"`
	Total          decimal.Decimal `gorm:"column:total" json:"total_worth"`
	DepartmentName string          `gorm:"column:department_name" json:"department_name"`
}

type DepartmentAssetDistributionResponse struct {
	Stats []*DepartmentAssetDistribution `json:"stats"`
}

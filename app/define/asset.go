package define

import (
	"asset-management/app/model"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type AssetInfo struct {
	AssetID     uint             `json:"asset_id" copier:"ID"`
	AssetName   string           `json:"asset_name" copier:"Name"`
	ParentID    uint             `json:"parent_id"`
	User        model.User       `json:"user"`
	Department  model.Department `json:"department"`
	Maintainer  model.User       `json:"maintainer"`
	Price       decimal.Decimal  `json:"price"`
	Description string           `json:"description"`
	Position    string           `json:"position"`
	Expire      uint             `json:"expire"`
	Class       model.AssetClass `json:"asset_class"`
	Number      int              `json:"count"`
	Type        int              `json:"type"`
	Children    []*AssetInfo     `json:"children"`
	State       uint             `json:"state"`
	Property    datatypes.JSON   `json:"property"`
	NetWorth    decimal.Decimal  `json:"net_worth"`
	CreatedAt   *model.ModelTime `json:"created_at"`
}

type ModifyAssetInfoReq struct {
	AssetName   string          `json:"asset_name"`
	ParentID    *uint           `json:"parent_id"`
	Price       decimal.Decimal `json:"price"`
	Description string          `json:"description"`
	Position    string          `json:"position"`
	ClassID     uint            `json:"class_id"`
	Type        int             `json:"type"`
	Number      int             `json:"count"`
	Expire      uint            `json:"expire" binding:"gte=0"`
}

type CreateAssetReq struct {
	AssetName   string            `json:"asset_name"`
	Price       decimal.Decimal   `json:"price"`
	Description string            `json:"description"`
	Position    string            `json:"position"`
	ClassID     uint              `json:"class_id"`
	Number      int               `json:"count"`
	Type        int               `json:"type"`
	ParentID    uint              `json:"parent_id"`
	Children    []*CreateAssetReq `json:"children"`
}

type CreateAssetListReq struct {
	AssetList []CreateAssetReq `json:"asset_list"`
}

type ExpireAssetReq struct {
	AssetID uint `json:"asset_id" copier:"ID"`
}

type ExpireAssetListReq struct {
	ExpireList []ExpireAssetReq `json:"asset_list"`
}

type AssetListResponse struct {
	AssetList []*AssetInfo `json:"asset_list"`
}

// 暂时借用 Expire 的请求体结构
type AssetTransferReq struct {
	UserID uint             `json:"user_id"`
	Assets []ExpireAssetReq `json:"assets"`
}

type AssetPropertyReq struct {
	Key   string `json:"key" bind:"required,gt=0"`
	Value string `json:"value"`
}

type DeleteAssetPropertyReq struct {
	Key string `json:"key" bind:"required,gt=0"`
}

type AssetHistory struct {
	Type               uint             `json:"type"`
	ReviewTime         *model.ModelTime `json:"time"`
	UserID             uint             `json:"user_id"`
	Username           string           `json:"username"`
	DepartmentID       uint             `json:"department_id"`
	TargetID           uint             `json:"target_user_id"`
	TargetName         string           `json:"target_username"`
	TargetDepartmentID uint             `json:"target_department_id"`
}

type AssetHistoryResponse struct {
	History []*AssetHistory `json:"history"`
}

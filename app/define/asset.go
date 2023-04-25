package define

import (
	"asset-management/app/model"

	"github.com/shopspring/decimal"
)

type AssetInfo struct {
	AssetID     uint             `json:"asset_id" copier:"ID"`
	AssetName   string           `json:"asset_name" copier:"Name"`
	ParentID    uint             `json:"parent_id"`
	User        model.User       `json:"user"`
	Department  model.Department `json:"department"`
	Price       decimal.Decimal  `json:"price"`
	Description string           `json:"description"`
	Position    string           `json:"position"`
	Expire      bool             `json:"expire"`
	Class       model.AssetClass `json:"asset_class"`
	Number      int              `json:"count"`
	Type        int              `json:"type"`
	Children    []*AssetInfo     `json:"children"`
	State       uint             `json:"state"`
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
}

type CreateAssetReq struct {
	AssetName   string            `json:"asset_name"`
	Price       decimal.Decimal   `json:"price"`
	Description string            `json:"description"`
	Position    string            `json:"position"`
	ClassID     uint              `json:"class_id"`
	Number      int               `json:"count"`
	Type        int               `json:"type"`
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

package define

type CreateAssetClassReq struct {
	ClassName string `json:"class_name" binding:"required"`
	ParentID  uint   `json:"parent_id" binding:"required"`
	Type      int    `json:"type" binding:"required"`
}

type ModifyAssetClassReq struct {
	ClassName string `json:"class_name"`
	ParentID  uint   `json:"parent_id"`
	Type      int    `json:"type"`
}

type AssetClassTreeNode struct {
	ClassID   uint                  `json:"class_id" copier:"ID"`
	ClassName string                `json:"class_name" copier:"Name"`
	ParentID  uint                  `json:"parent_id"`
	Type      int                   `json:"type"`
	Children  []*AssetClassTreeNode `json:"children"`
}

type AssetClassTreeResponse struct {
	AssetClassTree []*AssetClassTreeNode `json:"asset_class_tree"`
}

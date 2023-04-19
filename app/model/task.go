package model

type Task struct {
	ID              uint   `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"task_id"`
	TaskType        uint   `gorm:"column:task_type" json:"task_type"` // 0领用、1退库、2维保、3转移
	TaskDescription string `gorm:"column:task_description" json:"task_description"`
	UserID          uint   `gorm:"default:null;column:user_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_id"`
	UserName        string `gorm:"column:user_name" json:"user_name"`
	User            User   `gorm:"foreignKey:UserID;references:ID;default:null" json:"user"`
	// Task的回应中还带有UserName，这里不设额外的外键了，用User.UserName在回应时设置，不设约束
	TargetID   uint     `gorm:"default:null;column:target_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"target_id"`
	TargetName string   `gorm:"column:target_name" json:"target_name"`
	Target     User     `gorm:"foreignKey:TargetID;references:ID;default:null" json:"target"`
	AssetList  []*Asset `json:"asset_list"`
}

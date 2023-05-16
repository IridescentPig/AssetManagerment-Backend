package model

type Task struct {
	ID              uint       `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"task_id"`
	TaskType        uint       `gorm:"column:task_type" json:"task_type"` // 0领用、1退库、2维保、3转移
	TaskDescription string     `gorm:"column:task_description" json:"task_description"`
	UserID          uint       `gorm:"default:null;column:user_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_id"`
	User            User       `gorm:"foreignKey:UserID;references:ID;default:null" json:"user"`
	TargetID        uint       `gorm:"default:null;column:target_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"target_id"`
	Target          User       `gorm:"foreignKey:TargetID;references:ID;default:null" json:"target"`
	DepartmentID    uint       `gorm:"default:null;column:department_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"department_id"`
	Department      Department `gorm:"foreignKey:DepartmentID;references:ID;default:null" json:"department"`
	AssetList       []*Asset   `gorm:"many2many:task_assets;" json:"asset_list"`
	State           uint       `gorm:"default:0;colunm:state" json:"state"` // 0提交未审批、1批准、2不通过、3自行撤销
	CreatedAt       *ModelTime `gorm:"column:created_at" json:"created_at"`
	ReviewAt        *ModelTime `gorm:"column:review_at" json:"review_at"`
}

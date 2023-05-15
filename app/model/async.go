package model

type AsyncTask struct {
	ID           uint       `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Type         uint       `gorm:"type" json:"type"` // 0-import 1-export login-log 2-export modify-log
	UserID       uint       `gorm:"default:null;column:user_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_id"`
	User         User       `gorm:"foreignKey:UserID;references:ID;default:null" json:"user"`
	DepartmentID uint       `gorm:"default:null;column:department_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"department_id"`
	Department   Department `gorm:"foreignKey:DepartmentID;references:ID;default:null" json:"department"`
	EntityID     uint       `gorm:"default:null;column:entity_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entity_id"`
	Entity       Entity     `gorm:"foreignKey:EntityID;references:ID;default:null" json:"-"`
	DownloadLink string     `gorm:"download_link" json:"download_link"`
	State        uint       `gorm:"state" json:"state"` // 0-Pending 1-Running 2-Success 3-Fail
	Message      string     `gorm:"message" json:"message"`
	FromTime     *ModelTime `gorm:"from_time" json:"from_time"` // use for export log
}

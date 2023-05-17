package model

type AsyncTask struct {
	ID           uint       `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Type         uint       `gorm:"column:type" json:"type"` // 0-import 1-export login-log 2-export modify-log
	UserID       uint       `gorm:"default:null;column:user_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_id"`
	User         User       `gorm:"foreignKey:UserID;references:ID;default:null" json:"user"`
	DepartmentID uint       `gorm:"default:null;column:department_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"department_id"`
	Department   Department `gorm:"foreignKey:DepartmentID;references:ID;default:null" json:"department"`
	EntityID     uint       `gorm:"default:null;column:entity_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entity_id"`
	Entity       Entity     `gorm:"foreignKey:EntityID;references:ID;default:null" json:"-"`
	ObjectKey    string     `gorm:"column:object_key" json:"object_key"`
	DownloadLink string     `gorm:"column:download_link" json:"download_link"`
	State        uint       `gorm:"column:state" json:"state"` // 0-Pending 1-Running 2-Success 3-Fail 4-cancel
	Message      string     `gorm:"column:message" json:"message"`
	FromTime     *ModelTime `gorm:"column:from_time" json:"from_time"` // use for export log
	LogType      uint       `gorm:"column:log_type" json:"log_type"`   // 0-all logs 1-success 2-failed
}

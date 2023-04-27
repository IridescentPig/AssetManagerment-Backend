package model

type Log struct {
	ID           uint       `gorm:"primaryKey;column:id;AUTO_INCREMENT"`
	Method       string     `gorm:"column:method" json:"method"`
	URL          string     `gorm:"column:url" json:"url"`
	Status       int        `gorm:"column:status" json:"status"`
	ErrorCode    int        `gorm:"column:error_code" json:"error_code"`
	ErrorMessage string     `gorm:"default:None;column:error_message" json:"error_message"`
	UserID       uint       `gorm:"column:user_id" json:"user_id"`
	Username     string     `gorm:"column:username" json:"username"`
	EntityID     uint       `gorm:"column:entity_id" json:"entity_id"`
	DepartmentID uint       `gorm:"column:department_id" json:"department_id"`
	Time         *ModelTime `gorm:"column:time" json:"time"`
	Level        string     `gorm:"column:level" json:"level"`
	Message      string     `gorm:"column:message" json:"message"`
}

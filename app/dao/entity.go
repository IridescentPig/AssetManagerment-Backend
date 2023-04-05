package dao

type entityDao struct {
}

var EntityDao *entityDao

func newEntityDao() *entityDao {
	return &entityDao{}
}

func init() {
	EntityDao = newEntityDao()
}

// func (entity *entityDao) Create(entityname string) error {
// 	newEntity := model.Entity{Name: entityname}
// 	result := db.Model(&model.Entity{}).Create(&newEntity)
// 	return utils.DBError(result)
// }

// func (entity *entityDao) Update(id uint, data map[string]interface{}) error {
// 	result := db.Model(&model.Entity{}).Where("id = ?", id).Updates(data)
// 	return utils.DBError(result)
// }

// func (entity *entityDao) Delete(id []uint) error {
// 	result := db.Model(&model.Entity{}).Where("id in (?)", id).Delete(&model.Entity{})
// 	return utils.DBError(result)
// }

// func (entity *entityDao) AllEntity() (list []model.Entity, err error) {
// 	result := db.Model(&model.Entity{}).Find(&list)
// 	err = utils.DBError(result)
// 	return
// }

// func (entity *entityDao) AllEntityWhere(query interface{}, args ...interface{}) (list []model.Entity, err error) {
// 	result := db.Model(&model.Entity{}).Where(query, args...).Find(&list)
// 	err = utils.DBError(result)
// 	return
// }

// func (entity *entityDao) OneEntityWhere(query interface{}, args ...interface{}) (record model.Entity, err error) {
// 	result := db.Model(&model.Entity{}).Where(query, args...).First(&record)
// 	err = utils.DBError(result)
// 	return
// }

// func (entity *entityDao) EntityCount() (count int64, err error) {
// 	result := db.Model(&model.Entity{}).Count(&count)
// 	err = utils.DBError(result)
// 	return
// }

// func (entity *entityDao) EntityCountWhere(query interface{}, args ...interface{}) (count int64, err error) {
// 	result := db.Model(&model.Entity{}).Where(query, args...).Count(&count)
// 	err = utils.DBError(result)
// 	return
// }

package dao

import (
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/utils"
)

type statDao struct {
}

var StatDao *statDao

func newStatDao() *statDao {
	return &statDao{}
}

func init() {
	StatDao = newStatDao()
}

func (stat *statDao) GetAllAssetStat() ([]*model.Stat, error) {
	var stats []*model.Stat

	result := db.Model(model.Asset{}).Select("department_id, SUM(net_worth) as total").Group("department_id").Scan(&stats)
	if result.Error != nil {
		return nil, utils.DBError(result)
	}

	return stats, nil
}

func (stat *statDao) CreateAssetStats(stats []*model.Stat) error {
	result := db.Model(&model.Stat{}).Create(stats)

	return utils.DBError(result)
}

func (stat *statDao) GetDepartmentStat(departmentID uint) (stats []*model.Stat, err error) {
	result := db.Model(&model.Stat{}).Where("department_id = ?", departmentID).Find(&stats)

	err = utils.DBError(result)
	return
}

func (stat *statDao) GetDepartmentAssetDistribution(departmentID uint) (distribution []*define.AssetDistribution, err error) {
	result := db.Model(&model.Asset{}).Select("state, COUNT(*) as count, SUM(net_worth) as total").Group("state").Scan(&distribution)

	err = utils.DBError(result)
	return
}

func (stat *statDao) GetDepartmentsAssetDistribution(IDs []uint, distribution []*define.DepartmentAssetDistribution) (err error) {
	result := db.Model(&model.Asset{}).Where("department_id IN (?)", IDs).Select("department_id, departments.name as department_name, COUNT(*) as count, SUM(net_worth) as total").Joins("left join departments on assets.department_id = departments.id").Group("department_id").Scan(&distribution)
	err = utils.DBError(result)
	return
}

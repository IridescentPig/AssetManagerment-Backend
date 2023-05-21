package timing

import (
	"asset-management/app/dao"
	"asset-management/app/model"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

type AssetDepreciate struct {
}

func (depreciate *AssetDepreciate) Run() {
	assetList, err := dao.AssetDao.GetAllAssets()

	if err == nil {
		for _, asset := range assetList {
			if asset.State < 3 && asset.Expire != 0 {
				// log.Println(asset)
				interval := getDiffDays(time.Time(*asset.CreatedAt), time.Now())
				if interval >= int(asset.Expire) {
					err = dao.AssetDao.Update(asset.ID, map[string]interface{}{
						"net_worth": decimal.Zero,
						"state":     3,
						"parent_id": gorm.Expr("NULL"),
					})

					if err != nil {
						continue
					}

					subAssets, _, err := dao.AssetDao.GetSubAsset(asset.ID, -1, -1)
					if err == nil {
						subAssetIDs := funk.Map(subAssets, func(thisAsset *model.Asset) uint {
							return thisAsset.ID
						}).([]uint)

						err := dao.AssetDao.AllUpdate(subAssetIDs, map[string]interface{}{
							"parent_id": gorm.Expr("NULL"),
						})

						if err != nil {
							continue
						}
					}
				} else {
					rate := 1.0 - float64(interval)/float64(asset.Expire)
					asset.NetWorth = asset.Price.Mul(decimal.NewFromFloat(rate))
					asset.Warn = (int(asset.Expire) - interval) < int(asset.Threshold)

					err = dao.AssetDao.Update(asset.ID, map[string]interface{}{
						"net_worth": asset.NetWorth,
						"warn":      asset.Warn,
					})

					if err != nil {
						continue
					}
				}
			}
		}
		log.Println("AssetDepreciate Succeed")
	} else {
		log.Println("AssetDepreciate Failed")
	}
}

type AssetStat struct {
}

func (stat *AssetStat) Run() {
	stats, err := dao.StatDao.GetAllAssetStat()
	if err != nil {
		log.Println("AssetStat Failed")
		return
	}

	now := time.Now()

	for _, stat := range stats {
		stat.Time = model.ModelTime(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, timezone))
	}

	err = dao.StatDao.CreateAssetStats(stats)
	if err != nil {
		log.Println("AssetStat Failed")
		return
	}

	log.Println("AssetStat Succeed")
}

func getDiffDays(t1, t2 time.Time) int {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	timeDay1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, timezone)
	timeDay2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, timezone)

	return int(timeDay2.Sub(timeDay1).Hours() / 24)
}

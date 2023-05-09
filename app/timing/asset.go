package timing

import (
	"asset-management/app/dao"
	"asset-management/app/model"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

var timezone *time.Location

/*
This package is for timing task: asset depreciate and statistics
*/
func Init() *cron.Cron {
	timezone, _ = time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(timezone))

	// _, _ = c.AddFunc("@every 1s", func() {
	// 	log.Println("Hello")
	// })

	_, err := c.AddJob(
		"0 3 * * *",
		cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(&AssetDepreciate{}),
	)

	if err != nil {
		log.Println("Something error when register daily job")
	}

	_, err = c.AddJob(
		"0 4 * * *",
		cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(&AssetStat{}),
	)

	if err != nil {
		log.Println("Something error when register daily job")
	}

	return c
}

type AssetDepreciate struct {
}

func (depreciate *AssetDepreciate) Run() {
	assetList, err := dao.AssetDao.GetAllAssets()

	if err == nil {
		for _, asset := range assetList {
			_ = dao.AssetDao.SaveAsset(asset)
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
		*stat.Time = model.ModelTime(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, timezone))
	}

	err = dao.StatDao.CreateAssetStats(stats)
	if err != nil {
		log.Println("AssetStat Failed")
		return
	}

	log.Println("AssetStat Succeed")
}

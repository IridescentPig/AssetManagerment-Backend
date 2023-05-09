package timing

import (
	"asset-management/app/dao"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

/*
This package is for timing task: asset depreciate and statistics
*/
func Init() *cron.Cron {
	c := cron.New(cron.WithLocation(time.Local))

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
	}
}

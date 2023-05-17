package timing

import (
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

	_, err = c.AddJob(
		"@every 30s",
		cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)).Then(&GetPendingAsyncTask{}),
	)

	if err != nil {
		log.Println("Something error when async task job")
	}

	return c
}

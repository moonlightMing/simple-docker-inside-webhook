package cronjob

import (
	"github.com/moonlightming/simple-docker-inside-webhook/conf"
	"github.com/robfig/cron"
)

var (
	config        conf.Config
	cronScheduler *cron.Cron
)

func init() {
	config = conf.NewConfig()
	cronScheduler = cron.New()

	for _, job := range config.CronEvents {
		switch job.Event {
		case CleanNoneTagImage:
			if err := cronScheduler.AddJob(job.Spec, new(CleanNoneTagImageJob)); err != nil {
				panic(err)
			}
		}
	}
	cronScheduler.Start()
	select {}
}

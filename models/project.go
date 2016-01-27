package models

import (
	"time"

	"github.com/chronitonapp/chroniton/utils"
)

type Project struct {
	Id                  int64
	UserId              int64
	Name                string `form:"name"`
	GitIntegrationName  string `form:"gitIntegrationName"`
	GitRepoName         string `form:"gitRepoName"`
	PmIntegrationName   string `form:"pmIntegrationName"`
	NumRecievedWebhooks int
	CreatedAt           time.Time
}

func (p Project) TotalHours() float64 {
	var total float64
	row := utils.ORM.Table("time_trackeds").Select("round(sum(duration) / 3600, 1)").
		Where("project_id = ?", p.Id).Row()
	row.Scan(&total)
	// if err != nil {
	// 	utils.Log.Error("Failed to calc total time tracked hours: %v", err)
	// 	return 0
	// }

	return total
}

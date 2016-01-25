package models

import (
	"time"
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

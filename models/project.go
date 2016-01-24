package models

import (
	"time"
)

type Project struct {
	Id                 int64
	UserId             int64
	Name               string
	GitIntegrationName string
	GitRepoName        string
	PmIntegrationName  string
	CreatedAt          time.Time
}

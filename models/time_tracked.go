package models

import (
	"time"
)

type TimeTracked struct {
	Id        int64
	UserId    int64
	ProjectId int64
	CommitSHA string
	Duration  int64
	CreatedAt time.Time
}

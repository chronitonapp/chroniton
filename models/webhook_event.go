package models

import (
	"time"
)

type WebhookEvent struct {
	Id        int64
	ProjectId int64
	UserId    int64
	Payload   string `sql:"type:text"`
	CreatedAt time.Time
}

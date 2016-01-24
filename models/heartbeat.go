package models

import (
	"time"

	"github.com/gophergala2016/chroniton/utils"

	wakatime "github.com/jsimnz/go.wakatime"
)

type Heartbeat struct {
	Id          int64  // Database ID
	Project     string // Associated project on Wakatime (git project)
	Branch      string // Git branch
	Entity      string
	Wid         string // Wakatime heartbeat ID
	IsDebugging bool   // If the heartbeat was sent during a debug
	IsWrite     bool   // If the heartbeat was sent during a write
	Language    string // Language of the file being edited during the heartbeat
	Type        string
	Time        time.Time // Time of the heartbeat
}

func NewHeartbeatFromWakaTime(heartbeat wakatime.Heartbeat) Heartbeat {
	return Heartbeat{
		Project:     heartbeat.Project,
		Branch:      heartbeat.Branch,
		Entity:      heartbeat.Entity,
		Wid:         heartbeat.ID,
		IsDebugging: heartbeat.IsDebugging,
		IsWrite:     heartbeat.IsWrite,
		Language:    heartbeat.Language,
		Type:        heartbeat.Type,
	}
}

func SaveHeartbeatIfNotExist(heartbeat Heartbeat) bool {
	count := 0
	utils.ORM.Model(Heartbeat{}).Where("wid = ?", heartbeat.Wid).Count(&count)
	if count == 0 {
		utils.ORM.Save(heartbeat)
		return true
	}
	return false
}

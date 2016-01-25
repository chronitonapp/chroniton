package models

import (
	"time"

	"github.com/gophergala2016/chroniton/utils"

	"github.com/dustin/go-humanize"
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
		Time:        time.Time(*heartbeat.Time),
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

func (h Heartbeat) HumanTime() string {
	return humanize.Time(h.Time)
}

func CalcTotalHeartbeatsDuration(user User, heartbeats []Heartbeat) int64 {
	duration := int64(0)
	var lastHeartbeat Heartbeat
	var curHeartbeat Heartbeat

	lastHeartbeat = heartbeats[0]
	for i := 1; i <= len(heartbeats)-1; i++ {
		curHeartbeat = heartbeats[i]
		t := curHeartbeat.Time.Unix() - lastHeartbeat.Time.Unix()
		lastHeartbeat = curHeartbeat
		if t >= 60*15 { /* 15 minutes */
			continue
		}
		duration += t
	}

	return duration
}

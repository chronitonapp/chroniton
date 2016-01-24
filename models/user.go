package models

import (
	"math"
	"time"

	"github.com/gophergala2016/chroniton/utils"

	wakatime "github.com/jsimnz/go.wakatime"
)

type User struct {
	Id                int64  `form:"id"`
	Name              string `form:"name"`
	Email             string `form:"email"`
	Password          string `form:"password"`
	WakaTimeApiKey    string `form:"wakaTimeApiKey"`
	IsSyncingWakaTime bool
	LastWakaTimeSync  time.Time
	CreatedAt         time.Time
}

func (u User) Verify() []error {
	return make([]error, 0)
}

func (u User) PullNewestHeartbeats() {
	utils.Log.Debug("Pulling new heartbeats")
	wt := wakatime.NewWakaTime(u.WakaTimeApiKey)
	wtUser, err := wt.GetUser("")
	if err != nil {
		utils.Log.Critical("Failed to get wakatime user info: ", err)
		return
	}
	utils.Log.Debug("LastWakaTimeSync: %v", u.LastWakaTimeSync)
	if u.LastWakaTimeSync.Unix() < 0 {
		u.LastWakaTimeSync = *wtUser.CreatedAt
	}

	u.IsSyncingWakaTime = true
	utils.ORM.Save(&u)

	current := wtUser.LastHeartbeat
	daysToSync := int(math.Min(7, math.Ceil(current.Sub(u.LastWakaTimeSync).Hours()/float64(24))))

	curDate := time.Now().Add(-1 * time.Hour * 24 * time.Duration(daysToSync))
	utils.Log.Debug("Days of heartbeats to sync: %v", daysToSync)
	for i := 0; i < daysToSync; i++ {
		curDate = curDate.Add(time.Hour * 24) // add 1 day
		utils.Log.Debug("Getting heartbeats for: %v", curDate)
		heartbeats, err := wt.GetHeartbeats(&wakatime.HeartbeatParameters{
			Date: &curDate,
			Show: []string{"project", "branch", "entity", "id", "language", "type", "time"},
			User: "",
		})
		if err != nil {
			utils.Log.Error("failed to get heartbeats for user %v on date %v. Error: %v", u.Id, curDate, err)
			continue
		}
		for _, heartbeat := range heartbeats {
			utils.Log.Debug("heartbeat from wakatime: %v", heartbeat)
			hb := NewHeartbeatFromWakaTime(heartbeat)
			utils.Log.Debug("heartbeat formatted: %v", hb)
			SaveHeartbeatIfNotExist(hb)
		}
	}
	utils.Log.Debug("LastHearbeat: %v", *wtUser.LastHeartbeat)
	u.LastWakaTimeSync = *wtUser.LastHeartbeat
	u.IsSyncingWakaTime = false
	utils.ORM.Save(&u)
}

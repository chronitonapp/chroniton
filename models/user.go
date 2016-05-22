package models

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/chronitonapp/chroniton/utils"

	"github.com/chronitonapp/gormseries"
	"github.com/jinzhu/gorm"
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

	// Associated Models
	Projects []Project
}

func (u User) Verify() []error {
	return make([]error, 0)
}

func (u *User) AfterFind(db *gorm.DB) error {
	utils.ORM.Model(u).Related(&u.Projects)
	return nil
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
			hb.UserId = u.Id
			utils.Log.Debug("heartbeat formatted: %v", hb)
			SaveHeartbeatIfNotExist(hb)
		}
	}
	utils.Log.Debug("LastHearbeat: %v", *wtUser.LastHeartbeat)
	u.LastWakaTimeSync = *wtUser.LastHeartbeat
	u.IsSyncingWakaTime = false
	utils.ORM.Save(&u)
}

func (u User) NumReceivedPushes() int {
	count := 0
	row := utils.ORM.Table("projects").Where("user_id = ?", u.Id).Select("sum(num_recieved_webhooks)").Row()
	row.Scan(&count)
	return count
}

func (u User) Heartbeats() []Heartbeat {
	var heartbeats []Heartbeat
	utils.ORM.Table("heartbeats").Joins("JOIN projects ON projects.name = heartbeats.project JOIN users ON projects.user_id = users.id").
		Order("heartbeats.time DESC").Find(&heartbeats)
	return heartbeats
}

func (u User) TotalMonthTimeTacked() string {
	var total []float64

	err := utils.ORM.Table("time_trackeds").
		Where("(date_trunc('day', created_at) <= date_trunc('month', now()) + INTERVAL '1 Month - 1 Day')").
		Where("date_trunc('month', now()) <= date_trunc('day', created_at)").
		Pluck("round(sum(duration) / 3600,1)", &total).Error

	if err != nil {
		utils.Log.Error("failed to get time tracked month total: %v", err)
		return "0"
	}

	return utils.RoundFloat(float64(total[0])) // convert to hours
}

func (u User) TotalYearTimeTacked() string {
	var total []float64
	err := utils.ORM.Table("time_trackeds").
		Where("(date_trunc('day', created_at) <= date_trunc('year', now()) + INTERVAL '1 year - 1 Day')").
		Where("date_trunc('year', now()) <= date_trunc('day', created_at)").
		Pluck("round(sum(duration) / 3600,1)", &total).Error
	if err != nil {
		utils.Log.Error("failed to get time tracked month total: %v", err)
		return "0"
	}

	return utils.RoundFloat(total[0]) // convert to hours
}

func (u User) TimeTrackedChartData() string {
	var count []float64
	err := utils.ORM.TimeSeries(gormseries.Last7Days).Table("time_trackeds").Group("day").
		Pluck("round(sum(coalesce(duration,0)) / 3600, 1)", &count).Error

	if err != nil {
		utils.Log.Error("Failed to get time tracked chart data! %v", err)
		return "[]"
	}

	lst := fmt.Sprint(count)
	return strings.Replace(lst, " ", ",", -1)
}

func (u User) NumHeartbeatsChartData() string {
	var count []float64
	err := utils.ORM.TimeSeries(gormseries.Last7Days, "day = time").Table("heartbeats").Group("day").
		Pluck("count(*)-1", &count).Error

	if err != nil {
		utils.Log.Error("Failed to get time tracked chart data! %v", err)
		return "[]"
	}

	lst := fmt.Sprint(count)
	return strings.Replace(lst, " ", ",", -1)
}

func (u User) TopThreeLanguages() [][]string {

	type langResult struct {
		Language string
		Count    int
	}

	var tmpResults []langResult
	results := make([][]string, 0)
	// err := utils.ORM.Table("heartbeats").Select([]string{"laguage"}).
	// 	Group("language").Error

	err := utils.ORM.Raw(`select language, count(*) from heartbeats where user_id = ? GROUP BY language ORDER BY count DESC`, u.Id).Scan(&tmpResults).Error
	if err != nil {
		utils.Log.Error("Failed to get top three languages for user: %v", err)
		return results
	}

	sum := 0
	for i := 0; i < len(tmpResults); i++ {
		sum += tmpResults[i].Count
	}
	for _, topLangRes := range tmpResults {
		r := []string{
			topLangRes.Language,
			strconv.Itoa(
				int(math.Ceil((float64(topLangRes.Count) / float64(sum)) * 100))),
			strconv.Itoa(topLangRes.Count),
		}
		results = append(results, r)
	}

	return results[0:3]
}

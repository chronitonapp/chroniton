package models

import (
	"github.com/chronitonapp/chroniton/utils"
)

func init() {
	utils.ORM.AutoMigrate(&User{})
	utils.ORM.AutoMigrate(&Heartbeat{})
	utils.ORM.AutoMigrate(&Project{})
	utils.ORM.AutoMigrate(&WebhookEvent{})
	utils.ORM.AutoMigrate(&TimeTracked{})
}

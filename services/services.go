package services

import (
	"github.com/go-martini/martini"
)

type Service interface {
	Register(martini.Router)
}

var Services []Service = []Service{
	new(UserService),
	new(WebhookService),
	new(ProjectService),
}

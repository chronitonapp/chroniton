package services

import (
	"github.com/go-martini/martini"

	"github.com/chronitonapp/chroniton/models"
)

type Service interface {
	Register(martini.Router)
}

var Services []Service = []Service{
	new(UserService),
	new(WebhookService),
	new(ProjectService),
}

type ResponseCtx struct {
	Error       bool
	ErrMessage  string
	Vars        map[string]interface{}
	CurrentUser models.User
	models.User
}

func (rctx ResponseCtx) craft(vars []map[string]interface{}) ResponseCtx {
	if len(vars) > 0 {
		rctx.Vars = vars[0]
	}

	return rctx
}

func ErrorResponse(errMsg string, currentUser models.User, vars ...map[string]interface{}) ResponseCtx {
	rctx := ResponseCtx{
		Error:       true,
		ErrMessage:  errMsg,
		CurrentUser: currentUser,
		User:        currentUser,
	}

	return rctx.craft(vars)
}

func SuccessResponse(currentUser models.User, vars ...map[string]interface{}) ResponseCtx {
	rctx := ResponseCtx{
		Error:       false,
		CurrentUser: currentUser,
		User:        currentUser,
	}
	return rctx.craft(vars)
}

package services

import (
	"io/ioutil"
	"net/http"
	//"strings"
	//"time"

	// "github.com/gophergala2016/chroniton/models"
	"github.com/gophergala2016/chroniton/utils"

	"github.com/go-martini/martini"
	// "github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	// //"github.com/martini-contrib/sessions"
)

type WebhookService struct{}

func (ws WebhookService) Register(router martini.Router) {
	router.Post("/webhook/:git_service", ws.Handle)
}

func (ws WebhookService) Handle(r render.Render, req *http.Request, params martini.Params) {
	gitServiceStr := params["git_service"]
	gitService, exists := registeredGitServices[gitServiceStr]
	if !exists {
		utils.Log.Warning("No git integration found for recieved webhook %v", gitServiceStr)
		r.Status(200)
		return
	}

	defer req.Body.Close()
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.Log.Error("Failed to ready the webhook payload: %v", err)
		r.Status(200)
		return
	}

	go gitService.HandleWebhookEvent(payload)
	r.Status(200)
}

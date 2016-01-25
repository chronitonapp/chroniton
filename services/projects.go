package services

import (
	// "net/http"
	// "strings"
	"fmt"
	//"time"

	"github.com/gophergala2016/chroniton/models"
	"github.com/gophergala2016/chroniton/utils"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	//"github.com/martini-contrib/sessions"
	//"github.com/fatih/structs"
)

type ProjectService struct{}

func (ps ProjectService) Register(router martini.Router) {
	router.Group("/projects", func(rtr martini.Router) {
		rtr.Get("/new", ps.New)
		rtr.Post("", binding.Bind(models.Project{}), ps.Create)
	}, EnsureAuth)
}

func (ps ProjectService) New(r render.Render) {
	r.HTML(200, "project/new", nil)
}

func (ps ProjectService) Create(currentUser models.User, project models.Project, r render.Render) {
	project.UserId = currentUser.Id
	err := utils.ORM.Save(&project).Error
	if err != nil {
		utils.Log.Error("Failed to create project: %v", err)
		r.HTML(403, "project/new", project)
	}
	r.Redirect(fmt.Sprintf("projects/%v", project.Id))
}

package main

import (
	_ "github.com/gophergala2016/chroniton/models"
	"github.com/gophergala2016/chroniton/services"
	"github.com/gophergala2016/chroniton/utils"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

func main() {
	m := martini.Classic()

	m.Use(martini.Static("assets", martini.StaticOptions{
		Prefix: "/assets",
	}))
	m.Use(martini.Static("public"))
	m.Use(render.Renderer(render.Options{
		// Funcs:      []template.FuncMap{utils.TemplateHelpers},
		Layout:     "layout",
		Extensions: []string{".html"},
	}))

	m.Use(sessions.Sessions("chroniton", utils.SessionStore))

	for _, service := range services.Services {
		service.Register(m.Router)
	}

	utils.Log.Info("Running Chroniton Server...")
	m.Run()
}

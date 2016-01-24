package services

import (
	"net/http"

	"github.com/gophergala2016/chroniton/models"
	"github.com/gophergala2016/chroniton/utils"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

func PreventReauth(session sessions.Session, r render.Render) {
	_, ok := session.Get("id").(int64)
	if ok {
		session.AddFlash("warning: You are already signed in!")
		r.Redirect("/dashboard")
	}
}

func EnsureAuth(session sessions.Session, r render.Render, req *http.Request, c martini.Context) {
	id, ok := session.Get("id").(int64)
	if !ok || id == 0 {
		session.AddFlash("warning: You must login first!")
		session.Set("previous_url", req.RequestURI)
		r.Redirect("/signin")
	} else if ok {
		var user models.User
		err := utils.ORM.First(&user, id).Error
		if err != nil {
			r.Error(500)
			return
		}
		c.Map(user)
	}
}

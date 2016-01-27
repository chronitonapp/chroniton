package services

import (
	"net/http"
	"strings"
	//"time"

	"github.com/chronitonapp/chroniton/models"
	"github.com/chronitonapp/chroniton/utils"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	//"github.com/martini-contrib/sessions"
	"github.com/fatih/structs"
)

type UserService struct{}

func (us UserService) Register(router martini.Router) {
	utils.Log.Debug("Registering UserService!")

	router.Group("/users", func(rtr martini.Router) {
		rtr.Post("", PreventReauth, binding.Bind(models.User{}), us.Create)
		rtr.Post("/update", EnsureAuth, binding.Bind(models.User{}), us.Update)
	})

	router.Get("/signout", us.SignOut)
	router.Get("/signup", PreventReauth, us.SignUp)
	router.Get("/signin", PreventReauth, us.SignIn)
	router.Get("/dashboard", EnsureAuth, us.Dashboard)
	router.Get("/settings", EnsureAuth, us.Edit)
	router.Get("/settings/:category", EnsureAuth, us.Edit)
	router.Get("/wakatime/sync", EnsureAuth, us.SyncHeartbeats)
}

func (us UserService) SignUp(r render.Render) {
	r.HTML(200, "user/signup", nil)
}

func (us UserService) Create(user models.User, req *http.Request, rw http.ResponseWriter, r render.Render) {
	errs := utils.Verify(user)
	if len(errs) > 0 {
		// some error
		r.HTML(403, "user/signup", nil)
		return
	}

	// user.LastWakaTimeSync = time.Unix(0, 0)
	err := utils.ORM.Save(&user).Error
	if err != nil {
		r.HTML(403, "user/signup", nil)
		return
	}

	session, _ := utils.SessionStore.New(req, "chroniton")
	session.Values["id"] = user.Id
	session.Save(req, rw)

	r.Redirect("/dashboard")
}

func (us UserService) SignIn(r render.Render) {
	r.HTML(200, "user/signin", nil)
}

func (us UserService) SignOut(w http.ResponseWriter, req *http.Request, r render.Render) {
	// delete logged in session
	session, _ := utils.SessionStore.Get(req, "chroniton")
	session.Options.MaxAge = -1
	session.Save(req, w)
	r.Redirect("/")
}

func (us UserService) Edit(currentUser models.User, params martini.Params, r render.Render) {
	category := params["category"]
	if category == "" {
		r.HTML(200, "user/edit", currentUser)
	} else {
		r.HTML(200, "user/edit-"+category, currentUser)
	}
}

func (us UserService) Update(user models.User, req *http.Request, r render.Render) {
	var oldUser models.User
	err := utils.ORM.First(&oldUser, user.Id).Error
	if err != nil {
		r.HTML(500, "user/edit", nil)
		return
	}

	err = req.ParseForm()
	if err != nil {
		r.Error(5000)
	}
	structVal := structs.Map(user)
	changes := utils.GetDiffFormValues(req.PostForm, structVal)
	utils.Log.Debug("changed: %v", changes)
	err = utils.ORM.Table("users").Updates(changes).Error
	if err != nil {
		r.HTML(403, "user/edit", nil)
		return
	}
	var curUser models.User
	utils.ORM.First(&curUser, user.Id)

	if oldUser.WakaTimeApiKey != curUser.WakaTimeApiKey && curUser.WakaTimeApiKey != "" {
		utils.Log.Debug("Pulling down wakatime heartbeats")
		curUser.IsSyncingWakaTime = true
		utils.ORM.Save(&curUser)
		go curUser.PullNewestHeartbeats()
	}

	if strings.Contains(req.Referer(), "/settings") {
		r.Redirect(req.Referer())
	} else {
		r.Redirect("/settings")
	}
}

func (us UserService) SyncHeartbeats(currentUser models.User, r render.Render, req *http.Request) {
	currentUser.IsSyncingWakaTime = true
	utils.ORM.Save(&currentUser)
	go currentUser.PullNewestHeartbeats()

	if req.Referer() != "" {
		r.Redirect(req.Referer())
	} else {
		r.Redirect("/dashboard")
	}
}

func (us UserService) Dashboard(current_user models.User, r render.Render) {
	r.HTML(200, "dashboard", current_user)
}

package utils

import (
	"github.com/martini-contrib/sessions"
)

var (
	SessionStore sessions.CookieStore
)

func init() {
	SessionStore = sessions.NewCookieStore([]byte("chroniton-bad-secret"))
	SessionStore.Options(sessions.Options{
		MaxAge: 0,
	})
}

package auth

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"

	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func NewAuth() {
  clientId := os.Getenv("GOOGLE_CLIENT_ID")
  secret := os.Getenv("GOOGLE_CLIENT_SECRET")
  url := os.Getenv("BASE_URL") + os.Getenv("PORT") +"/auth/callback?provider=google"

//  maxAge := 86400 * 30
//  isProd := false       // Set to true when serving over https
//
  gothic.Store = sessions.NewCookieStore([]byte("randomString"))
//  store.MaxAge(maxAge)
//  store.Options.Path = "/"
//  store.Options.HttpOnly = true   // HttpOnly should always be enabled
//  store.Options.Secure = isProd

  goth.UseProviders(
    google.New(clientId, secret, url),
  )
}

package auth

import (
	"os"
	"time"

	"github.com/aaron-vasilev/diary-templ/src/model"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"

	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func NewAuth() {
  clientId := os.Getenv("GOOGLE_CLIENT_ID")
  secret := os.Getenv("GOOGLE_CLIENT_SECRET")
  url := os.Getenv("BASE_URL") + os.Getenv("PORT") +"/auth/callback?provider=google"

  gothic.Store = sessions.NewCookieStore([]byte("randomString"))

  goth.UseProviders(
    google.New(clientId, secret, url),
  )
}

type UserClaims struct {
  Id         int        `json:"id"`
  Email      string     `json:"email"`
  Name       string     `json:"name"`
  Role       model.Role `json:"role"`
  Subscribed bool       `json:"subscribed"`
  jwt.StandardClaims
}

func newAccessToken(claims UserClaims) (string, error) {
 accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

 return accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func EncodeJWT(u model.User) (string, error) {
  userClaims := UserClaims{
    Id: u.Id,
    Email: u.Email,
    Name: u.Name,
    Role: u.Role,
    Subscribed: u.Subscribed,
    StandardClaims: jwt.StandardClaims{
      IssuedAt: time.Now().Unix(),
      ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
    },
  }

  token, err := newAccessToken(userClaims)

  return token, err
}

func DecodeJWT(accessToken string) (*UserClaims, error) {
 parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
  return []byte(os.Getenv("JWT_SECRET")), nil
 })

 return parsedAccessToken.Claims.(*UserClaims), err
}

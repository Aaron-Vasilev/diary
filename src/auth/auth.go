package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/aaron-vasilev/diary/src/controller"
	"github.com/aaron-vasilev/diary/src/model"
	"github.com/aaron-vasilev/diary/src/utils"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"golang.org/x/crypto/bcrypt"

	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	accessTTL  = 15 * time.Minute
	refreshTTL = 3 * 24 * time.Hour
)

func NewAuth() {
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	url := os.Getenv("BASE_URL") + "auth/callback?provider=google"

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

type RefreshClaims struct {
	Id int `json:"id"`
	jwt.StandardClaims
}

func jwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func EncodeAccess(u model.User) (string, error) {
	claims := UserClaims{
		Id:         u.Id,
		Email:      u.Email,
		Name:       u.Name,
		Role:       u.Role,
		Subscribed: u.Subscribed,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(accessTTL).Unix(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret())
}

func EncodeRefresh(u model.User) (string, error) {
	claims := RefreshClaims{
		Id: u.Id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(refreshTTL).Unix(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret())
}

func DecodeAccess(token string) (*UserClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (any, error) {
		return jwtSecret(), nil
	})
	if err != nil || parsed == nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*UserClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid access token")
	}
	return claims, nil
}

func DecodeRefresh(token string) (*RefreshClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &RefreshClaims{}, func(t *jwt.Token) (any, error) {
		return jwtSecret(), nil
	})
	if err != nil || parsed == nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*RefreshClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}

func setAuthCookie(c echo.Context, name, value string, ttl time.Duration) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(ttl)
	cookie.MaxAge = int(ttl.Seconds())
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Secure = utils.IsProd()
	c.SetCookie(cookie)
}

func SetAuthCookies(c echo.Context, u model.User) error {
	access, err := EncodeAccess(u)
	if err != nil {
		return err
	}
	refresh, err := EncodeRefresh(u)
	if err != nil {
		return err
	}
	setAuthCookie(c, utils.TOKEN, access, accessTTL)
	setAuthCookie(c, utils.REFRESH_TOKEN, refresh, refreshTTL)
	return nil
}

func ClearAuthCookies(c echo.Context) {
	utils.DeleteCookie(c, utils.TOKEN)
	utils.DeleteCookie(c, utils.REFRESH_TOKEN)
}

func GetUserClaimsFromCtx(c echo.Context) (*UserClaims, error) {
	if accessCookie, err := c.Cookie(utils.TOKEN); err == nil {
		if claims, err := DecodeAccess(accessCookie.Value); err == nil {
			return claims, nil
		}
	}

	refreshCookie, err := c.Cookie(utils.REFRESH_TOKEN)
	if err != nil {
		ClearAuthCookies(c)
		return nil, err
	}

	refreshClaims, err := DecodeRefresh(refreshCookie.Value)
	if err != nil {
		ClearAuthCookies(c)
		return nil, err
	}

	user := controller.GetUserById(refreshClaims.Id)
	if user.Id == 0 {
		ClearAuthCookies(c)
		return nil, errors.New("user not found")
	}

	if err := SetAuthCookies(c, user); err != nil {
		return nil, err
	}

	return &UserClaims{
		Id:         user.Id,
		Email:      user.Email,
		Name:       user.Name,
		Role:       user.Role,
		Subscribed: user.Subscribed,
	}, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

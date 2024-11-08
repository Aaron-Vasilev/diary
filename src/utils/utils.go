package utils

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	TOKEN string = "token"
)

func DateStrIsValid(dateStr string) bool {
	if dateStr == "" {
		return false
	}

	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func DeleteCookie(c echo.Context, tokenName string) {
	cookie := new(http.Cookie)
	cookie.Name = tokenName
	cookie.Expires = time.Unix(0, 0)
	cookie.Path = "/"
	c.SetCookie(cookie)
}

func IsProd() bool {
	return os.Getenv("ENV") == "production"
}

func BeautyDate(date string) string {
	return strings.Split(date, "T")[0]
}

func Int(i int) string {
	return strconv.Itoa(i)
}

func PublicUrl(path string) string {
	return os.Getenv("PUBLIC_URL") + path
}

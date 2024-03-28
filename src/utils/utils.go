package utils

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
  TOKEN string = "token"
)

func ValidateDateString(dateStr string) bool {
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
package utils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

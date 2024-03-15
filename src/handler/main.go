package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/aaron-vasilev/diary-templ/src/controller"
	"github.com/aaron-vasilev/diary-templ/src/model"
	"github.com/aaron-vasilev/diary-templ/src/pages"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type HandlerCtx struct {
  Db *sql.DB
}

func (h HandlerCtx) QuestionListHandler(ctx echo.Context) error {
  questions, err := controller.GetQuestions(h.Db)

  if err != nil {
    return nil
  }

  return pages.QuestionList(questions).Render(ctx.Request().Context(), ctx.Response())
}

func (h HandlerCtx) Diary(c echo.Context) error {
//   shownDate := ctx.QueryParam("shown-date")
//   if shownDate == "" {
//     shownDate = "2023-10-10"
//   }
  user, ok := c.Get(string(userContextKey)).(model.User)
  fmt.Println("† line 34 user", user)
  fmt.Println("† line 34 ok", ok)

  return pages.Diary(pages.DiaryProps{
    User: user,
  }).Render(c.Request().Context(), c.Response())
}

func (h HandlerCtx) Login(ctx echo.Context) error {
  gothic.BeginAuthHandler(ctx.Response().Writer, ctx.Request())
  return nil
}

type contextKey string

const userContextKey contextKey = "user"

func (h HandlerCtx) AuthCallback(c echo.Context) error {
  googleUser, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
  if err != nil {
    fmt.Println("† line 44 err", err)
    return c.Redirect(http.StatusFound, "/login")
  }

  user := controller.GetUserByEmail(h.Db, googleUser.Email)

  c.Set(string(userContextKey), user)

  return c.Redirect(http.StatusFound, "/diary")
}

// func ValidateDateString(dateStr string) bool {
//     layout := "2006-01-02" // Reference layout for "YYYY-MM-DD"
//     _, err := time.Parse(layout, dateStr)
//     return err == nil // If there's no error, the format is correct
// }


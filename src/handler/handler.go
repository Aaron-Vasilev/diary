package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/aaron-vasilev/diary-templ/src/auth"
	"github.com/aaron-vasilev/diary-templ/src/controller"
	"github.com/aaron-vasilev/diary-templ/src/model"
	"github.com/aaron-vasilev/diary-templ/src/pages"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type HandlerCtx struct {
  Db *sql.DB
}

const (
  TOKEN string = "token"
)

func (h HandlerCtx) QuestionListHandler(ctx echo.Context) error {
  questions := controller.GetQuestions(h.Db)

  return pages.QuestionList(questions).Render(ctx.Request().Context(), ctx.Response())
}

func (h HandlerCtx) Diary(c echo.Context) error {
  var user     model.User
  var question model.Question
  var notes  []model.Note

  shownDate := c.QueryParam("shown-date")

  if shownDate == "" {
    question.ShownDate = time.Now().Format("2006-01-02") 
  }

  cookies, err := c.Cookie(TOKEN)
  
  if err == nil {
    token := cookies.Value
    userClaim, err := auth.DecodeJWT(token)
    
    if err != nil {
      deleteCookie(c, TOKEN)
    } else {
      user = controller.GetUserByEmail(h.Db, userClaim.Email)
      question = controller.GetQuestionByDate(h.Db, question.ShownDate)
      notes = controller.GetNotes(h.Db, user.Id, question.Id)
    }
  }

  if user.Id == 0 {
    user.Name = "Anon"
    user.Role = "user"
  }

  return pages.Diary(pages.DiaryProps{
    User: user,
    Question: question,
    Notes: notes,
  }).Render(c.Request().Context(), c.Response())
}

func (h HandlerCtx) LoginPage(c echo.Context) error {
  cookies, err := c.Cookie(TOKEN)

  if err == nil {
    token := cookies.Value
    _, err := auth.DecodeJWT(token)
    
    if err == nil {
      return c.Redirect(http.StatusFound, "/diary")
    }
  }

  return pages.Login().Render(c.Request().Context(), c.Response())
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
    return c.Redirect(http.StatusFound, "/login")
  }

  user := controller.GetUserByEmail(h.Db, googleUser.Email)
  token, err := auth.EncodeJWT(user)

  if err != nil {
    return c.Redirect(http.StatusFound, "/login")
  }

	cookie := new(http.Cookie)
	cookie.Name = TOKEN
	cookie.Value = token
	cookie.Path = "/"
  c.SetCookie(cookie)

  return c.Redirect(http.StatusFound, "/diary")
}

// func ValidateDateString(dateStr string) bool {
//     layout := "2006-01-02" // Reference layout for "YYYY-MM-DD"
//     _, err := time.Parse(layout, dateStr)
//     return err == nil // If there's no error, the format is correct
// }

func deleteCookie(c echo.Context, tokenName string) {
  cookie := new(http.Cookie)
  cookie.Name = tokenName
  cookie.Expires = time.Unix(0, 0)
  cookie.Path = "/"
  c.SetCookie(cookie)
}

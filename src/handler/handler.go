package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/aaron-vasilev/diary-templ/src/auth"
	"github.com/aaron-vasilev/diary-templ/src/components"
	"github.com/aaron-vasilev/diary-templ/src/controller"
	"github.com/aaron-vasilev/diary-templ/src/model"
	"github.com/aaron-vasilev/diary-templ/src/pages"
	"github.com/aaron-vasilev/diary-templ/src/utils"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type HandlerCtx struct {
  Db *sql.DB
}

func (h HandlerCtx) QuestionListHandler(ctx echo.Context) error {
  questions := controller.GetQuestions(h.Db)

  return pages.QuestionList(false, questions).Render(ctx.Request().Context(), ctx.Response())
}

func (h HandlerCtx) Diary(c echo.Context) error {
  isLogin := false
  var question model.Question
  var notes  []model.Note
  user := model.User{
    Name: "Anon",
    Role: "user",
  }

  shownDate := c.QueryParam("shown-date")

  if shownDate == "" {
    question.ShownDate = time.Now().Format("2006-01-02") 
  }

  userClaims, err := auth.GetUserClaimsFromCtx(c)
  
  if err == nil {
    isLogin = true
    user = controller.GetUserByEmail(h.Db, userClaims.Email)
    question = controller.GetQuestionByDate(h.Db, question.ShownDate)
    notes = controller.GetNotes(h.Db, user.Id, question.Id)
  }

  return pages.Diary(components.DiaryProps{
    IsLogin: isLogin,
    User: user,
    Question: question,
    Notes: notes,
  }).Render(c.Request().Context(), c.Response())
}

func (h HandlerCtx) LoginPage(c echo.Context) error {
  logoutStr := c.QueryParam("logout")
  logout, err := strconv.ParseBool(logoutStr)

  if err == nil && logout {
    utils.DeleteCookie(c, utils.TOKEN)
    return c.Redirect(http.StatusFound, "/login")
  }

  _, err = auth.GetUserClaimsFromCtx(c)
    
  if err == nil {
    return c.Redirect(http.StatusFound, "/diary")
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
	cookie.Name = utils.TOKEN
	cookie.Value = token
	cookie.Path = "/"
  c.SetCookie(cookie)

  return c.Redirect(http.StatusFound, "/diary")
}

package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/aaron-vasilev/diary/src/auth"
	"github.com/aaron-vasilev/diary/src/components"
	"github.com/aaron-vasilev/diary/src/controller"
	"github.com/aaron-vasilev/diary/src/model"
	"github.com/aaron-vasilev/diary/src/pages"
	"github.com/aaron-vasilev/diary/src/utils"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type HandlerCtx struct {
	Db *sql.DB
}

func (h HandlerCtx) Home(c echo.Context) error {
	question := controller.GetQuestionByDate(h.Db, "2023-08-05")

	_, err := auth.GetUserClaimsFromCtx(c)

	if err == nil {
		return c.Redirect(http.StatusFound, "/diary")
	}

	return pages.Home(pages.HomeProps{
		Question: question,
	}).Render(c.Request().Context(), c.Response())
}

func (h HandlerCtx) QuestionListHandler(c echo.Context) error {
	questions := controller.GetQuestions(h.Db)

	_, err := auth.GetUserClaimsFromCtx(c)

	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	return pages.QuestionList(questions).Render(c.Request().Context(), c.Response())
}

func (h HandlerCtx) Diary(c echo.Context) error {
	var question model.Question
	var notes []model.Note
	user := model.User{
		Name: "Anon",
		Role: "user",
	}

	shownDate := c.QueryParam("shown-date")

	if !utils.DateStrIsValid(shownDate) {
		question.ShownDate = time.Now().Format("2006-01-02")
	} else {
		question.ShownDate = shownDate
	}

	userClaims, err := auth.GetUserClaimsFromCtx(c)

	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	question = controller.GetQuestionByDate(h.Db, question.ShownDate)
	user, err = controller.GetUserByEmail(h.Db, userClaims.Email)
	notes = controller.GetNotes(h.Db, user.Id, question.Id)

	return pages.Diary(components.DiaryProps{
		User:     user,
		Question: question,
		Notes:    notes,
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
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	user, err := controller.GetUserByEmail(h.Db, email)

	if err != nil {
		return err
	}

	if user.Password == nil {
		return nil
	}

	isValidPassword := auth.CheckPassword(password, *user.Password)

	if !isValidPassword {
		return nil
	}

	token, err := auth.EncodeJWT(user)

	if err != nil {
		return nil
	}

	cookie := new(http.Cookie)
	cookie.Name = utils.TOKEN
	cookie.Value = token
	cookie.Path = "/"
	ctx.SetCookie(cookie)

	return ctx.Redirect(http.StatusFound, "/diary")
}

func (h HandlerCtx) Register(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	name := ctx.FormValue("name")

	_, err := controller.CreateUser(h.Db, email, password, name)

	if err != nil {
		return nil
	}

	return ctx.Redirect(http.StatusFound, "/diary")
}

type contextKey string

const userContextKey contextKey = "user"

func (h HandlerCtx) AuthCallback(c echo.Context) error {
	googleUser, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())

	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	user, err := controller.GetUserByEmail(h.Db, googleUser.Email)
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

func (h HandlerCtx) UpdateQuestion(c echo.Context) error {
	_, err := auth.GetUserClaimsFromCtx(c)

	if err != nil {
		return c.Redirect(http.StatusFound, "/diary")
	}

	var question model.Question
	shownDate := c.QueryParam("shown-date")

	if !utils.DateStrIsValid(shownDate) {
		question.ShownDate = time.Now().Format("2006-01-02")
	} else {
		question.ShownDate = shownDate
	}

	question = controller.GetQuestionByDate(h.Db, question.ShownDate)

	return pages.UpdateQuestion(pages.UpdateQuestionProps{
		Question: question,
		User: model.User{
			Id:   1,
			Name: "Aaron",
		},
	}).Render(c.Request().Context(), c.Response())
}

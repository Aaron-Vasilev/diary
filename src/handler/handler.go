package handler

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aaron-vasilev/diary/src/auth"
	"github.com/aaron-vasilev/diary/src/components"
	"github.com/aaron-vasilev/diary/src/controller"
	"github.com/aaron-vasilev/diary/src/model"
	"github.com/aaron-vasilev/diary/src/pages"
	"github.com/aaron-vasilev/diary/src/telegram"
	"github.com/aaron-vasilev/diary/src/utils"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type HandlerCtx struct{}

func (h HandlerCtx) Home(c echo.Context) error {
	question := controller.GetQuestionByDate("2023-08-05")

	_, err := auth.GetUserClaimsFromCtx(c)

	if err == nil {
		return c.Redirect(http.StatusFound, "/diary")
	}

	return pages.Home(pages.HomeProps{
		Question: question,
	}).Render(c.Request().Context(), c.Response())
}

func (h HandlerCtx) QuestionListHandler(c echo.Context) error {
	_, err := auth.GetUserClaimsFromCtx(c)

	if err != nil {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	questions, err := controller.GetQuestions()

	if err != nil {
		log.Printf("QuestionListHandler: GetQuestions failed: %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return pages.QuestionList(questions).Render(c.Request().Context(), c.Response())
}

func (h HandlerCtx) NoteListHandler(c echo.Context) error {
	_, err := auth.GetUserClaimsFromCtx(c)

	if err != nil {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	var notes []model.Note

	return pages.NoteList(notes).Render(c.Request().Context(), c.Response())
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

	question = controller.GetQuestionByDate(question.ShownDate)
	user, err = controller.GetUserByEmail(userClaims.Email)
	notes = controller.GetNotes(user.Id, question.Id)

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
		auth.ClearAuthCookies(c)
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

	user, err := controller.GetUserByEmail(email)

	if err != nil {
		log.Printf("Login: GetUserByEmail(%q) failed: %v", email, err)
		return ctx.Redirect(http.StatusFound, "/login")
	}

	if user.Password == nil {
		log.Printf("Login: user %q has no password set", email)
		return ctx.Redirect(http.StatusFound, "/login")
	}

	isValidPassword := auth.CheckPassword(password, *user.Password)

	if !isValidPassword {
		log.Printf("Login: bad password for %q", email)
		return ctx.Redirect(http.StatusFound, "/login")
	}

	if err := auth.SetAuthCookies(ctx, user); err != nil {
		log.Printf("Login: SetAuthCookies for %q failed: %v", email, err)
		return ctx.NoContent(http.StatusUnauthorized)
	}

	return ctx.Redirect(http.StatusFound, "/diary")
}

func (h HandlerCtx) Register(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	name := ctx.FormValue("name")

	_, err := controller.CreateUser(email, password, name)

	if err != nil {
		log.Printf("Register: CreateUser(%q) failed: %v", email, err)
		return ctx.Redirect(http.StatusFound, "/login")
	}

	return ctx.Redirect(http.StatusFound, "/diary")
}

type contextKey string

const userContextKey contextKey = "user"

func (h HandlerCtx) AuthCallback(c echo.Context) error {
	googleUser, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())

	if err != nil {
		log.Printf("AuthCallback: CompleteUserAuth failed: %v", err)
		return c.Redirect(http.StatusFound, "/login")
	}

	user, err := controller.GetUserByEmail(googleUser.Email)

	if err != nil {
		log.Printf("AuthCallback: GetUserByEmail(%q) failed: %v", googleUser.Email, err)
		return c.Redirect(http.StatusFound, "/login")
	}

	if err := auth.SetAuthCookies(c, user); err != nil {
		log.Printf("AuthCallback: SetAuthCookies for %q failed: %v", googleUser.Email, err)
		return c.Redirect(http.StatusFound, "/login")
	}

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

	question = controller.GetQuestionByDate(question.ShownDate)

	return pages.UpdateQuestion(pages.UpdateQuestionProps{
		Question: question,
		User: model.User{
			Id:   1,
			Name: "Aaron",
		},
	}).Render(c.Request().Context(), c.Response())
}

func (h HandlerCtx) TelegramLogin(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("TelegramLogin: read body failed: %v", err)
		return c.NoContent(http.StatusBadRequest)
	}

	tgUser, err := telegram.ParseInitData(string(body))
	if err != nil {
		log.Printf("TelegramLogin: ParseInitData failed: %v", err)
		return c.NoContent(http.StatusUnauthorized)
	}

	name := tgUser.FirstName
	if tgUser.LastName != "" {
		name += " " + tgUser.LastName
	}

	user, err := controller.UpsertTelegramUser(tgUser.ID, name)
	if err != nil {
		log.Printf("TelegramLogin: UpsertTelegramUser(%d, %q) failed: %v", tgUser.ID, name, err)
		return c.NoContent(http.StatusInternalServerError)
	}

	if err := auth.SetAuthCookies(c, user); err != nil {
		log.Printf("TelegramLogin: SetAuthCookies for tg user %d failed: %v", tgUser.ID, err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusFound, "/diary")
}

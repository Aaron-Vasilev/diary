package router

import (
	"database/sql"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/aaron-vasilev/diary/src/auth"
	"github.com/aaron-vasilev/diary/src/components"
	"github.com/aaron-vasilev/diary/src/controller"
	"github.com/aaron-vasilev/diary/src/handler"
	"github.com/aaron-vasilev/diary/src/model"
	"github.com/aaron-vasilev/diary/src/pages"
	"github.com/labstack/echo/v4"
)

func ConnectRoutes(app *echo.Echo, db *sql.DB) {
	// Pages
	app.GET("/", handler.HandlerCtx{Db: db}.Home)
	app.GET("/question-list", handler.HandlerCtx{Db: db}.QuestionListHandler)
	app.GET("/note-list", handler.HandlerCtx{Db: db}.NoteListHandler)
	app.GET("/diary", handler.HandlerCtx{Db: db}.Diary)
	app.GET("/login", handler.HandlerCtx{Db: db}.LoginPage)
	app.GET("/update-question", handler.HandlerCtx{Db: db}.UpdateQuestion)

	app.GET("/auth/login", handler.HandlerCtx{Db: db}.Login)
	app.GET("/auth/callback", handler.HandlerCtx{Db: db}.AuthCallback) // for Google
	app.GET("/test", func(c echo.Context) error {
		fmt.Println("✡️  line 24 test")
		return c.String(http.StatusOK, "test")
	})
	app.GET("/*", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/")
	})

	// Handlers
	app.POST("/note", func(c echo.Context) error {
		noteText := c.Request().FormValue("note")
		date := c.Request().FormValue("createdDate")
		questionIdStr := c.QueryParam("question_id")
		questionId, err := strconv.Atoi(questionIdStr)

		if err != nil {
			return c.NoContent(http.StatusNotAcceptable)
		}

		userClaims, err := auth.GetUserClaimsFromCtx(c)

		if err == nil && noteText != "" && date != "" {
			n := controller.CreateNote(db, userClaims.Id, questionId, date, noteText)

			return components.Note(n).Render(c.Request().Context(), c.Response())
		}

		return c.NoContent(http.StatusNotAcceptable)
	})

	app.GET("/note/:id", func(c echo.Context) error {
		_, err := auth.GetUserClaimsFromCtx(c)

		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		id, _ := strconv.Atoi(c.Param("id"))
		n := controller.GetNoteById(db, id)

		return components.Note(n).Render(c.Request().Context(), c.Response())
	})

	app.PUT("/note/:id", func(c echo.Context) error {
		var n model.Note
		_, err := auth.GetUserClaimsFromCtx(c)

		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.NoContent(http.StatusNotAcceptable)
		}

		changedText := c.Request().FormValue("text")

		if changedText == "" {
			n = controller.GetNoteById(db, id)

			return components.EditNote(n).Render(c.Request().Context(), c.Response())
		} else {
			n = controller.UpdateNote(db, id, changedText)

			return components.Note(n).Render(c.Request().Context(), c.Response())
		}
	})

	app.DELETE("/note/:id", func(c echo.Context) error {
		_, err := auth.GetUserClaimsFromCtx(c)

		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.NoContent(http.StatusNotAcceptable)
		}

		controller.DeleteNote(db, id)

		return c.NoContent(http.StatusOK)
	})

	app.POST("/change-date", func(c echo.Context) error {
		question := controller.GetQuestionByDate(db, c.FormValue("date"))
		var notes []model.Note
		user := model.User{
			Name: "Anon",
			Role: "user",
		}

		userClaims, err := auth.GetUserClaimsFromCtx(c)

		if err == nil {
			user, err = controller.GetUserByEmail(db, userClaims.Email)
			notes = controller.GetNotes(db, user.Id, question.Id)
		}

		return components.Diary(components.DiaryProps{
			User:     user,
			Question: question,
			Notes:    notes,
		}).Render(c.Request().Context(), c.Response())
	})

	app.POST("/question-search", func(c echo.Context) error {
		search := c.FormValue("search")

		_, err := auth.GetUserClaimsFromCtx(c)

		if err != nil {
			return c.Redirect(http.StatusUnauthorized, "/login")
		}

		questions, err := controller.GetQuestionsLike(db, search)

		if err != nil {
			fmt.Println("✡️  line 151 err", err)
			c.String(http.StatusBadRequest, err.Error())
		}

		return components.QuestionList(questions).Render(c.Request().Context(), c.Response())
	})

	app.POST("/update-question", func(c echo.Context) error {
		question := controller.GetQuestionByDate(db, c.FormValue("date"))
		user := model.User{
			Name: "Aaron",
			Id:   1,
		}

		return components.Question(question, user).Render(c.Request().Context(), c.Response())
	})

	app.PUT("/update-question", func(c echo.Context) error {
		_, err := auth.GetUserClaimsFromCtx(c)

		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		newQuestion := c.FormValue("question")
		questionIdStr := c.QueryParam("id")
		id, _ := strconv.Atoi(questionIdStr)

		question := controller.UpdateQuestion(db, id, newQuestion)
		user := model.User{
			Name: "Aaron",
			Id:   1,
		}

		return components.Question(question, user).Render(c.Request().Context(), c.Response())
	})

	app.GET("/random-question", func(c echo.Context) error {
		var question model.Question

		question = controller.GetQuestion(db, rand.IntN(360-1)+1)

		return components.RandomQuestion(question).Render(c.Request().Context(), c.Response())
	})

	app.POST("/note-search", func(c echo.Context) error {
		var notes []model.Note
		search := c.FormValue("search")
		user, err := auth.GetUserClaimsFromCtx(c)

		if err != nil {
			return c.Redirect(http.StatusUnauthorized, "/login")
		} else if len(search) > 1 {
			notes, err = controller.GetNotesByText(db, user.Id, search)

			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}

		return pages.NoteHistory(notes).Render(c.Request().Context(), c.Response())
	})
}

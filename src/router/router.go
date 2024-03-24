package router

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aaron-vasilev/diary-templ/src/auth"
	"github.com/aaron-vasilev/diary-templ/src/components"
	"github.com/aaron-vasilev/diary-templ/src/controller"
	"github.com/aaron-vasilev/diary-templ/src/handler"
	"github.com/aaron-vasilev/diary-templ/src/model"
	"github.com/labstack/echo/v4"
)

func ConnectRoutes(app *echo.Echo, db *sql.DB) {
  app.GET("/question-list", handler.HandlerCtx{ Db: db, }.QuestionListHandler)
  app.GET("/diary", handler.HandlerCtx{ Db: db, }.Diary)
  app.GET("/login", handler.HandlerCtx{ Db: db, }.LoginPage)

  app.GET("/auth/login", handler.HandlerCtx{ Db: db }.Login)
  app.GET("/auth/callback", handler.HandlerCtx{ Db: db }.AuthCallback)
  app.GET("/test", func(c echo.Context) error {
    fmt.Println("✡️  line 24 test")
    return c.String(http.StatusOK, "test")
  })


  app.POST("/note", func(c echo.Context) error {
    noteText := c.Request().FormValue("note")
    date := c.QueryParam("created_date")
    questionIdStr := c.QueryParam("question_id")
    questionId, err := strconv.Atoi(questionIdStr)

    if err != nil { return nil }

    userClaims, err := auth.GetUserClaimsFromCtx(c)

    if err == nil && noteText != "" {
      n := controller.CreateNote(db, userClaims.Id, questionId, date, noteText)

      return components.Note(n).Render(c.Request().Context(), c.Response())
    }

    return nil
  })

  app.GET("/note/:id", func(c echo.Context) error {
    _, err := auth.GetUserClaimsFromCtx(c)

    if err != nil { return c.NoContent(http.StatusUnauthorized) }

    id, _ := strconv.Atoi(c.Param("id"))
    n := controller.GetNoteById(db, id)

    return components.Note(n).Render(c.Request().Context(), c.Response())
  })

  app.PUT("/note/:id", func(c echo.Context) error {
    var n model.Note
    _, err := auth.GetUserClaimsFromCtx(c)

    if err != nil { return c.NoContent(http.StatusUnauthorized) }

    id, err := strconv.Atoi(c.Param("id"))

    if err != nil { return c.NoContent(http.StatusNotAcceptable) }

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

    if err != nil { return c.NoContent(http.StatusUnauthorized) }

    id, err := strconv.Atoi(c.Param("id"))

    if err != nil { return c.NoContent(http.StatusNotAcceptable) }

    controller.DeleteNote(db, id)

    return c.NoContent(http.StatusOK)
  })

   app.POST("/change-date", func(c echo.Context) error {
    var question model.Question
    var notes  []model.Note
    user := model.User{
      Name: "Anon",
      Role: "user",
    }

    userClaims, err := auth.GetUserClaimsFromCtx(c)

    if err == nil {
      user = controller.GetUserByEmail(db, userClaims.Email)
      question = controller.GetQuestionByDate(db, c.FormValue("date"))
      notes = controller.GetNotes(db, user.Id, question.Id)
    }

    return components.Diary(components.DiaryProps{
      User: user,
      Question: question,
      Notes: notes,
    }).Render(c.Request().Context(), c.Response())
  })
}

package router

import (
	"database/sql"

	"github.com/aaron-vasilev/diary-templ/src/components"
	"github.com/aaron-vasilev/diary-templ/src/handler"
	"github.com/aaron-vasilev/diary-templ/src/model"
	"github.com/labstack/echo/v4"
)

func Connect(app *echo.Echo, db *sql.DB) {
  app.GET("/question-list", handler.HandlerCtx{ Db: db, }.QuestionListHandler)
  app.GET("/diary", handler.HandlerCtx{ Db: db, }.Diary)
  app.GET("/login", handler.HandlerCtx{ Db: db, }.LoginPage)

  app.GET("/auth/login", handler.HandlerCtx{ Db: db }.Login)
  app.GET("/auth/callback", handler.HandlerCtx{ Db: db }.AuthCallback)

  app.POST("/note", func(c echo.Context) error {
    noteText := c.Request().FormValue("note")

    n := model.Note{
      Id: 88,
      UserId: 1,
      Text: noteText,
      CreatedDate: "2022-01-01",
      QuestionId: 1,
    }

    return components.Note(n).Render(c.Request().Context(), c.Response())
  })
}

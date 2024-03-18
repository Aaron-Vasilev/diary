package router

import (
	"database/sql"
	"fmt"

	"github.com/aaron-vasilev/diary-templ/src/auth"
	"github.com/aaron-vasilev/diary-templ/src/components"
	"github.com/aaron-vasilev/diary-templ/src/controller"
	"github.com/aaron-vasilev/diary-templ/src/handler"
	"github.com/aaron-vasilev/diary-templ/src/model"
	"github.com/aaron-vasilev/diary-templ/src/utils"
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

   app.POST("/change-date", func(c echo.Context) error {
    var user     model.User
    var question model.Question
    var notes  []model.Note

    cookies, err := c.Cookie(handler.TOKEN)

    if err == nil {
      token := cookies.Value
      userClaim, err := auth.DecodeJWT(token)

      if err != nil {
        utils.DeleteCookie(c, handler.TOKEN)
      } else {
        user = controller.GetUserByEmail(db, userClaim.Email)
        question = controller.GetQuestionByDate(db, c.FormValue("date"))
        notes = controller.GetNotes(db, user.Id, question.Id)
      }
    }
    fmt.Println("† line 53 question", question.ShownDate)

    if user.Id == 0 {
      user.Name = "Anon"
      user.Role = "user"
    }

    return components.Diary(components.DiaryProps{
      User: user,
      Question: question,
      Notes: notes,
    }).Render(c.Request().Context(), c.Response())
  })
}

package handler

import (
	"database/sql"

	"github.com/aaron-vasilev/diary-templ/src/controller"
	"github.com/aaron-vasilev/diary-templ/src/pages"
	"github.com/labstack/echo/v4"
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

func (h HandlerCtx) DiaryHandler(ctx echo.Context) error {
  shownDate := ctx.QueryParam("shown-date")

  if shownDate == "" {
    shownDate = "2023-10-10"
  }

  return pages.Diary(shownDate).Render(ctx.Request().Context(), ctx.Response())
}

// func ValidateDateString(dateStr string) bool {
//     layout := "2006-01-02" // Reference layout for "YYYY-MM-DD"
//     _, err := time.Parse(layout, dateStr)
//     return err == nil // If there's no error, the format is correct
// }


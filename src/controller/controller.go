package controller

import (
	"database/sql"
	"fmt"
)

type Question struct {
  ID int
  Text string
  ShownDate string
}

func GetQuestionByDate(db *sql.DB, date string) Question {
  var q Question
  query := `SELECT * FROM diary.question q WHERE q.shown_date=$1;`
  row := db.QueryRow(query, date).Scan(&q.ID, &q.Text, &q.ShownDate)

  if row == sql.ErrNoRows {
    fmt.Println("No question with that date")
  }

  return q
}

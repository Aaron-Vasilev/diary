package controller

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/aaron-vasilev/diary-templ/src/model"
)

func GetQuestionByDate(db *sql.DB, date string) model.Question {
  var q model.Question
  query := `SELECT * FROM diary.question q WHERE q.shown_date=$1;`
  row := db.QueryRow(query, date).Scan(&q.Id, &q.Text, &q.ShownDate)

  if row == sql.ErrNoRows {
    fmt.Println("No question with that date")
  }

  return q
}

func GetUserById(db *sql.DB, id int) model.User {
  var u model.User
  query := `SELECT * FROM diary.user u WHERE u.id=$1;`
  row := db.QueryRow(query, id).Scan(
    &u.Id,
    &u.CreatedAt,
    &u.Email,
    &u.Name,
    &u.Role,
    &u.SubId,
    &u.Subscribed)

  if row == sql.ErrNoRows {
    fmt.Println("No question with that id")
  }

  return u
}

func GetQuestions(db *sql.DB) ([]model.Question, error) {
  var questions []model.Question
  query := `SELECT * FROM diary.question ORDER BY shown_date ASC;`

  rows, err := db.Query(query)

  if err != nil {
    return nil, err
  }
  defer rows.Close()

  for rows.Next() {
    var q model.Question
    var dateString string

    rows.Scan(&q.Id, &q.Text, &dateString)

    dateWithoutTime := strings.Split(dateString, "T")
    q.ShownDate = dateWithoutTime[0]

    questions = append(questions, q)
  }

  return questions, nil
}

type diaryPageData struct {

}

func FetchDataForDiaryPage(db *sql.DB) {
//   query := `SELECT * FROM diary.question q LEFT JOIN diary.note n 
//   ON q.id = n.question_id WHERE q.shown_date=$1 AND n.user_id=$2;`
//   rows, err := db.Query(query, "2024-03-09", 1)
}

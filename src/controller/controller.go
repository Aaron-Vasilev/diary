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

  q.ShownDate = strings.Split(q.ShownDate, "T")[0]

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

func GetUserByEmail(db *sql.DB, email string) model.User {
  var u model.User
  query := `SELECT * FROM diary.user u WHERE u.email=$1;`
  row := db.QueryRow(query, email).Scan(
    &u.Id,
    &u.CreatedAt,
    &u.Email,
    &u.Name,
    &u.Role,
    &u.SubId,
    &u.Subscribed)

  if row == sql.ErrNoRows {
    fmt.Println("No question with that email")
  }

  return u
}

func GetQuestions(db *sql.DB) []model.Question {
  var questions []model.Question
  query := `SELECT * FROM diary.question ORDER BY shown_date ASC;`

  rows, err := db.Query(query)

  if err == nil {
    for rows.Next() {
      var q model.Question
      var dateString string

      rows.Scan(&q.Id, &q.Text, &dateString)

      dateWithoutTime := strings.Split(dateString, "T")
      q.ShownDate = dateWithoutTime[0]

      questions = append(questions, q)
    }
  }
  defer rows.Close()

  return questions
}

func GetNotes(db *sql.DB, userId, questionId int) []model.Note {
  var notes []model.Note
  query := `SELECT * FROM diary.note WHERE user_id=$1 AND question_id=$2;`

  rows, err := db.Query(query, userId, questionId)

  if err == nil {
    for rows.Next() {
      var n model.Note

      rows.Scan(&n.Id, &n.UserId, &n.Text, &n.CreatedDate, &n.QuestionId)

      n.CreatedDate = strings.Split(n.CreatedDate, "T")[0]

      notes = append(notes, n)
    }
  }
  defer rows.Close()

  return notes
}

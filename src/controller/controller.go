package controller

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/aaron-vasilev/diary-templ/src/model"
	"github.com/aaron-vasilev/diary-templ/src/utils"
)

func GetQuestionByDate(db *sql.DB, date string) model.Question {
  var q model.Question
  q.ShownDate = date

  if !utils.DateStrIsValid(date) {
    return q
  } 

  query := `SELECT id, text FROM diary.question q WHERE q.shown_date=$1;`
  splitedDate := strings.Split(date, "-")
  querDate := "2024-" + splitedDate[1] + "-" + splitedDate[2]

  row := db.QueryRow(query, querDate).Scan(&q.Id, &q.Text)

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

func GetQuestionsLike(db *sql.DB, search string) []model.Question {
  var questions []model.Question
  query := `SELECT * FROM diary.question q WHERE q.text LIKE $1 ORDER BY shown_date ASC;`

  rows, err := db.Query(query, "%" + search + "%")

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

func UpdateQuestion(db *sql.DB, id int, text string) model.Question {
  var q model.Question
  query := `UPDATE diary.question SET text=$1 WHERE id=$2 RETURNING *;`

  db.QueryRow(query, text, id).Scan(
    &q.Id,
    &q.Text,
    &q.ShownDate,
  )

  q.ShownDate = utils.BeautyDate(q.ShownDate)

  return q
}


func GetNoteById(db *sql.DB, id int) model.Note {
  var n model.Note
  query := `SELECT * FROM diary.note WHERE id=$1`

  db.QueryRow(query, id).Scan(
    &n.Id,
    &n.UserId,
    &n.Text,
    &n.CreatedDate,
    &n.QuestionId,
  )

  n.CreatedDate = strings.Split(n.CreatedDate, "T")[0]

  return n
}

func UpdateNote(db *sql.DB, id int, text string) model.Note {
  var n model.Note
  query := `UPDATE diary.note SET text=$1 WHERE id=$2 RETURNING *;`

  db.QueryRow(query, text, id).Scan(
    &n.Id,
    &n.UserId,
    &n.Text,
    &n.CreatedDate,
    &n.QuestionId,
  )

  n.CreatedDate = utils.BeautyDate(n.CreatedDate)

  return n
}

func GetNotes(db *sql.DB, userId, questionId int) []model.Note {
  var notes []model.Note
  query := `SELECT * FROM diary.note WHERE user_id=$1 AND question_id=$2 ORDER BY id ASC;`

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

func CreateNote(db *sql.DB, userId, questionId int, createdDate, text string) model.Note {
  note := model.Note{
    UserId: userId,
    QuestionId: questionId,
    Text: text,
    CreatedDate: createdDate,
  }
  query := "INSERT INTO diary.note (user_id, text, created_date, question_id) VALUES ($1, $2, $3, $4) returning id;"

  db.QueryRow(query, userId, text, createdDate, questionId).Scan(&note.Id)

  return note
}

func DeleteNote(db *sql.DB, id int) {
  db.Exec("DELETE FROM diary.note WHERE id=$1", id)
}

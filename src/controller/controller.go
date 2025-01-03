package controller

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/aaron-vasilev/diary/src/model"
	"github.com/aaron-vasilev/diary/src/utils"
)

func GetQuestion(db *sql.DB, id int) model.Question {
	var q model.Question

	query := `SELECT * FROM diary.question q WHERE q.id=$1;`

	db.QueryRow(query, id).Scan(&q.Id, &q.Text, &q.ShownDate)

	return q
}

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
		&u.Subscribed,
		&u.Password,
	)

	if row == sql.ErrNoRows {
		fmt.Println("No question with that id")
	}

	return u
}

func GetUserByEmail(db *sql.DB, email string) (model.User, error) {
	var u model.User
	query := `SELECT * FROM diary.user u WHERE u.email=$1;`
	err := db.QueryRow(query, email).Scan(
		&u.Id,
		&u.CreatedAt,
		&u.Email,
		&u.Name,
		&u.Role,
		&u.SubId,
		&u.Subscribed,
		&u.Password,
	)

	return u, err
}

func CreateUser(db *sql.DB, email, password, name string) (model.User, error) {
	var u model.User
	query := `INSERT INTO diary.user (email, password, name) VALUES($1, $2, $3) RETURNING *;`

	err := db.QueryRow(query, email).Scan(
		&u.Id,
		&u.CreatedAt,
		&u.Email,
		&u.Name,
		&u.Role,
		&u.SubId,
		&u.Subscribed)
	return u, err
}

func GetQuestions(db *sql.DB) ([]model.Question, error) {
	var questions []model.Question
	query := `SELECT * FROM diary.question ORDER BY shown_date ASC;`

	rows, err := db.Query(query)

	if err != nil {
		return questions, err
	}

	for rows.Next() {
		var q model.Question
		var dateString string

		err = rows.Scan(&q.Id, &q.Text, &dateString)

		if err != nil {
			return questions, err
		}

		dateWithoutTime := strings.Split(dateString, "T")
		q.ShownDate = dateWithoutTime[0]

		questions = append(questions, q)
	}
	defer rows.Close()

	return questions, nil
}

func GetQuestionsLike(db *sql.DB, search string) ([]model.Question, error) {
	var questions []model.Question
	query := `SELECT * FROM diary.question WHERE text ILIKE $1 ORDER BY shown_date ASC;`

	rows, err := db.Query(query, "%"+search+"%")

	if err != nil {
		return questions, err
	}

	for rows.Next() {
		var q model.Question
		var dateString string

		err = rows.Scan(&q.Id, &q.Text, &dateString)

		if err != nil {
			return questions, err
		}

		dateWithoutTime := strings.Split(dateString, "T")
		q.ShownDate = dateWithoutTime[0]

		questions = append(questions, q)
	}
	defer rows.Close()

	return questions, nil
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
		UserId:      userId,
		QuestionId:  questionId,
		Text:        text,
		CreatedDate: createdDate,
	}
	query := "INSERT INTO diary.note (user_id, text, created_date, question_id) VALUES ($1, $2, $3, $4) returning id;"

	db.QueryRow(query, userId, text, createdDate, questionId).Scan(&note.Id)

	return note
}

func DeleteNote(db *sql.DB, id int) {
	db.Exec("DELETE FROM diary.note WHERE id=$1", id)
}

func GetNotesByText(db *sql.DB, userId int, search string) ([]model.Note, error) {
	var notes []model.Note
	query := `SELECT * FROM diary.note WHERE user_id=$1 AND text ILIKE $2;`

	rows, err := db.Query(query, userId, "%"+search+"%")

	if err != nil {
		return notes, err
	}

	for rows.Next() {
		var n model.Note

		rows.Scan(&n.Id, &n.UserId, &n.Text, &n.CreatedDate, &n.QuestionId)

		n.CreatedDate = strings.Split(n.CreatedDate, "T")[0]

		notes = append(notes, n)
	}
	defer rows.Close()

	return notes, nil
}

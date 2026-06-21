package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aaron-vasilev/diary/src/db"
	"github.com/aaron-vasilev/diary/src/model"
	"github.com/aaron-vasilev/diary/src/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetQuestion(id int) model.Question {
	q, err := db.Query.GetQuestion(context.Background(), int32(id))
	if err != nil {
		return model.Question{}
	}
	return questionToModel(q)
}

func GetQuestionByDate(date string) model.Question {
	q := model.Question{ShownDate: date}
	if !utils.DateStrIsValid(date) {
		return q
	}
	parts := strings.Split(date, "-")
	queryDate := "2024-" + parts[1] + "-" + parts[2]
	t, err := time.Parse("2006-01-02", queryDate)
	if err != nil {
		return q
	}
	result, err := db.Query.GetQuestionByDate(context.Background(), pgtype.Date{Time: t, Valid: true})
	if err != nil {
		return q
	}
	q.Id = int(result.ID)
	q.Text = result.Text
	return q
}

func GetUserById(id int) model.User {
	u, err := db.Query.GetUserById(context.Background(), int32(id))
	if err != nil {
		return model.User{}
	}
	return userToModel(u)
}

func GetUserByEmail(email string) (model.User, error) {
	u, err := db.Query.GetUserByEmail(context.Background(), email)
	if err != nil {
		return model.User{}, err
	}
	return userToModel(u), nil
}

func CreateUser(email, password, name string) (model.User, error) {
	u, err := db.Query.CreateUser(context.Background(), db.CreateUserParams{
		Email:    email,
		Password: pgtype.Text{String: password, Valid: true},
		Name:     name,
	})
	if err != nil {
		return model.User{}, err
	}
	return userToModel(u), nil
}

func GetQuestions() ([]model.Question, error) {
	rows, err := db.Query.GetQuestions(context.Background())
	if err != nil {
		return nil, err
	}
	questions := make([]model.Question, len(rows))
	for i, q := range rows {
		questions[i] = questionToModel(q)
	}
	return questions, nil
}

func GetQuestionsLike(search string) ([]model.Question, error) {
	rows, err := db.Query.GetQuestionsLike(context.Background(), "%"+search+"%")
	if err != nil {
		return nil, err
	}
	questions := make([]model.Question, len(rows))
	for i, q := range rows {
		questions[i] = questionToModel(q)
	}
	return questions, nil
}

func UpdateQuestion(id int, text string) model.Question {
	q, err := db.Query.UpdateQuestion(context.Background(), db.UpdateQuestionParams{
		Text: text,
		ID:   int32(id),
	})
	if err != nil {
		return model.Question{}
	}
	return questionToModel(q)
}

func GetNoteById(id int) model.Note {
	n, err := db.Query.GetNoteById(context.Background(), int32(id))
	if err != nil {
		return model.Note{}
	}
	return noteToModel(n)
}

func UpdateNote(id int, text string) model.Note {
	n, err := db.Query.UpdateNote(context.Background(), db.UpdateNoteParams{
		Text: text,
		ID:   int32(id),
	})
	if err != nil {
		return model.Note{}
	}
	return noteToModel(n)
}

func GetNotes(userId, questionId int) []model.Note {
	rows, err := db.Query.GetNotes(context.Background(), db.GetNotesParams{
		UserID:     userId,
		QuestionID: pgtype.Int4{Int32: int32(questionId), Valid: true},
	})
	if err != nil {
		return nil
	}
	notes := make([]model.Note, len(rows))
	for i, n := range rows {
		notes[i] = noteToModel(n)
	}
	return notes
}

func CreateNote(userId, questionId int, createdDate, text string) model.Note {
	t, err := time.Parse("2006-01-02", createdDate)
	if err != nil {
		return model.Note{UserId: userId, QuestionId: questionId, Text: text, CreatedDate: createdDate}
	}
	n, err := db.Query.CreateNote(context.Background(), db.CreateNoteParams{
		UserID:      userId,
		Text:        text,
		CreatedDate: t,
		QuestionID:  pgtype.Int4{Int32: int32(questionId), Valid: true},
	})
	if err != nil {
		return model.Note{UserId: userId, QuestionId: questionId, Text: text, CreatedDate: createdDate}
	}
	return noteToModel(n)
}

func DeleteNote(id int) {
	db.Query.DeleteNote(context.Background(), int32(id))
}

func GetNotesByText(userId int, search string) ([]model.Note, error) {
	rows, err := db.Query.GetNotesByText(context.Background(), db.GetNotesByTextParams{
		UserID: userId,
		Text:   "%" + search + "%",
	})
	if err != nil {
		return nil, err
	}
	notes := make([]model.Note, len(rows))
	for i, n := range rows {
		notes[i] = noteToModel(n)
	}
	return notes, nil
}

func questionToModel(q db.DiaryQuestion) model.Question {
	shownDate := ""
	if q.ShownDate.Valid {
		shownDate = q.ShownDate.Time.Format("2006-01-02")
	}
	return model.Question{
		Id:        int(q.ID),
		Text:      q.Text,
		ShownDate: shownDate,
	}
}

func noteToModel(n db.DiaryNote) model.Note {
	qid := 0
	if n.QuestionID.Valid {
		qid = int(n.QuestionID.Int32)
	}
	return model.Note{
		Id:          int(n.ID),
		UserId:      n.UserID,
		QuestionId:  qid,
		Text:        n.Text,
		CreatedDate: n.CreatedDate.Format("2006-01-02"),
	}
}

func userToModel(u db.DiaryUser) model.User {
	var subId *string
	if u.SubID.Valid {
		subId = &u.SubID.String
	}
	var password *string
	if u.Password.Valid {
		password = &u.Password.String
	}
	var telegramId *int64
	if u.TelegramID.Valid {
		telegramId = &u.TelegramID.Int64
	}
	return model.User{
		Id:         int(u.ID),
		CreatedAt:  u.CreatedAt.Format(time.RFC3339),
		Email:      u.Email,
		Name:       u.Name,
		Role:       model.Role(u.Role),
		SubId:      subId,
		Subscribed: u.Subscribed,
		Password:   password,
		TelegramId: telegramId,
	}
}

func UpsertTelegramUser(telegramId int64, name string) (model.User, error) {
	email := fmt.Sprintf("telegram_%d@telegram.placeholder", telegramId)
	u, err := db.Query.UpsertTelegramUser(context.Background(), db.UpsertTelegramUserParams{
		Email:      email,
		Name:       name,
		TelegramID: pgtype.Int8{Int64: telegramId, Valid: true},
	})
	if err != nil {
		return model.User{}, err
	}
	return userToModel(u), nil
}

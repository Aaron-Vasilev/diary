-- name: GetQuestion :one
SELECT * FROM diary.question WHERE id=$1;

-- name: GetQuestionByDate :one
SELECT * FROM diary.question WHERE shown_date=$1;

-- name: GetUserById :one
SELECT * FROM diary.user WHERE id=$1;

-- name: GetUserByEmail :one
SELECT * FROM diary.user WHERE email=$1;

-- name: CreateUser :one
INSERT INTO diary.user (email, password, name) VALUES($1, $2, $3) RETURNING *;

-- name: GetQuestions :many
SELECT * FROM diary.question ORDER BY shown_date ASC;

-- name: GetQuestionsLike :many
SELECT * FROM diary.question WHERE text ILIKE $1 ORDER BY shown_date ASC;

-- name: UpdateQuestion :one
UPDATE diary.question SET text=$1 WHERE id=$2 RETURNING *;

-- name: GetNoteById :one
SELECT * FROM diary.note WHERE id=$1;

-- name: UpdateNote :one
UPDATE diary.note SET text=$1 WHERE id=$2 RETURNING *;

-- name: GetNotes :many
SELECT * FROM diary.note WHERE user_id=$1 AND question_id=$2 ORDER BY id ASC;

-- name: CreateNote :one
INSERT INTO diary.note (user_id, text, created_date, question_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: DeleteNote :exec
DELETE FROM diary.note WHERE id=$1;

-- name: GetNotesByText :many
SELECT * FROM diary.note WHERE user_id=$1 AND text ILIKE $2;

-- name: GetUserByTelegramId :one
SELECT * FROM diary.user WHERE telegram_id=$1;

-- name: UpsertTelegramUser :one
INSERT INTO diary.user (email, name, telegram_id)
VALUES ($1, $2, $3)
ON CONFLICT (telegram_id) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

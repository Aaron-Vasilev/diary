CREATE SCHEMA IF NOT EXISTS diary;

CREATE TABLE diary.user (
    id          SERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    email       TEXT        NOT NULL UNIQUE,
    name        TEXT        NOT NULL,
    role        TEXT        NOT NULL DEFAULT 'user',
    sub_id      TEXT,
    subscribed  BOOLEAN     NOT NULL DEFAULT FALSE,
    password    TEXT,
    telegram_id BIGINT      UNIQUE
);

CREATE TABLE diary.question (
    id          SERIAL PRIMARY KEY,
    text        TEXT        NOT NULL,
    shown_date  DATE
);

CREATE TABLE diary.note (
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER     NOT NULL REFERENCES diary.user(id) ON DELETE CASCADE,
    text          TEXT        NOT NULL,
    created_date  DATE        NOT NULL DEFAULT CURRENT_DATE,
    question_id   INTEGER     REFERENCES diary.question(id) ON DELETE SET NULL
);

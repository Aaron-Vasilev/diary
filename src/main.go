package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/aaron-vasilev/diary-templ/src/auth"
	"github.com/aaron-vasilev/diary-templ/src/router"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
  loadEnv()
  PORT := os.Getenv("PORT")
  app := echo.New()
  db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

  if err != nil {
    log.Fatal("Error connecting to db", err)
  }
  defer db.Close()

  auth.NewAuth()
  app.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

  router.Connect(app, db)

  app.Static("src/styles", "src/styles")
  app.Static("public/", "public/")
  fmt.Printf("Server started at localhost%s\n", PORT)

  err = app.Start(PORT)

  log.Fatal("† line 33 err", err)
}

func loadEnv() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("No .env")
  }
}


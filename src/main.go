package main

import (
	"fmt"
  "log"
	"os"

	"database/sql"
	"github.com/aaron-vasilev/diary-templ/src/handler"
	"github.com/joho/godotenv"
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

  app.GET("/question-list", handler.HandlerCtx{
    Db: db,
  }.QuestionListHandler)

  app.GET("/diary", handler.HandlerCtx{
    Db: db,
  }.DiaryHandler)

  app.POST("/post", func(c echo.Context) error {
    return c.HTML(200, "<h1>POST</h1>")
  })

  app.Static("src/styles", "src/styles")
  fmt.Printf("Server started at localhost%s\n", PORT)

  err = app.Start(PORT)

  fmt.Println("† line 33 err", err)
}

func withUser(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {
//  r  c.Set
    return next(c)
  }
}

func loadEnv() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("No .env")
  }
}


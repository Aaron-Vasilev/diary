package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
  "diary/src/controller"
)

func main() {
  loadEnv()
  setPaths()
  PORT := os.Getenv("PORT")
  db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

  if err != nil {
    log.Fatal("Error connecting to db", err)
  }
  defer db.Close()

  question := controller.GetQuestionByDate(db, "2024-03-01")

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    components, err := template.ParseFiles(
      "./src/components/layout.html",
      "./src/pages/diary.html",
      "./src/components/question.html",
    )

    if err != nil {
      http.Error(w, "Error loading main template", http.StatusInternalServerError)
      return
    }

    tmpl := template.Must(components.Clone())

    err = tmpl.Execute(w, question)
    if err != nil {
      http.Error(w, "Error executing template", http.StatusInternalServerError)
      return
    }
  })

  fmt.Printf("Server started at localhost%s\n", PORT)

  err = http.ListenAndServe(PORT, nil)
  if err != nil  {
    log.Fatal("Error while starting server", err)
  }
}

func loadEnv() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("No .env")
  }
}

func setPaths() {
  http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./src/styles/"))))
}

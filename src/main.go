package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/aaron-vasilev/diary/src/auth"
	"github.com/aaron-vasilev/diary/src/router"
	"github.com/aaron-vasilev/diary/src/utils"
	"github.com/akrylysov/algnhsa"
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

	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting to db", err)
	}
	defer db.Close()

	auth.NewAuth()
	app.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	router.ConnectRoutes(app, db)

	if utils.IsProd() {
		app.Use(middleware)
		algnhsa.ListenAndServe(app, nil)
	} else {
		app.Static("public/", "public/")
		fmt.Printf("Server started at localhost%s\n", PORT)

		err = app.Start(PORT)
		log.Fatal("† line 33 err", err)
	}
}

func loadEnv() {
	env := os.Getenv("ENV")

	if env == "" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("No .env")
		}
	}
}

func middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
		return next(c)
	}
}

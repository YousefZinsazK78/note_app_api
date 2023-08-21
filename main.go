package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/yousefzinsazk78/note_app_api/api"
	notedb "github.com/yousefzinsazk78/note_app_api/database/note_db"
	sessiontbl "github.com/yousefzinsazk78/note_app_api/database/session_tbl"
	usertbl "github.com/yousefzinsazk78/note_app_api/database/user_tbl"
)

const (
	username = "root"
	password = "13781378"
	hostname = "127.0.0.1"
	dbname   = "noteappdb"
	port     = ":5000"
)

//final task
//todo : user authorization
//todo : validation for input forms and jsons and ...

func main() {
	db, err := sql.Open("mysql", generateDsn())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var (
		tmplEngine = html.New("./views", ".html")
		app        = fiber.New(
			fiber.Config{
				ErrorHandler: api.ErrorHandler,
				Views:        tmplEngine,
			},
		)
		mysqlNoteStorer    = notedb.NewMysqlNoteStorer(db)
		mysqlUserStorer    = usertbl.NewMysqlUserStorer(db)
		mysqlSessionStorer = sessiontbl.NewMysqlSessionStorer(db)
		apiV1              = api.NewApi(mysqlNoteStorer, mysqlUserStorer, mysqlSessionStorer)
		v1                 = app.Group("/api", apiV1.AuthMiddleware)
	)

	app.Static("/static", "./views")

	//sign-in and logout
	app.Get("/welcome", apiV1.HandleWelcome)
	app.Post("/signin", apiV1.HandleSignIn)
	app.Get("/signin", apiV1.HandleSignInGet)
	app.Post("/signup", apiV1.HandleSignUp)
	app.Get("/signup", apiV1.HandleSignUpGet)
	app.Get("/refresh", apiV1.HandleRefresh)
	app.Get("/logout", apiV1.HandleLogout)

	//html template version
	app.Get("/", apiV1.HandleIndex)
	app.Get("/create", apiV1.HandleCreate)
	app.Post("/create", apiV1.HandleCreatePost)
	app.Post("/", apiV1.HandleDeletePost)
	app.Get("/edit", apiV1.HandleEdit)
	app.Post("/edit", apiV1.HandleEditPost)

	//json version 1
	v1.Post("/v1/notes", apiV1.HandleCreateNote)
	v1.Get("/v1/notes", apiV1.HandleNotes)
	v1.Get("/v1/notes/:id", apiV1.HandleNoteID)
	v1.Put("/v1/notes/:id", apiV1.HandleUpdateNote)
	v1.Delete("/v1/notes/:id/delete", apiV1.HandleDeleteNote)
	log.Fatal(app.Listen(port))

}

func generateDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, hostname, dbname)
}

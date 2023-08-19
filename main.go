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
		api                = api.NewApi(mysqlNoteStorer, mysqlUserStorer, mysqlSessionStorer)
		v1                 = app.Group("/api")
	)

	app.Static("/static", "./views")

	//sign-in and logout
	app.Get("/welcome", api.HandleWelcome)
	app.Post("/signin", api.HandleSignIn)
	app.Post("/signup", api.HandleSignUp)
	app.Get("/refresh", api.HandleRefresh)

	//html template version
	app.Get("/", api.HandleIndex)
	app.Get("/create", api.HandleCreate)
	app.Post("/create", api.HandleCreatePost)
	app.Post("/", api.HandleDeletePost)
	app.Get("/edit", api.HandleEdit)
	app.Post("/edit", api.HandleEditPost)

	//json version 1
	v1.Post("/v1/notes", api.HandleCreateNote)
	v1.Get("/v1/notes", api.HandleNotes)
	v1.Get("/v1/notes/:id", api.HandleNoteID)
	v1.Put("/v1/notes/:id", api.HandleUpdateNote)
	v1.Delete("/v1/notes/:id/delete", api.HandleDeleteNote)
	log.Fatal(app.Listen(port))

}

func generateDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, hostname, dbname)
}

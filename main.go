package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/yousefzinsazk78/note_app_api/api"
	notedb "github.com/yousefzinsazk78/note_app_api/database/note_db"
)

const (
	username = "root"
	password = "13781378"
	hostname = "127.0.0.1"
	dbname   = "noteappdb"
	port     = ":5000"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {

	db, err := sql.Open("mysql", generateDsn())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var (
		app = fiber.New(
			config,
		)
		mysqlNoteStorer = notedb.NewMysqlNoteStorer(db)
		api             = api.NewApi(mysqlNoteStorer)
		v1              = app.Group("/api")
	)

	app.Static("/", "./static")

	v1.Post("/v1/notes", api.HandleCreateNote)
	v1.Get("/v1/notes", api.HandleNotes)
	v1.Get("/v1/notes/:id", api.HandleNoteID)
	v1.Put("/v1/notes/:id", api.HandleUpdateNote)
	v1.Delete("/v1/notes/:id/delete", api.HandleDeleteNote)
	log.Fatal(app.Listen(port))
}

func generateDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

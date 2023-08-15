package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/yousefzinsazk78/note_app_api/api"
	notedb "github.com/yousefzinsazk78/note_app_api/database/note_db"
	"github.com/yousefzinsazk78/note_app_api/types"
)

const (
	username = "root"
	passowrd = "13781378"
	hostname = "127.0.0.1"
	dbname   = "noteappdb"
)

var createtblquery = `CREATE TABLE IF NOT EXISTS note_table(id int primary key auto_increment, title text, description text);`
var insertvalueintotbl = `INSERT INTO note_table(title, description) VALUES (?, ?);`
var readvaluefromtbl = `SELECT * FROM note_table;`
var updatevaluetotbl = `UPDATE note_table SET title=?, description=? WHERE id=?;`
var deletevaluefromtbl = `DELETE FROM note_table WHERE ID=?`

// TODO: make first json api end point
// ? create note / update note / read note / delete note
func main() {

	db, err := sql.Open("mysql", generateDsn())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var (
		app             = fiber.New()
		mysqlNoteStorer = notedb.NewMysqlNoteStorer(db)
		api             = api.NewApi(mysqlNoteStorer)
		v1              = app.Group("/api")
	)

	v1.Post("/v1/notes", api.HandleCreateNote)
	v1.Get("/v1/notes", api.HandleNotes)
	v1.Get("/v1/notes/:id", api.HandleNoteID)
	v1.Put("/v1/notes/:id", api.HandleUpdateNote)
	v1.Delete("/v1/notes/:id/delete", api.HandleDeleteNote)

	app.Listen(":5000")
}

func deleteValueFromTbl(db *sql.DB, note types.Note) {
	stmt, err := db.Prepare(deletevaluefromtbl)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(note.ID)
	if err != nil {
		log.Fatal(err)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(rows, "deleted successfully!")
}

func updateValueToNoteTable(db *sql.DB, note types.Note) {
	stmt, err := db.Prepare(updatevaluetotbl)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(note.Title, note.Description, note.ID)
	if err != nil {
		log.Fatal(err)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(rows, "updated successfully!")
}

func readValueFromNoteTable(db *sql.DB) {
	res, err := db.Query(readvaluefromtbl)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer res.Close()

	for res.Next() {
		var note types.Note
		err := res.Scan(&note.ID, &note.Title, &note.Description)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", note)
	}
}

func insertValuesIntoNoteTable(db *sql.DB, note types.Note) {
	stmt, err := db.Prepare(insertvalueintotbl)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(note.Title, note.Description)
	if err != nil {
		log.Fatal(err)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("%d note inserted successfully...", rows)
}

func createNoteTable(db *sql.DB) {
	//create table
	res, err := db.Exec(createtblquery)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(rows)
}

func generateDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, passowrd, hostname, dbname)
}

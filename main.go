package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	passowrd = "13781378"
	hostname = "127.0.0.1"
	dbname   = "noteappdb"
)

// TODO: make connection to database
// TODO: make first json api end point
// ? create note / update note / read note / delete note
func main() {
	db, err := sql.Open("mysql", generateDsn())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("database successfully connected!")
}

func generateDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, passowrd, hostname, dbname)
}

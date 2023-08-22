package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	usertbl "github.com/yousefzinsazk78/note_app_api/database/user_tbl"
	"github.com/yousefzinsazk78/note_app_api/types"
)

const (
	username = "root"
	password = "13781378"
	hostname = "127.0.0.1"
	dbname   = "noteappdb"
	port     = ":5000"
)

func CreateAdminUser(userStore usertbl.UserStorer, username, password string, isAdmin bool) error {
	var user = types.User{
		Username: username,
		Password: password,
		IsAdmin:  isAdmin,
	}
	if err := user.ValidateUser(); err != nil {
		log.Fatal(err)
	}

	if err := userStore.InsertUser(user); err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	db, err := sql.Open("mysql", generateDsn())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var userStore = usertbl.NewMysqlUserStorer(db)

	if err := CreateAdminUser(userStore, "yousefAdmin", "yousefAdmin1234", true); err != nil {
		log.Fatal(err)
	}
}

func generateDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, hostname, dbname)
}

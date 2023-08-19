package usertbl

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/yousefzinsazk78/note_app_api/types"
)

type UserStorer interface {
	InsertUser(types.User) error
	DeleteUser(int) error
	GetUserByUsername(string) (*types.User, error)
}

type MysqlUserStorer struct {
	db *sql.DB
}

func NewMysqlUserStorer(db *sql.DB) *MysqlUserStorer {
	return &MysqlUserStorer{
		db: db,
	}
}

func (m *MysqlUserStorer) InsertUser(user types.User) error {
	insertQuery := `INSERT INTO user_tbl(Username, Password) VALUES (?, ?);`

	stmt, err := m.db.Prepare(insertQuery)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (m *MysqlUserStorer) DeleteUser(id int) error {
	deleteQuery := `DELETE FROM user_tbl WHERE ID=?;`

	stmt, err := m.db.Prepare(deleteQuery)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (m *MysqlUserStorer) GetUserByUsername(username string) (*types.User, error) {
	selectQuery := `SELECT * FROM user_tbl WHERE Username=?;`
	res, err := m.db.Query(selectQuery, username)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var user types.User
	for res.Next() {
		err := res.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%s : not found!", username)
			}
			return nil, err
		}
	}
	return &user, nil
}

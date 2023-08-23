package sessiontbl

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/yousefzinsazk78/note_app_api/types"
)

type SessionStorer interface {
	InsertSession(string, bool, time.Time, string) error
	DeleteSession(string) error
	GetSession(string) (*types.Session, error)
}

type MysqlSessionStorer struct {
	db *sql.DB
}

func NewMysqlSessionStorer(db *sql.DB) *MysqlSessionStorer {
	return &MysqlSessionStorer{
		db: db,
	}
}

func (m *MysqlSessionStorer) InsertSession(username string, isAdmin bool, expireTime time.Time, sessionToken string) error {
	insertSession := `INSERT INTO session_tbl(Username, IsAdmin, SessionExpiry, SessionToken) VALUES (?, ?, ?, ?);`

	stmt, err := m.db.Prepare(insertSession)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(username, isAdmin, expireTime, sessionToken)
	if err != nil {
		return err
	}
	return nil
}

func (m *MysqlSessionStorer) DeleteSession(sessionToken string) error {
	deleteQuery := `DELETE FROM session_tbl WHERE SessionToken=?;`

	stmt, err := m.db.Prepare(deleteQuery)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sessionToken)
	if err != nil {
		return err
	}
	return nil
}

func (m *MysqlSessionStorer) GetSession(SessionToken string) (*types.Session, error) {
	selectQuery := fmt.Sprintf(`SELECT * FROM session_tbl WHERE SessionToken = '%s';`, SessionToken)
	res, err := m.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var session types.Session
	for res.Next() {
		err := res.Scan(&session.Username, &session.IsAdmin, &session.SessionExpiry, &session.SessionToken)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%s : not found!", SessionToken)
			}
			return nil, err
		}

	}
	return &session, nil
}

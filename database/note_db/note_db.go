package notedb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/yousefzinsazk78/note_app_api/types"
)

type NoteStorer interface {
	InsertNote(*types.Note) (int, error)
	UpdateNote(*types.Note, int) (int, error)
	DeleteNote(int) (int, error)
	GetNoteByID(int) (*types.Note, error)
	GetNotes() ([]types.Note, error)
}

type MySqlNoteStorer struct {
	db *sql.DB
}

func NewMysqlNoteStorer(db *sql.DB) *MySqlNoteStorer {
	return &MySqlNoteStorer{
		db: db,
	}
}

func (m *MySqlNoteStorer) InsertNote(note *types.Note) (int, error) {
	const insertvalueintotbl = `INSERT INTO note_table(title, description) VALUES (?, ?);`

	stmt, err := m.db.Prepare(insertvalueintotbl)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(note.Title, note.Description)
	if err != nil {
		return 0, nil
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	log.Printf("%d note inserted successfully...", rows)
	return int(rows), nil
}

func (m *MySqlNoteStorer) UpdateNote(note *types.Note, id int) (int, error) {
	const updatevaluetotbl = `UPDATE note_table SET title=?, description=? WHERE id=?;`
	stmt, err := m.db.Prepare(updatevaluetotbl)
	if err != nil {
		return 0, nil
	}
	defer stmt.Close()
	res, err := stmt.Exec(note.Title, note.Description, id)
	if err != nil {
		return 0, nil
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	log.Println(rows, "updated successfully!")
	return int(rows), nil
}

func (m *MySqlNoteStorer) DeleteNote(id int) (int, error) {
	const deletevaluefromtbl = `DELETE FROM note_table WHERE ID=?`
	stmt, err := m.db.Prepare(deletevaluefromtbl)
	if err != nil {
		return 0, nil
	}
	defer stmt.Close()
	res, err := stmt.Exec(id)
	if err != nil {
		return 0, nil
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	log.Println(rows, "deleted successfully!")
	return int(rows), nil
}

func (m *MySqlNoteStorer) GetNotes() ([]types.Note, error) {
	const readvaluefromtbl = `SELECT * FROM note_table;`
	res, err := m.db.Query(readvaluefromtbl)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var notes []types.Note

	for res.Next() {
		var note types.Note
		err := res.Scan(&note.ID, &note.Title, &note.Description)
		if err != nil {
			log.Fatal(err)

		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (m *MySqlNoteStorer) GetNoteByID(id int) (*types.Note, error) {
	const readvaluefromtbl = `SELECT * FROM note_table WHERE id=?;`
	res, err := m.db.Query(readvaluefromtbl, id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var note types.Note
	for res.Next() {
		err := res.Scan(&note.ID, &note.Title, &note.Description)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%d : note not found!", id)
			}
			return nil, err
		}

	}
	return &note, nil
}

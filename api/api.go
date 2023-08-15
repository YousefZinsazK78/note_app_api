package api

import notedb "github.com/yousefzinsazk78/note_app_api/database/note_db"

type Api struct {
	NoteStorer notedb.NoteStorer
}

func NewApi(notestorer notedb.NoteStorer) *Api {
	return &Api{
		NoteStorer: notestorer,
	}
}

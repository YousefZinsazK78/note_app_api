package api

import (
	notedb "github.com/yousefzinsazk78/note_app_api/database/note_db"
	sessiontbl "github.com/yousefzinsazk78/note_app_api/database/session_tbl"
	usertbl "github.com/yousefzinsazk78/note_app_api/database/user_tbl"
)

type Api struct {
	NoteStorer    notedb.NoteStorer
	UserStorer    usertbl.UserStorer
	SessionStorer sessiontbl.SessionStorer
}

func NewApi(notestorer notedb.NoteStorer, userstorer usertbl.UserStorer, sessionstorer sessiontbl.SessionStorer) *Api {
	return &Api{
		NoteStorer:    notestorer,
		UserStorer:    userstorer,
		SessionStorer: sessionstorer,
	}
}

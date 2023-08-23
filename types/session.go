package types

import "time"

type Session struct {
	Username      string
	IsAdmin       bool
	SessionExpiry time.Time
	SessionToken  string
}

func (s Session) IsExpired() bool {
	return s.SessionExpiry.Before(time.Now())
}

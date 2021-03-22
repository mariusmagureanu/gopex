package ds

import "time"

type Session struct {
	SessionKey  string    `json:"session_key"`
	SessionData string    `json:"session_data"`
	ExpireDate  time.Time `json:"expire_date"`
}

func (Session) TableName() string {
	return "django_session"
}

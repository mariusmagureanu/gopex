package dbl

import (
	"gorm.io/gorm"

	"github.com/mariusmagureanu/gopex/pkg/ds"
)

type SessionDao interface {
	GetByKey(*ds.Session, string) error
}

type sessionDao struct {
	db *gorm.DB
}

func (s sessionDao) GetByKey(sess *ds.Session, sessionKey string) error {
	return s.db.Where("session_key=?", sessionKey).First(sess).Error
}

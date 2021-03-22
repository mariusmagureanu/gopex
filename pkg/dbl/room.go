package dbl

import (
	"gorm.io/gorm"

	"bitbucket.org/kinlydev/gopex/pkg/ds"
)

type RoomDao interface {
	GetByID(*ds.Room, int) error
}

type roomDao struct {
	db *gorm.DB
}

func (r roomDao) GetByID(room *ds.Room, roomID int) error {
	return r.db.Where("id=?", roomID).First(room).Error
}

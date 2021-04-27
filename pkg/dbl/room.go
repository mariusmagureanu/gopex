package dbl

import (
	"gorm.io/gorm"

	"github.com/mariusmagureanu/gopex/pkg/ds"
)

// RoomDao is in interface which exhibits
// CRUD operations for the Room type.
type RoomDao interface {
	GetByID(*ds.Room, int) error
	GetAll(*[]ds.Room) error
	Create(room *ds.Room) error
}

type roomDao struct {
	db *gorm.DB
}

func (r roomDao) GetByID(room *ds.Room, roomID int) error {
	return r.db.Where("id=?", roomID).First(room).Error
}

func (r roomDao) GetAll(rooms *[]ds.Room) error {
	return r.db.Find(rooms).Error
}

func (r roomDao) Create(room *ds.Room) error {
	return r.db.Create(room).Error
}

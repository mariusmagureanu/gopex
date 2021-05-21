package dbl

import (
	"errors"

	"gorm.io/gorm"

	"github.com/mariusmagureanu/gopex/pkg/ds"
	e "github.com/mariusmagureanu/gopex/pkg/errors"
)

// RoomDao is in interface which exhibits
// CRUD operations for the Room type.
type RoomDao interface {
	GetByID(*ds.Room, uint) error
	GetByName(*ds.Room, string) error
	GetAll(*[]ds.Room) error
	Create(*ds.Room) error
	Save(*ds.Room) error
	Delete(*ds.Room) error
}

type roomDao struct {
	db *gorm.DB
}

func (r roomDao) GetByID(room *ds.Room, roomID uint) error {
	err := r.db.Where("id=?", roomID).First(room).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return e.ErrRecordNotFound
	}

	return err
}

func (r roomDao) GetByName(room *ds.Room, name string) error {
	err := r.db.Where("name=?", name).First(room).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return e.ErrRecordNotFound
	}

	return err
}

func (r roomDao) GetAll(rooms *[]ds.Room) error {
	return r.db.Find(rooms).Error
}

func (r roomDao) Create(room *ds.Room) error {
	return r.db.Create(room).Error
}

func (r roomDao) Save(room *ds.Room) error {
	return r.db.Save(room).Error
}

func (r roomDao) Delete(room *ds.Room) error {
	return r.db.Delete(room).Error
}

// Package dbl handles the persistence of
// various objects in the pexip monitor.
package dbl

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DAO struct {
	dbSession *gorm.DB

	roomDao    RoomDao
}

func (d *DAO) InitSqlite(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	d.initDataAccessObjects(db)

	return nil
}

func (d *DAO) InitPostgres(dsn string, maxIdle int, maxOpen int, maxLifetime time.Duration) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(maxLifetime)

	d.initDataAccessObjects(db)

	return sqlDB.Ping()
}

func (d *DAO) Rooms() RoomDao {
	return d.roomDao
}

func (d *DAO) initDataAccessObjects(gdb *gorm.DB) {
	d.dbSession = gdb

	d.roomDao = roomDao{db: gdb}
}

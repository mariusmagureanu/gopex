// Package dbl handles the persistence of
// various objects in the pexip monitor.
package dbl

import (
	"time"

	"github.com/mariusmagureanu/gopex/pkg/ds"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DAO is a type that manages the connection to a database
// and holds objects which provide CRUD functionality
// for all involved models.
// The entire implementation is based on
// https://gorm.io/docs/v2_release_note.html
type DAO struct {
	dbSession *gorm.DB

	roomDao RoomDao
}

// InitSqlite initializes a connection against
// a sqlite database.
func (d *DAO) InitSqlite(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	d.initDataAccessObjects(db)

	return nil
}

// InitPostgres initializes a connection against
// a Postgres server.
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

// DropTables instructs gorm to drop all
// tables for the configured database.
func (d *DAO) DropTables() error {
	return d.dbSession.Migrator().DropTable(ds.Room{})
}

// CreateTables instructs gorm to create
// tables for all the corresponding model types.
func (d *DAO) CreateTables() error {
	return d.dbSession.Migrator().AutoMigrate(ds.Room{})
}

// Rooms returns a dao object which exposes
// CRUD functionality for the Room type.
func (d *DAO) Rooms() RoomDao {
	return d.roomDao
}

func (d *DAO) initDataAccessObjects(gdb *gorm.DB) {
	d.dbSession = gdb

	d.roomDao = roomDao{db: gdb}
}

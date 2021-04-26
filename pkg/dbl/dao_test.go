package dbl

import (
	"testing"
)

func TestDAOInit(t *testing.T) {
	dao, err := tearUp()

	if err != nil {
		t.Error(err)
	}

	err = tearDown(dao)

	if err != nil {
		t.Error(err)
	}
}

func tearUp() (DAO, error) {
	dao := DAO{}

	err := dao.InitSqlite("db.sqlite3")

	if err != nil {
		return dao, err
	}

	err = dao.CreateTables()

	return dao, err
}

func tearDown(dao DAO) error {
	return dao.DropTables()
}

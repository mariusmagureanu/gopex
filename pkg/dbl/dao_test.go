package dbl

import (
	"testing"
)

func TestDAO_Init(t *testing.T) {
	dao := DAO{}

	err := dao.InitSqlite("db.sqlite3")

	if err != nil {
		t.Fatal(err)
	}

}

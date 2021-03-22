package dbl

import (
	"testing"

	"bitbucket.org/kinlydev/gopex/pkg/ds"
)

func TestGetSessionByKey(t *testing.T) {
	d := DAO{}

	testSessKey := "mvfuzx78098qre4pv634wzd3s7amasnz"

	err := d.InitSqlite("db.sqlite3")

	if err != nil {
		t.Error(err)
	}

	var s ds.Session

	err = d.Sessions().GetByKey(&s, testSessKey)

	if err != nil {
		t.Error(err)
	}

	if s.SessionKey != testSessKey {
		t.Errorf("Was expecting %s, got %s instead", testSessKey, s.SessionKey)
	}
}

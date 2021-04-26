package dbl

import (
	"fmt"
	"testing"

	"github.com/mariusmagureanu/gopex/pkg/ds"
	"github.com/stretchr/testify/assert"
)

func TestRoomDaoCreate(t *testing.T) {
	dao, err := tearUp()

	if err != nil {
		t.Error(err)
	}

	defer func(dao DAO) {
		err := tearDown(dao)
		if err != nil {
			t.Error(err)
		}
	}(dao)

	var room ds.Room
	room.AliasID = 5
	room.AllowGuests = true
	room.HostPin = "1923"
	room.Locked = false
	room.Name = "foo@bar.com"
	err = dao.Rooms().Create(&room)

	if err != nil {
		t.Fatal(err)
	}

	var savedRoom ds.Room
	err = dao.Rooms().GetByID(&savedRoom, room.ID)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, room.AliasID, savedRoom.AliasID)
	assert.Equal(t, room.Name, savedRoom.Name)
	assert.Equal(t, room.Locked, savedRoom.Locked)
	assert.Equal(t, room.HostPin, savedRoom.HostPin)
}

func TestRoomDaoGetAll(t *testing.T) {
	dao, err := tearUp()

	if err != nil {
		t.Error(err)
	}

	defer func(dao DAO) {
		err := tearDown(dao)
		if err != nil {
			t.Error(err)
		}
	}(dao)

	for i := 0; i < 50; i++ {
		var room ds.Room
		room.AliasID = i
		room.AllowGuests = i%2 == 0
		room.HostPin = fmt.Sprintf("192%d", i)
		room.Locked = i%2 != 0
		room.Name = fmt.Sprintf("foo-%d@email.com", i)
		err = dao.Rooms().Create(&room)

		if err != nil {
			t.Fatal(err)
		}
	}

	var rooms []ds.Room

	err = dao.Rooms().GetAll(&rooms)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 50, len(rooms))
}

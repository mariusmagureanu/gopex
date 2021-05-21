package mux

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mariusmagureanu/gopex/pkg/ds"

	"github.com/mariusmagureanu/gopex/pexip"
	"github.com/mariusmagureanu/gopex/pkg/errors"
	logger "github.com/mariusmagureanu/gopex/pkg/log"

	"github.com/gorilla/mux"
)

func createNewRoomHandler(w http.ResponseWriter, r *http.Request) {
	var room ds.Room

	mimeType := r.Header.Get("Content-Type")

	if mimeType != "application/json" {
		logger.Warning("invalid content type", mimeType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	err = json.Unmarshal(b, &room)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = dao.Rooms().Create(&room)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	location := fmt.Sprintf("%s/%s/%s", apiV1Prefix, "room", room.Name)
	w.Header().Set("Content-Location", location)
	w.Header().Set("Content-Length", "0")

	w.WriteHeader(http.StatusCreated)
}

func getAllRoomsHandler(w http.ResponseWriter, r *http.Request) {
	var rooms []ds.Room

	outputType := r.Header.Get("Accept")

	if outputType != "application/json" {
		logger.Error("cannot serve requested mime type", outputType)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err := dao.Rooms().GetAll(&rooms)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(&rooms)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(out)
}

func getRoomHandler(w http.ResponseWriter, r *http.Request) {
	outputType := r.Header.Get("Accept")

	if outputType != "application/json" {
		logger.Error("cannot serve requested mime type", outputType)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	var room ds.Room
	vars := mux.Vars(r)

	confName := vars["room"]
	err := dao.Rooms().GetByName(&room, confName)

	if err != nil {
		switch err {
		case errors.ErrRecordNotFound:
			logger.Warning(err)
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	out, err := json.Marshal(&room)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(out)
}

func deleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	var room ds.Room
	vars := mux.Vars(r)

	confName := vars["room"]

	err := dao.Rooms().GetByName(&room, confName)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = dao.Rooms().Delete(&room)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func monitorStartHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {

	go func() {
		err := sseManager.Listen(conf.Name, token)
		if err != nil {
			logger.Error(err)
			return
		}
	}()

	return []byte{}, nil
}

func monitorStopHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {
	var (
		err error
		out []byte
	)

	err = sseManager.Stop(conf.Name)
	if err != nil {
		return out, err
	}

	err = tokenStore.Release(conf)
	return out, err
}

func conferenceLockHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {
	return conf.Lock(token)
}

func conferenceUnLockHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {
	return conf.Unlock(token)
}

func conferenceMuteGuestsHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {
	return conf.MuteGuests(token)
}

func conferenceUnmuteGuestsHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {
	return conf.UnmuteGuests(token)
}

func conferenceStatusHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {
	return conf.Status(token)
}

func conferenceDisconnectHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {
	return conf.Disconnect(token)
}

func conferenceDialHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {

	var (
		err error
		out []byte
	)

	dp, err := ioutil.ReadAll(r)

	if err != nil {
		return out, err
	}

	out, err = conf.Dial(token, dp)

	if err != nil {
		return out, err
	}

	var dpr pexip.DialResponse

	err = json.Unmarshal(out, &dpr)

	if err != nil {
		return out, err
	}

	participantStore.AddMultiple(dpr.Result)

	return out, nil
}

func conferenceParticipantsHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {
	return conf.Participants(token)
}

func conferenceMessageHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {
	var out []byte

	msg, err := ioutil.ReadAll(r)

	if err != nil {
		return out, err
	}

	return conf.Message(token, msg)
}

func transformLayoutHandler(conf *pexip.Conference, r io.Reader, token string) ([]byte, error) {
	var out []byte

	layout, err := ioutil.ReadAll(r)

	if err != nil {
		return out, err
	}

	return conf.Message(token, layout)
}

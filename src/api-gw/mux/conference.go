package mux

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mariusmagureanu/gopex/pexip"
	"github.com/mariusmagureanu/gopex/pkg/errors"
	logger "github.com/mariusmagureanu/gopex/pkg/log"

	"github.com/gorilla/mux"
)

func monitorStartHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug(r.Method, r.Proto, r.Host+r.RequestURI)
	vars := mux.Vars(r)
	confName := vars["room"]

	conf, err := confStore.Get(confName)

	if err != nil {
		logger.Warning(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = tokenStore.Watch(conf)

	if err != nil {
		if err == errors.ErrorRoomAlreadyStarted {
			logger.Info(err)
			w.WriteHeader(http.StatusOK)
			return
		}

		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func monitorStopHandler(w http.ResponseWriter, r *http.Request, conf *pexip.Conference, token string) {

	err := tokenStore.Release(conf)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func conferenceCmdHandler(w http.ResponseWriter, r *http.Request, conf *pexip.Conference, token string) {
	var (
		err     error
		cmdResp []byte
	)

	vars := mux.Vars(r)
	cmd := vars["cmd"]

	switch cmd {
	case pexip.CommandLock:
		cmdResp, err = conf.Lock(token)
		break
	case pexip.CommandUnlock:
		cmdResp, err = conf.Unlock(token)
		break
	case pexip.CommandMuteGuests:
		cmdResp, err = conf.MuteGuests(token)
		break
	case pexip.CommandUnmuteGuests:
		cmdResp, err = conf.UnmuteGuests(token)
		break
	default:
		logger.Warning("unsupported command", cmd)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(cmdResp)
}

func conferenceStatusHandler(w http.ResponseWriter, r *http.Request, conf *pexip.Conference, token string) {

	statusResp, err := conf.Status(token)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(statusResp)
}

func conferenceDisconnectHandler(w http.ResponseWriter, r *http.Request, conf *pexip.Conference, token string) {

	disconnectResp, err := conf.Disconnect(token)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(disconnectResp)
}

func conferenceDialHandler(w http.ResponseWriter, r *http.Request, conf *pexip.Conference, token string) {

	dp, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	dialResp, err := conf.Dial(token, dp)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var dpr pexip.DialResponse

	err = json.Unmarshal(dialResp, &dpr)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	participantStore.AddMultiple(dpr.Result)

	w.WriteHeader(http.StatusOK)
	w.Write(dialResp)
}

func conferenceParticipantsHandler(w http.ResponseWriter, r *http.Request, conf *pexip.Conference, token string) {

	participantsResp, err := conf.Participants(token)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(participantsResp)
}

func conferenceMessageHandler(w http.ResponseWriter, r *http.Request, conf *pexip.Conference, token string) {

	msg, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	messageResp, err := conf.Message(token, msg)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(messageResp)
}

func transformLayoutHandler(w http.ResponseWriter, r *http.Request, conf *pexip.Conference, token string) {

	layout, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	layoutResp, err := conf.Message(token, layout)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(layoutResp)
}

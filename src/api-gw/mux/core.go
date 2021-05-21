// Package mux holds the entire implementation
// of all the request handlers that are to be
// supported by the exposed rest interface.
package mux

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/mariusmagureanu/gopex/pkg/errors"

	"github.com/gorilla/mux"
	"github.com/mariusmagureanu/gopex/pkg/dbl"
	"github.com/mariusmagureanu/gopex/pkg/ds"

	"github.com/mariusmagureanu/gopex/pexip"
	logger "github.com/mariusmagureanu/gopex/pkg/log"
)

const (
	apiV1Prefix = "/api/v1"
)

var (
	tokenStore       = pexip.InitTokenStore()
	participantStore = pexip.InitParticipantStore()
	sseManager       = pexip.InitSSEManager()

	dao *dbl.DAO

	pexipNameSpace = fmt.Sprintf("%s/%s", apiV1Prefix, "pexip")

)

type confWrappedFunc = func(*pexip.Conference, io.Reader, string) ([]byte, error)
type partWrappedFunc = func(*pexip.Participant, string, string) ([]byte, error)

// InitMux initializes a new router for
// the rest api webserver.
func InitMux(db *dbl.DAO) (*mux.Router, error) {

	// This is really meh, should not be here!
	dao = db

	mgmtMux := mux.NewRouter()

	mgmtMux.HandleFunc(apiV1Prefix+"/ping", wrapRequestHandler(pingReqHandler)).Methods(http.MethodHead, http.MethodGet)

	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/participants/{part_uuid}/{cmd}", wrapRequestHandler(pingReqHandler))
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/override_layout", wrapRequestHandler(pingReqHandler))
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/participants/{part_uuid}/transfer", wrapRequestHandler(pingReqHandler))
	mgmtMux.HandleFunc(pexipNameSpace+"/room/{room}/participants/{part_uuid}/role", wrapRequestHandler(pingReqHandler))

	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/start", wrapConferenceHandler(monitorStartHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/stop", wrapConferenceHandler(monitorStopHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/participants", wrapConferenceHandler(conferenceParticipantsHandler)).Methods(http.MethodGet)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/dial", wrapConferenceHandler(conferenceDialHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/lock", wrapConferenceHandler(conferenceLockHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/unlock", wrapConferenceHandler(conferenceUnLockHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/muteguests", wrapConferenceHandler(conferenceMuteGuestsHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/unmuteguests", wrapConferenceHandler(conferenceUnmuteGuestsHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/transform_layout", wrapConferenceHandler(transformLayoutHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/disconnect", wrapConferenceHandler(conferenceDisconnectHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/participants/{part_uuid}/spotlighton", wrapParticipantHandler(participantSpotlightOnHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/participants/{part_uuid}/spotlightoff", wrapParticipantHandler(participantSpotlightOffHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/status", wrapConferenceHandler(conferenceStatusHandler)).Methods(http.MethodGet)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/message", wrapConferenceHandler(conferenceMessageHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(pexipNameSpace+"/rooms/{room}/participants/{part_id}/disconnect", wrapParticipantHandler(participantDisconnectHandler)).Methods(http.MethodPost)

	mgmtMux.HandleFunc(apiV1Prefix+"/rooms", wrapRequestHandler(createNewRoomHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/rooms", wrapRequestHandler(getAllRoomsHandler)).Methods(http.MethodGet)
	mgmtMux.HandleFunc(apiV1Prefix+"/rooms/{room}", wrapRequestHandler(getRoomHandler)).Methods(http.MethodGet)
	mgmtMux.HandleFunc(apiV1Prefix+"/rooms/{room}", wrapRequestHandler(deleteRoomHandler)).Methods(http.MethodDelete)

	return mgmtMux, nil
}

func wrapRequestHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(r.Method, r.Proto, r.Host+r.RequestURI)
		f(w, r)
	}
}

func getConferenceAndToken(roomName string) (*pexip.Conference, string, error) {

	var (
		room  ds.Room
		conf  *pexip.Conference
		token string
		err   error
	)

	err = dao.Rooms().GetByName(&room, roomName)

	if err != nil {
		return nil, token, err
	}

	conf = &pexip.Conference{}
	conf.Name = room.Name
	conf.Pin = room.HostPin
	conf.Alias = room.Name

	token, err = tokenStore.Get(roomName)

	if err != nil {
		switch err {
		case errors.ErrExpiredPexipToken, errors.ErrNoPexipToken:
			logger.Warning(err)
			err = tokenStore.Fetch(conf.Alias, conf.Pin)

			if err != nil {
				return nil, token, err
			}

			token, err = tokenStore.Get(roomName)

			if err != nil {
				return nil, token, err
			}

		default:
			return nil, token, err
		}
	}

	return conf, token, err
}

func wrapConferenceHandler(f confWrappedFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(r.Method, r.Proto, r.Host+r.RequestURI)
		vars := mux.Vars(r)
		confName := vars["room"]

		conf, token, err := getConferenceAndToken(confName)

		if err != nil {
			switch err {
			case errors.ErrRecordNotFound:
				logger.Error(err)
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(r.Body)

		respPayload, err := f(conf, r.Body, token)

		if err != nil {
			logger.Error(err)

			switch err {
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if len(respPayload) != 0 {
			w.Write(respPayload)
		}
	}
}

func wrapParticipantHandler(f partWrappedFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(r.Method, r.Proto, r.Host+r.RequestURI)
		vars := mux.Vars(r)
		confName := vars["room"]
		partUUID := vars["partid"]

		conf, token, err := getConferenceAndToken(confName)

		if err != nil {
			switch err {
			case errors.ErrRecordNotFound:
				logger.Error(err)
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		participant, err := participantStore.Get(partUUID)

		if err != nil {
			logger.Warning(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		out, err := f(participant, conf.Name, token)

		if err != nil {
			logger.Error(err)

			switch err {
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if len(out) > 0 {
			w.Write(out)
		}
	}
}

func pingReqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Pong", strconv.FormatInt(time.Now().Unix(), 10))
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(http.StatusOK)
}

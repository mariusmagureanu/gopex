package mux

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"bitbucket.org/kinlydev/gopex/pexip"
	logger "bitbucket.org/kinlydev/gopex/pkg/log"
)

const (
	apiV1Prefix = "/api/v1"
)

var (
	urlNameSpace = fmt.Sprintf("%s%s%s", "/pexip_monitor", apiV1Prefix, "/rooms")

	confStore        = pexip.InitConfStore()
	tokenStore       = pexip.InitTokenStore()
	participantStore = pexip.InitParticipantStore()
)

type confWrappedFunc = func(http.ResponseWriter, *http.Request, *pexip.Conference, string)
type partWrappedFunc = func(http.ResponseWriter, *http.Request, *pexip.Participant, string, string)

//TODO: to be removed
func dummyConferences() {
	mc := pexip.Conference{}
	mc.Name = "marius@test.dev.kinlycloud.net"
	mc.Pin = "6421"

	confStore.Set(&mc)
}

// InitMux initializes a new router for
// the rest api webserver.
func InitMux() (*mux.Router, error) {
	dummyConferences()

	mgmtMux := mux.NewRouter()
	mgmtMux.HandleFunc(apiV1Prefix+"/ping", wrapConferenceHandler(pingReqHandler)).Methods(http.MethodHead)

	// rest interface taken from pexws
	mgmtMux.HandleFunc(urlNameSpace+"/{room}", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants", wrapConferenceHandler(conferenceParticipantsHandler)).Methods(http.MethodGet)
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/disconnect_all", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/override_layout", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/{cmd:lock|unlock}", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/{cmd:muteguests|unmuteguests}", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/disconnect_part", wrapParticipantHandler(participantDisconnectHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/demote_host", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/promote_guest", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/lock_part", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/{cmd:unmute_part|mute_part}", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/{cmd}", wrapConferenceHandler(pingReqHandler))

	// rest interface taken from pexwebrtc
	mgmtMux.HandleFunc(apiV1Prefix+"/monitor/start/{room}", monitorStartHandler).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/monitor/stop/{room}", wrapConferenceHandler(monitorStopHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/dial", wrapConferenceHandler(conferenceDialHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/{cmd:lock|unlock|muteguests|unmuteguests}", wrapConferenceHandler(conferenceCmdHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/transform_layout", wrapConferenceHandler(transformLayoutHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/override_layout", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/disconnect", wrapConferenceHandler(conferenceDisconnectHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/transfer", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/role", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/{cmd:spotlighton|spotlightoff}", wrapConferenceHandler(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/{cmd}", wrapConferenceHandler(pingReqHandler))

	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/status", wrapConferenceHandler(conferenceStatusHandler)).Methods(http.MethodGet)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/message", wrapConferenceHandler(conferenceMessageHandler)).Methods(http.MethodPost)

	return mgmtMux, nil
}

func wrapConferenceHandler(f confWrappedFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(r.Method, r.Proto, r.Host+r.RequestURI)
		vars := mux.Vars(r)
		confName := vars["room"]

		conf, err := confStore.Get(confName)

		if err != nil {
			logger.Warning(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		token, err := tokenStore.Get(confName)

		if err != nil {
			logger.Warning(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		f(w, r, conf, token)
	}
}

func wrapParticipantHandler(f partWrappedFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(r.Method, r.Proto, r.Host+r.RequestURI)
		vars := mux.Vars(r)
		confName := vars["room"]
		partUUID := vars["partid"]

		_, err := confStore.Get(confName)

		if err != nil {
			logger.Warning(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		token, err := tokenStore.Get(confName)

		if err != nil {
			logger.Warning(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		participant, err := participantStore.Get(partUUID)

		if err != nil {
			logger.Warning(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		f(w, r, participant, confName, token)
	}
}

func pingReqHandler(w http.ResponseWriter, r *http.Request, conf *pexip.Conference, token string) {
	w.Header().Set("X-Pong", strconv.FormatInt(time.Now().Unix(), 10))
	w.Header().Set("Content-Length", "0")
}

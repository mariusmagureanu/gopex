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

	confStore  = pexip.InitConfStore()
	tokenStore = pexip.InitTokenStore()
)

type wrappedFunc = func(http.ResponseWriter, *http.Request, *pexip.Conference, string)

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
	mgmtMux.HandleFunc(apiV1Prefix+"/ping", wrapRequest(pingReqHandler)).Methods(http.MethodHead)

	// rest interface taken from pexws
	mgmtMux.HandleFunc(urlNameSpace+"/{room}", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants", wrapRequest(conferenceParticipantsHandler)).Methods(http.MethodGet)
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/disconnect_all", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/override_layout", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/{cmd:lock|unlock}", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/{cmd:muteguests|unmuteguests}", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/disconnect_part", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/demote_host", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/promote_guest", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/lock_part", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/{cmd:unmute_part|mute_part}", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants/{partid}/{cmd}", wrapRequest(pingReqHandler))

	// rest interface taken from pexwebrtc
	mgmtMux.HandleFunc(apiV1Prefix+"/monitor/start/{room}", monitorStartHandler).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/monitor/stop/{room}", wrapRequest(monitorStopHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/dial", wrapRequest(conferenceDialHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/{cmd:lock|unlock|muteguests|unmuteguests}", wrapRequest(conferenceCmdHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/transform_layout", wrapRequest(transformLayoutHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/override_layout", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/disconnect", wrapRequest(conferenceDisconnectHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/transfer", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/role", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/{cmd:spotlighton|spotlightoff}", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/{cmd}", wrapRequest(pingReqHandler))

	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/status", wrapRequest(conferenceStatusHandler)).Methods(http.MethodGet)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/message", wrapRequest(conferenceMessageHandler)).Methods(http.MethodPost)

	return mgmtMux, nil
}

func wrapRequest(f wrappedFunc) http.HandlerFunc {

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

func pingReqHandler(w http.ResponseWriter, r *http.Request, conf *pexip.Conference, token string) {
	w.Header().Set("X-Pong", strconv.FormatInt(time.Now().Unix(), 10))
	w.Header().Set("Content-Length", "0")
}

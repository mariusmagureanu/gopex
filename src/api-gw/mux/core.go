package mux

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"bitbucket.org/kinlydev/gopex/pexip"
	"bitbucket.org/kinlydev/gopex/pkg/errors"
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
	mgmtMux.HandleFunc(urlNameSpace+"/{room}/participants", wrapRequest(pingReqHandler))
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
	mgmtMux.HandleFunc(apiV1Prefix+"/monitor/start/{room}", wrapRequest(monitorStartHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/monitor/stop/{room}", wrapRequest(monitorStopHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/dial", wrapRequest(conferenceDialHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/{cmd:lock|unlock|muteguests|unmuteguests}", wrapRequest(conferenceCmdHandler)).Methods(http.MethodPost)
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/transform_layout", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/override_layout", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/disconnect", wrapRequest(conferenceDisconnectHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/transfer", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/role", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/{cmd:spotlighton|spotlightoff}", wrapRequest(pingReqHandler))
	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/participants/{part_uuid}/{cmd}", wrapRequest(pingReqHandler))

	mgmtMux.HandleFunc(apiV1Prefix+"/room/{room}/status", wrapRequest(conferenceStatusHandler)).Methods(http.MethodGet)

	return mgmtMux, nil
}

func wrapRequest(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(r.Method, r.Proto, r.Host+r.RequestURI)
		f(w, r)
	}
}

func pingReqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Pong", strconv.FormatInt(time.Now().Unix(), 10))
	w.Header().Set("Content-Length", "0")
}

func monitorStartHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["room"]

	conf, err := confStore.Get(name)

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

func monitorStopHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["room"]

	conf, err := confStore.Get(name)

	if err != nil {
		logger.Warning(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = tokenStore.Release(conf)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func conferenceCmdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	confName := vars["room"]
	cmd := vars["cmd"]

	var cmdResp []byte
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

func conferenceStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	confName := vars["room"]

	var statusResp []byte
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

	statusResp, err = conf.Status(token)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(statusResp)
}

func conferenceDisconnectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	confName := vars["room"]

	var disconnectResp []byte
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

	disconnectResp, err = conf.Disconnect(token)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(disconnectResp)

}

func conferenceDialHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	confName := vars["room"]

	var dialResp []byte
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

	dp, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	dialResp, err = conf.Dial(token, dp)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dialResp)
}

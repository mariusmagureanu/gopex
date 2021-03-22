package mux

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"bitbucket.org/kinlydev/gopex/pkg/log"
)

const (
	apiV1Prefix = "/api/v1"
)

var (
	urlNameSpace = fmt.Sprintf("%s%s%s", "/pexip_monitor", apiV1Prefix, "/rooms")
)

func InitMux() (*mux.Router, error) {
	mgmtMux := mux.NewRouter()
	mgmtMux.HandleFunc(apiV1Prefix+"/ping", wrapRequest(pingReqHandler)).Methods(http.MethodHead)

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

	return mgmtMux, nil
}

func wrapRequest(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug(r.Method, r.Proto, r.Host+r.RequestURI)
		f(w, r)
	}
}

func pingReqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Pong", strconv.FormatInt(time.Now().Unix(), 10))
	w.Header().Set("Content-Length", "0")
}

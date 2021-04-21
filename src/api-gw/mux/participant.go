package mux

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/mariusmagureanu/gopex/pexip"
	logger "github.com/mariusmagureanu/gopex/pkg/log"
)

func participantDisconnectHandler(w http.ResponseWriter, r *http.Request, p *pexip.Participant, confName, token string) {
	disconnectResp, err := p.Disconnect(confName, token)

	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(disconnectResp)
}

func participantSpotlightHandler(w http.ResponseWriter, r *http.Request, p *pexip.Participant, confName, token string) {

	var (
		err     error
		cmdResp []byte
	)

	vars := mux.Vars(r)
	cmd := vars["cmd"]

	switch cmd {
	case pexip.ParticipantSpotlightOff:
		cmdResp, err = p.SpotlightOff(confName, token)
		break
	case pexip.ParticipantSpotlightOn:
		cmdResp, err = p.SpotlightOn(confName, token)
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

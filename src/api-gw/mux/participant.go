package mux

import (
	"net/http"

	"bitbucket.org/kinlydev/gopex/pexip"
	logger "bitbucket.org/kinlydev/gopex/pkg/log"
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

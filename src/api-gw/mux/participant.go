package mux

import (
	"github.com/mariusmagureanu/gopex/pexip"
)

func participantDisconnectHandler(p *pexip.Participant, confName, token string) ([]byte, error) {
	return p.Disconnect(confName, token)
}

func participantSpotlightOnHandler(p *pexip.Participant, confName, token string) ([]byte, error) {
	return p.SpotlightOn(confName, token)
}

func participantSpotlightOffHandler(p *pexip.Participant, confName, token string) ([]byte, error) {
	return p.SpotlightOff(confName, token)
}

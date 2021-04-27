package pexip

import (
	"fmt"
	"net/http"
	"sync"

	logger "github.com/mariusmagureanu/gopex/pkg/log"
)

// ParticipantStore is a type that acts as a storage
// to keep track of currently active participants.
type ParticipantStore struct {
	store map[string]*Participant
	sync.RWMutex
}

// Set adds a new Participant to the store, using its
// uuid as key and itself as value.
func (ps *ParticipantStore) Set(uuid string) {
	ps.Lock()
	ps.store[uuid] = &Participant{UUID: uuid}
	ps.Unlock()
}

// AddMultiple creates and adds new Participants in the store
// given their assigned uuid's.
func (ps *ParticipantStore) AddMultiple(uuids []string) {
	ps.Lock()
	for _, u := range uuids {
		ps.store[u] = &Participant{UUID: u}
	}
	ps.Unlock()
}

// Get returns a Participant from the store
// given its uuid.
func (ps *ParticipantStore) Get(uuid string) (*Participant, error) {
	ps.RLock()
	defer ps.RUnlock()

	if p, found := ps.store[uuid]; found {
		return p, nil
	}

	return nil, fmt.Errorf("could not find a participant in the store, no participant found by [%s]", uuid)
}

// Remove removes a Participant from the store.
func (ps *ParticipantStore) Remove(uuid string) {
	ps.Lock()
	delete(ps.store, uuid)
	ps.Unlock()
}

// Participant is a type the represents a person
// that takes part in an active conference.
type Participant struct {
	UUID string
}

// Disconnect will disconnect the participant from the
// specified room.
func (p *Participant) Disconnect(roomName, token string) ([]byte, error) {
	logger.Debug("disconnecting participant ", p.UUID, "from room", roomName)
	disconnectURL := fmt.Sprintf("%s/%s/%s/%s/%s", urlNameSpace, roomName, "participants", p.UUID, ParticipantDisconnect)

	return doRequest(http.MethodPost, disconnectURL, token, "", []byte{})
}

// SpotlightOff will remove the spotlight for the calling participant
// in the specified room.
func (p *Participant) SpotlightOff(roomName, token string) ([]byte, error) {
	logger.Debug("setting spotlight off for participant ", p.UUID, "from room", roomName)
	spotlightOffURL := fmt.Sprintf("%s/%s/%s/%s/%s", urlNameSpace, roomName, "participants", p.UUID, ParticipantSpotlightOff)

	return doRequest(http.MethodPost, spotlightOffURL, token, "", []byte{})
}

// SpotlightOn will set the spotlight on the calling participant
// for the specified room.
func (p *Participant) SpotlightOn(roomName, token string) ([]byte, error) {
	logger.Debug("setting spotlight on for participant ", p.UUID, "from room", roomName)
	spotlightOnURL := fmt.Sprintf("%s/%s/%s/%s/%s", urlNameSpace, roomName, "participants", p.UUID, ParticipantSpotlightOn)

	return doRequest(http.MethodPost, spotlightOnURL, token, "", []byte{})
}

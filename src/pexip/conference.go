package pexip

import (
	"fmt"
	"net/http"
	"sync"

	logger "bitbucket.org/kinlydev/gopex/pkg/log"
)

// ConferenceStore is a type that acts as storage
// to keep track of running conferences.
type ConferenceStore struct {
	store map[string]*Conference
	sync.RWMutex
}

// Set adds a new conference to the store, using its
// name as key and itself as value.
func (cs *ConferenceStore) Set(conference *Conference) {
	cs.Lock()
	cs.store[conference.Name] = conference
	cs.Unlock()
}

// Get returns a Conference from the store
// given its name.
func (cs *ConferenceStore) Get(roomName string) (*Conference, error) {
	cs.RLock()
	defer cs.RUnlock()

	if conf, found := cs.store[roomName]; found {
		return conf, nil
	}

	return nil, fmt.Errorf("could not find a conference in the store, no conference found by [%s]", roomName)
}

// Remove removes a Conference from the store.
func (cs *ConferenceStore) Remove(roomName string) {
	cs.Lock()
	delete(cs.store, roomName)
	cs.Unlock()
}

type Conference struct {
	Alias string
	Name  string
	Pin   string
}

func (c *Conference) Dial(token string, dp []byte) ([]byte, error) {
	logger.Debug("dialing in into room", c.Name)
	dialUrl := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, ConferenceDial)

	return doRequest(http.MethodPost, dialUrl, token, "", dp)
}

func (c *Conference) Status(token string) ([]byte, error) {
	logger.Debug("fetching status for room", c.Name)
	statusUrl := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, ConferenceStatus)

	return doRequest(http.MethodGet, statusUrl, token, "", []byte{})
}

func (c *Conference) Lock(token string) ([]byte, error) {
	logger.Debug("locking room", c.Name)
	lockUrl := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, CommandLock)

	return doRequest(http.MethodPost, lockUrl, token, "", []byte{})
}

func (c *Conference) Unlock(token string) ([]byte, error) {
	logger.Debug("unlocking room", c.Name)
	unlockUrl := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, CommandUnlock)

	return doRequest(http.MethodPost, unlockUrl, token, "", []byte{})
}

func (c *Conference) Start() error {
	return nil
}

func (c *Conference) MuteGuests(token string) ([]byte, error) {
	logger.Debug("muting guests for room", c.Name)
	muteGuestsUrl := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, CommandMuteGuests)

	return doRequest(http.MethodPost, muteGuestsUrl, token, "", []byte{})
}

func (c *Conference) UnmuteGuests(token string) ([]byte, error) {
	logger.Debug("unmuting guests for room", c.Name)
	unmuteGuestsUrl := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, CommandUnmuteGuests)

	return doRequest(http.MethodPost, unmuteGuestsUrl, token, "", []byte{})
}

func (c *Conference) Disconnect(token string) ([]byte, error) {
	logger.Debug("disconnecting room", c.Name)
	disconnectUrl := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, ConferenceDisconnect)

	return doRequest(http.MethodPost, disconnectUrl, token, "", []byte{})
}

func (c *Conference) Message(token string, message []byte) ([]byte, error) {
	logger.Debug("sending a message to room", c.Name)
	messageUrl := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, ConferenceMessage)

	return doRequest(http.MethodPost, messageUrl, token, "", message)
}

func (c *Conference) Participants(token string) ([]byte, error) {
	logger.Debug("fetching participants from room", c.Name)
	participantsUrl := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, ConferenceParticipants)

	return doRequest(http.MethodGet, participantsUrl, token, "", []byte{})
}

func (c *Conference) TransformLayout() error {
	return nil
}

package pexip

import (
	"fmt"
	"net/http"
	"sync"

	logger "github.com/mariusmagureanu/gopex/pkg/log"
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

// Conference is a type which represents a room in Pexip.
// The "room" term will often be used as a conference as well,
// though they point to the same concept.
type Conference struct {
	Alias string
	Name  string
	Pin   string
}

// Dial will send a dial request for the specified room.
func (c *Conference) Dial(token string, dp []byte) ([]byte, error) {
	logger.Debug("dialing in into room", c.Name)
	dialURL := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, ConferenceDial)

	return doRequest(http.MethodPost, dialURL, token, "", dp)
}

// Status fetches the status for the given conference.
func (c *Conference) Status(token string) ([]byte, error) {
	logger.Debug("fetching status for room", c.Name)
	statusURL := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, ConferenceStatus)

	return doRequest(http.MethodGet, statusURL, token, "", []byte{})
}

// Lock sends a lock request for the given room.
func (c *Conference) Lock(token string) ([]byte, error) {
	logger.Debug("locking room", c.Name)
	lockURL := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, CommandLock)

	return doRequest(http.MethodPost, lockURL, token, "", []byte{})
}

// Unlock sends an unlock request for the given room.
func (c *Conference) Unlock(token string) ([]byte, error) {
	logger.Debug("unlocking room", c.Name)
	unlockURL := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, CommandUnlock)

	return doRequest(http.MethodPost, unlockURL, token, "", []byte{})
}

// Start is not yet implemented
func (c *Conference) Start() error {
	return nil
}

// MuteGuests mutes all guests in the given room.
func (c *Conference) MuteGuests(token string) ([]byte, error) {
	logger.Debug("muting guests for room", c.Name)
	muteGuestsURL := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, CommandMuteGuests)

	return doRequest(http.MethodPost, muteGuestsURL, token, "", []byte{})
}

// UnmuteGuests unmutes all guests in the given room.
func (c *Conference) UnmuteGuests(token string) ([]byte, error) {
	logger.Debug("unmuting guests for room", c.Name)
	unmuteGuestsURL := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, CommandUnmuteGuests)

	return doRequest(http.MethodPost, unmuteGuestsURL, token, "", []byte{})
}

// Disconnect disconnects the calling room.
func (c *Conference) Disconnect(token string) ([]byte, error) {
	logger.Debug("disconnecting room", c.Name)
	disconnectURL := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, ConferenceDisconnect)

	return doRequest(http.MethodPost, disconnectURL, token, "", []byte{})
}

// Message sends a text message to the given room.
func (c *Conference) Message(token string, message []byte) ([]byte, error) {
	logger.Debug("sending a message to room", c.Name)
	messageURL := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, ConferenceMessage)

	return doRequest(http.MethodPost, messageURL, token, "", message)
}

// Participants returns a slice of all participants in the calling room.
func (c *Conference) Participants(token string) ([]byte, error) {
	logger.Debug("fetching participants from room", c.Name)
	participantsURL := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, ConferenceParticipants)

	return doRequest(http.MethodGet, participantsURL, token, "", []byte{})
}

// TransformLayout sends a layout transform request for the calling room.
func (c *Conference) TransformLayout(token string, layout []byte) ([]byte, error) {
	logger.Debug("transforming the layout for room", c.Name)
	layoutURL := fmt.Sprintf("%s/%s/%s", urlNameSpace, c.Name, ConferenceTransformLayout)

	return doRequest(http.MethodPost, layoutURL, token, "", layout)
}

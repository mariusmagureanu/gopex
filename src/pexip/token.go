package pexip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"bitbucket.org/kinlydev/gopex/pkg/errors"
	logger "bitbucket.org/kinlydev/gopex/pkg/log"
)

// TokenStore is a type that handles the storage
// and lifecycle of a token.
type TokenStore struct {

	// map that stores the conference name as key
	// and token as value.
	store map[string]string

	// map that stores the conference name as key
	// and an unbuffered chanel as value, the channel
	// is used to signal a stop over refreshing tokens.
	doneListeners map[string]chan bool

	sync.RWMutex
}

// set updates the storage with a conference name
// and a new token value.
func (ts TokenStore) set(roomName, token string) {
	ts.Lock()
	ts.store[roomName] = token
	ts.Unlock()
}

// remove deletes from the storage the conference name,
// both the token and the chanel are removed.
func (ts TokenStore) remove(roomName string) {
	ts.Lock()
	delete(ts.store, roomName)
	delete(ts.doneListeners, roomName)
	ts.Unlock()
}

// refresh performs a http request against the pexip node
// and asks for a **refresh_token** for a specific conference.
func (ts TokenStore) refresh(room *Conference) error {
	urlReq := fmt.Sprintf("%s/%s/%s", urlNameSpace, room.Name, RefreshToken)
	logger.Debug("refreshing token for room", room.Name)

	currentToken, err := ts.Get(room.Name)
	if err != nil {
		return err
	}

	refreshResp, err := doRequest(http.MethodPost, urlReq, currentToken, "", []byte{})

	var refreshTokenResp tokenResponse

	err = json.Unmarshal(refreshResp, &refreshTokenResp)

	if err != nil {
		return err
	}

	ts.set(room.Name, refreshTokenResp.Result.Token)
	return nil
}

// request performs a http request against the pexip node
// and asks for a new initial **request_token** given a specific conference.
func (ts TokenStore) request(room *Conference) error {
	urlReq := fmt.Sprintf("%s/%s/%s", urlNameSpace, room.Name, RequestToken)
	logger.Debug("started watching room", room.Name)

	p := payload{DisplayName: room.Alias}
	pb, err := json.Marshal(&p)

	if err != nil {
		return err
	}

	respBody, err := doRequest(http.MethodPost, urlReq, "", room.Pin, pb)

	var tokenResp tokenResponse

	err = json.Unmarshal(respBody, &tokenResp)

	if err != nil {
		return err
	}

	ts.set(room.Name, tokenResp.Result.Token)

	return nil
}

// Get returns the current token for a given conference.
func (ts TokenStore) Get(roomName string) (string, error) {
	ts.RLock()
	defer ts.RUnlock()

	if token, found := ts.store[roomName]; found {
		return token, nil
	}

	return "", fmt.Errorf("could not find token in store, no room found by [%s]", roomName)
}

// Release performs a http request against the pexip node
// in order to revoke the current token, if succesfull the token
// storage will also clear it from its own storage.
func (ts TokenStore) Release(room *Conference) error {
	revokeUrl := fmt.Sprintf("%s/%s/%s/%s", pexipHost, urlNameSpace, room.Name, ReleaseToken)
	logger.Debug("releasing token for room", room.Name)

	currentToken, err := ts.Get(room.Name)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, revokeUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Token", currentToken)

	_, err = client.Do(req)

	if err != nil {
		return err
	}

	close(ts.doneListeners[room.Name])
	ts.remove(room.Name)

	return nil
}

// Watch retrieves an initial **request_token** for a given
// conference and starts a goroutine which will keep on
// fetching a **refresh_token** every **refreshInterval**.
func (ts TokenStore) Watch(room *Conference) error {

	currentToken, _ := ts.Get(room.Name)

	if currentToken != "" {
		return errors.ErrorRoomAlreadyStarted
	}

	err := ts.request(room)

	if err != nil {
		return err
	}

	ticker := time.NewTicker(refreshInterval)

	done := make(chan bool)
	ts.doneListeners[room.Name] = done

	go func(c *Conference) {
		for {
			select {
			case <-done:
				logger.Debug("stopped watching room", c.Name)
				return
			case <-ticker.C:
				err := ts.refresh(c)
				if err != nil {
					logger.Error(err)
					return
				}
			}
		}
	}(room)

	return nil
}

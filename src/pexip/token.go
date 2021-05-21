package pexip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mariusmagureanu/gopex/pkg/errors"
	logger "github.com/mariusmagureanu/gopex/pkg/log"
)

type token struct {
	Value     string
	Timestamp time.Time
}

// TokenStore is a type that handles the storage
// and lifecycle of a token.
type TokenStore struct {
	sync.RWMutex

	// map that stores the conference name as key
	// and token as value.
	store map[string]token
}

// set updates the storage with a conference name
// and a new token value.
func (ts *TokenStore) set(roomName, tks string) {
	ts.Lock()
	tk := token{Value: tks, Timestamp: time.Now()}
	ts.store[roomName] = tk
	ts.Unlock()
}

// remove deletes from the storage the conference name,
// both the token and the channel are removed.
func (ts *TokenStore) remove(roomName string) {
	ts.Lock()
	delete(ts.store, roomName)
	ts.Unlock()
}

// refresh performs a http request against the pexip node
// and asks for a **refresh_token** for a specific conference.
func (ts *TokenStore) refresh(room *Conference) error {
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

// Fetch performs a http request against a pexip node
// and asks for a new initial **request_token** given a specific conference.
func (ts *TokenStore) Fetch(roomAlias, pin string) error {
	urlReq := fmt.Sprintf("%s/%s/%s", urlNameSpace, roomAlias, RequestToken)
	logger.Debug("request pexip token for room", roomAlias)

	p := payload{DisplayName: roomAlias}
	pb, err := json.Marshal(&p)

	if err != nil {
		return err
	}

	respBody, err := doRequest(http.MethodPost, urlReq, "", pin, pb)

	var tokenResp tokenResponse

	err = json.Unmarshal(respBody, &tokenResp)

	if err != nil {
		return err
	}

	ts.set(roomAlias, tokenResp.Result.Token)

	return nil
}

// Get returns the current token for a given conference.
func (ts *TokenStore) Get(roomName string) (string, error) {
	ts.RLock()
	defer ts.RUnlock()

	if token, found := ts.store[roomName]; found {
		if token.Timestamp.Add(2 * time.Minute).Before(time.Now()) {
			delete(ts.store, roomName)
			return "", errors.ErrExpiredPexipToken
		}

		return token.Value, nil
	}

	return "", errors.ErrNoPexipToken
}

// Release performs a http request against the pexip node
// in order to revoke the current token, if succesfull the token
// storage will also clear it from its own storage.
func (ts *TokenStore) Release(room *Conference) error {
	revokeURL := fmt.Sprintf("%s/%s/%s/%s", pexipHost, urlNameSpace, room.Name, ReleaseToken)
	logger.Debug("releasing token for room", room.Name)

	currentToken, err := ts.Get(room.Name)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, revokeURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Token", currentToken)

	_, err = client.Do(req)

	if err != nil {
		return err
	}

	ts.remove(room.Name)

	return nil
}

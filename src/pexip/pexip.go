// Package pexip acts as a client for accesing
// Pexip's client api (https://docs.pexip.com/api_client/api_rest.html)
// The main consumer of this package is api-gw.
package pexip

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	logger "bitbucket.org/kinlydev/gopex/pkg/log"
)

const (
	CommandLock               = "lock"
	CommandUnlock             = "unlock"
	CommandMuteGuests         = "muteguests"
	CommandUnmuteGuests       = "unmuteguests"
	ConferenceStatus          = "conference_status"
	ConferenceDisconnect      = "disconnect"
	ConferenceDial            = "dial"
	ConferenceMessage         = "message"
	ConferenceParticipants    = "participants"
	ConferenceTransformLayout = "transform_layout"
	ParticipantDisconnect     = "disconnect"
	ParticipantSpotlightOff   = "spotlightoff"
	ParticipantSpotlightOn    = "spotlighton"
	RequestToken              = "request_token"
	RefreshToken              = "refresh_token"
	ReleaseToken              = "release_token"

	urlNameSpace = "api/client/v2/conferences"
)

var (
	pexipHost       string
	client          *http.Client
	refreshInterval time.Duration
)

// InitConfStore initializes a new ConferenceStore.
func InitConfStore() *ConferenceStore {
	cs := ConferenceStore{}
	cs.store = make(map[string]*Conference)

	return &cs
}

// InitTokenStore initializes a new TokenStore.
func InitTokenStore() *TokenStore {
	ts := TokenStore{}
	ts.store = make(map[string]string)
	ts.doneListeners = make(map[string](chan bool))

	return &ts
}

func InitParticipantStore() *ParticipantStore {
	ps := ParticipantStore{}
	ps.store = make(map[string]*Participant)

	return &ps
}

// InitPexipClient intializes a new http client used
// for all comunication against the Pexip api.
func InitPexipClient(host string, timeout time.Duration, maxConns, maxIdleConns int, tokenRefresh time.Duration) error {

	pexipHost = host
	refreshInterval = tokenRefresh

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConns
	t.MaxIdleConnsPerHost = maxIdleConns

	client = &http.Client{
		Timeout:   timeout,
		Transport: t,
	}

	return nil
}

// doRequest is a helper function used to perform http
// requests against the Pexip api.
// * method - http method
// * url - request path
// * token - current token, if available
// * pin - conference pin, if available
// * payload - request body, if available
func doRequest(method, url, token, pin string, payload []byte) ([]byte, error) {

	var (
		out         []byte
		payloadBuff *bytes.Buffer
		req         *http.Request
		err         error
	)

	reqURI := fmt.Sprintf("%s/%s", pexipHost, url)
	logger.Debug(method, reqURI)

	if len(payload) > 0 {
		payloadBuff = bytes.NewBuffer(payload)
		req, err = http.NewRequest(method, reqURI, payloadBuff)
	} else {
		req, err = http.NewRequest(method, reqURI, nil)
	}

	if err != nil {
		return out, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Token", token)
	}

	if pin != "" {
		req.Header.Set("Pin", pin)
	}

	resp, err := client.Do(req)

	if err != nil {
		return out, err
	}

	defer resp.Body.Close()

	out, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return out, err
	}

	return out, nil
}

// Package pexip acts as a client for accesing
// Pexip's client api (https://docs.pexip.com/api_client/api_rest.html)
// The main consumer of this package is api-gw.
package pexip

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"

	logger "github.com/mariusmagureanu/gopex/pkg/log"
)

// Const values used for requests
// against Pexip
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
	natsConn        *nats.Conn

	headerID    = []byte("id:")
	headerData  = []byte("data:")
	headerEvent = []byte("event:")
	headerRetry = []byte("retry:")
)

// InitTokenStore initializes a new TokenStore.
func InitTokenStore() *TokenStore {
	ts := TokenStore{}
	ts.store = make(map[string]token)

	return &ts
}

//InitParticipantStore initializes a new ParticipantStore.
func InitParticipantStore() *ParticipantStore {
	ps := ParticipantStore{}
	ps.store = make(map[string]*Participant)

	return &ps
}

func InitSSEManager() *SSEManager {
	sse := SSEManager{}
	sse.sseClient = &http.Client{}
	sse.cancelFuncs = make(map[string]context.CancelFunc)

	return &sse
}

// InitPexipClient initializes a new http client used
// for all communication against the Pexip api.
func InitPexipClient(host string, timeout time.Duration, maxConns, maxIdleConns int, nc *nats.Conn) error {

	natsConn = nc
	pexipHost = host

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

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logger.Error(err)
		}
	}(resp.Body)

	out, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return out, err
	}

	return out, nil
}

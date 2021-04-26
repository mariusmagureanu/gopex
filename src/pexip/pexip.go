// Package pexip acts as a client for accesing
// Pexip's client api (https://docs.pexip.com/api_client/api_rest.html)
// The main consumer of this package is api-gw.
package pexip

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	logger "github.com/mariusmagureanu/gopex/pkg/log"
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
	sseClient       *http.Client
	refreshInterval time.Duration

	headerID    = []byte("id:")
	headerData  = []byte("data:")
	headerEvent = []byte("event:")
	headerRetry = []byte("retry:")
)

// Event is a type corresponding to a server
// sent event sent by Pexip.
type Event struct {
	ID    string
	Data  string
	Event string
	Retry string
}

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
	ts.doneListeners = make(map[string]chan bool)

	return &ts
}

func InitParticipantStore() *ParticipantStore {
	ps := ParticipantStore{}
	ps.store = make(map[string]*Participant)

	return &ps
}

// InitPexipClient initializes a new http client used
// for all communication against the Pexip api.
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

	sseClient = &http.Client{
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

func sse(roomName, token string) error {
	logger.Debug("starting sse for room", roomName)
	sseURL := fmt.Sprintf("%s/%s/%s/%s", pexipHost, "api/client/v2/conferences", roomName, "events")
	req, err := http.NewRequest(http.MethodGet, sseURL, nil)

	if err != nil {
		return err
	}

	ctx, _ := context.WithCancel(context.Background())

	req = req.WithContext(ctx)

	query := req.URL.Query()
	query.Add("token", token)
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	sseResp, err := sseClient.Do(req)

	if err != nil {
		logger.Error(err)
		return err
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logger.Error(err)
		}
	}(sseResp.Body)

	bodyReader := bufio.NewReader(sseResp.Body)

	for {
		msg, err := bodyReader.ReadBytes('\n')
		if err != nil {
			logger.Error(err)
			break
		}

		_, err = processEvent(msg)

		if err != nil {
			logger.Error(err)
			break
		}
	}

	return nil
}

func processEvent(msg []byte) (Event, error) {
	var (
		e     Event
		id    []byte
		data  []byte
		event []byte
		retry []byte
	)

	if len(msg) == 0 {
		return e, fmt.Errorf("event message was empty")
	}

	bytes.Replace(msg, []byte("\n\r"), []byte("\n"), -1)
	for _, line := range bytes.FieldsFunc(msg, func(r rune) bool { return r == '\n' || r == '\r' }) {
		switch {
		case bytes.HasPrefix(line, headerID):
			id = append([]byte(nil), trimHeader(len(headerID), line)...)
		case bytes.HasPrefix(line, headerData):
			data = append(append(trimHeader(len(headerData), line), data[:]...), byte('\n'))
		case bytes.Equal(line, bytes.TrimSuffix(headerData, []byte(":"))):
			data = append(data, byte('\n'))
		case bytes.HasPrefix(line, headerEvent):
			event = append([]byte(nil), trimHeader(len(headerEvent), line)...)
		case bytes.HasPrefix(line, headerRetry):
			retry = append([]byte(nil), trimHeader(len(headerRetry), line)...)
		default:
		}
	}

	e.ID = string(id)
	e.Event = string(event)
	e.Data = string(bytes.TrimSuffix(data, []byte("\n")))
	e.Retry = string(retry)

	return e, nil
}

func trimHeader(size int, data []byte) []byte {
	data = data[size:]
	if data[0] == 32 {
		data = data[1:]
	}
	if len(data) > 0 && data[len(data)-1] == 10 {
		data = data[:len(data)-1]
	}
	return data
}

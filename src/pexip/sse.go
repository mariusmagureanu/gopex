package pexip

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/mariusmagureanu/gopex/pkg/errors"
	logger "github.com/mariusmagureanu/gopex/pkg/log"
)

// Event is a type corresponding to a server
// sent event sent by Pexip.
type Event struct {
	ID    string
	Data  string
	Event string
	Retry string
}

// SSEManager is a type that handles server sent
// events by Pexip.
type SSEManager struct {
	cancelFuncs map[string]context.CancelFunc
	sseClient   *http.Client
	sync.RWMutex
}

func (s *SSEManager) addCancelable(roomName string, cf context.CancelFunc) {
	s.Lock()
	s.cancelFuncs[roomName] = cf
	s.Unlock()
}

func (s *SSEManager) removeCancelable(roomName string) {
	s.Lock()
	delete(s.cancelFuncs, roomName)
	s.Unlock()
}

func (s *SSEManager) getCancelable(roomName string) (context.CancelFunc, error) {
	s.RLock()
	if cf, ok := s.cancelFuncs[roomName]; ok {
		return cf, nil
	}
	s.RUnlock()

	return nil, fmt.Errorf("dds")
}

// Stop cancels a sse request for the specified room.
func (s *SSEManager) Stop(roomName string) error {
	cf, err := s.getCancelable(roomName)

	if err != nil {
		return err
	}

	cf()
	s.removeCancelable(roomName)

	return nil
}

// Listen sends a request to pexip in order to listen
// for sse's for the specified room.
// This method is blocking, the caller should consider
// running it as a goroutine.
func (s *SSEManager) Listen(roomName, token string) error {
	logger.Debug("starting sse for room", roomName)

	sseURL := fmt.Sprintf("%s/%s/%s/%s", pexipHost, "api/client/v2/conferences", roomName, "events")
	req, err := http.NewRequest(http.MethodGet, sseURL, nil)

	if err != nil {
		return err
	}

	ctx, cf := context.WithCancel(context.Background())

	s.addCancelable(roomName, cf)

	req = req.WithContext(ctx)

	query := req.URL.Query()
	query.Add("token", token)
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	sseResp, err := s.sseClient.Do(req)

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
		return e, errors.ErrSSEBodyIsEmpty
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
	if data[0] == 2<<4 {
		data = data[1:]
	}
	if len(data) > 0 && data[len(data)-1] == 10 {
		data = data[:len(data)-1]
	}
	return data
}

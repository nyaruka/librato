package librato

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

// The endpoint we post librato to
var metricsEndpont = "https://metrics-api.librato.com/v1/metrics"

// basic metrics collector which can batch events
type collector struct {
	url      string
	username string
	token    string
	source   string
	timeout  time.Duration

	httpClient *http.Client
	waitGroup  *sync.WaitGroup
	stop       chan bool
	buffer     chan gauge
}

type gauge struct {
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
	MeasureTime int64   `json:"measure_time"`
}

type payload struct {
	MeasureTime int64   `json:"measure_time"`
	Source      string  `json:"source"`
	Gauges      []gauge `json:"gauges"`
}

// NewCollector creates a new librato Sender with the passed in parameters
func NewCollector(username string, token string, source string, timeout time.Duration, waitGroup *sync.WaitGroup) Collector {
	return &collector{
		url:      metricsEndpont,
		username: username,
		token:    token,
		source:   source,
		timeout:  timeout,

		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		waitGroup: waitGroup,
		stop:      make(chan bool),
		buffer:    make(chan gauge, 10000),
	}
}

// Start starts our librato sender, callers can use Stop to stop it
func (c *collector) Start() {
	c.waitGroup.Add(1)
	go func() {
		defer c.waitGroup.Done()

		slog.Info("started collector", "username", c.username, "comp", "librato")
		for {
			select {
			case <-c.stop:
				for len(c.buffer) > 0 {
					c.flush(250)
				}
				slog.Info("stopped", "comp", "librato")
				return

			case <-time.After(c.timeout):
				for i := 0; i < 4; i++ {
					c.flush(250)
				}
			}
		}
	}()
}

// Gauge can be used to add a new gauge to be sent to librato
func (c *collector) Gauge(name string, value float64) {
	// our buffer is full, log an error but continue
	if len(c.buffer) >= cap(c.buffer) {
		slog.Error("unable to add new gauges, buffer full, you may want to increase your buffer size or decrease your timeout")
		return
	}

	c.buffer <- gauge{Name: strings.ToLower(name), Value: value, MeasureTime: time.Now().Unix()}
}

// Stop stops our sender, callers can use the WaitGroup used during initialization to block for stop
func (c *collector) Stop() {
	close(c.stop)
}

func (c *collector) flush(count int) {
	if len(c.buffer) <= 0 {
		return
	}

	// build our payload
	reqPayload := &payload{
		MeasureTime: time.Now().Unix(),
		Source:      c.source,
		Gauges:      make([]gauge, 0, len(c.buffer)),
	}

	// read up to our count of gauges
readCounts:
	for i := 0; i < count; i++ {
		select {
		case g := <-c.buffer:
			reqPayload.Gauges = append(reqPayload.Gauges, g)
		default:
			break readCounts
		}
	}

	// send it off
	encoded, err := json.Marshal(reqPayload)
	if err != nil {
		slog.Error("error encoding librato metrics", "error", err, "comp", "librato")
		return
	}

	req, _ := http.NewRequest("POST", c.url, bytes.NewReader(encoded))
	req.SetBasicAuth(c.username, c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("error sending librato metrics", "error", err, "comp", "librato")
		return
	}
	// read our entire body and always close so we reuse connections
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	// non 200 or 201 are errors
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		slog.Error("non 200 returned when posting librato metrics", "comp", "librato", "status_code", resp.StatusCode)
		return
	}

	slog.Debug("flushed to librato", "comp", "librato", "count", len(reqPayload.Gauges))
}

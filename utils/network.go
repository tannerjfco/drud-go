package utils

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

// HTTPOptions defines the URL and other common HTTP options for EnsureHTTPStatus.
type HTTPOptions struct {
	URL            string
	username       string
	password       string
	timeout        time.Duration
	tickerInterval time.Duration
	expectedStatus int
	headers        map[string]string
}

// Returns a new HTTPOptions struct with some sane defaults.
func NewHTTPOptions(URL string) *HTTPOptions {
	o := HTTPOptions{
		URL:            URL,
		tickerInterval: 20,
		timeout:        60,
		expectedStatus: http.StatusOK,
	}
	return &o
}

// EnsureHTTPStatus will verify a URL responds with a given response code within the timeout period (in seconds)
func EnsureHTTPStatus(o HTTPOptions) error {
	//targetURL string, username string, password string, timeout int, expectedStatus int
	giveUp := make(chan bool)
	go func() {
		time.Sleep(time.Second * o.timeout)
		giveUp <- true
	}()

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	// Set some sane defaults if they are not passed in.
	if o.tickerInterval == 0 {
		o.tickerInterval = 20
	}
	if o.expectedStatus == 0 {
		o.expectedStatus = http.StatusOK
	}

	queryTicker := time.NewTicker(time.Second * o.tickerInterval).C
	for {
		select {
		case <-queryTicker:
			req, err := http.NewRequest("GET", o.URL, nil)
			if o.username != "" && o.password != "" {
				req.SetBasicAuth(o.username, o.password)
			}

			// Make the request
			resp, err := client.Do(req)

			if len(o.headers) > 0 {
				for header, value := range o.headers {
					req.Header.Add(header, value)
				}
			}
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == o.expectedStatus {
					// Log expected vs. actual if we do not get a match.
					log.WithFields(log.Fields{
						"URL":      o.URL,
						"expected": o.expectedStatus,
						"got":      resp.StatusCode,
					}).Info("HTTP Status code matched expectations")
					return nil
				}

				// Log expected vs. actual if we do not get a match.
				log.WithFields(log.Fields{
					"URL":      o.URL,
					"expected": o.expectedStatus,
					"got":      resp.StatusCode,
				}).Info("HTTP Status could not be matched")
			}

		case <-giveUp:
			return fmt.Errorf("No deployment found after waiting %d seconds", o.timeout)
		}
	}
}

// IsTCPPortAvailable checks a port to see if anythign is listening
// hostPort is a string like 'localhost:80'
func IsTCPPortAvailable(port int) bool {
	conn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

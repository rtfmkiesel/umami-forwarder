package umami

import (
	"net/http"

	logger "github.com/rtfmkiesel/kisslog"
)

var forwardLog = logger.New("umami/forwarder.go")

type Forwarder struct {
	client *Client
}

func (c *Client) Forward() *Forwarder {
	return &Forwarder{
		client: c,
	}
}

func (f *Forwarder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	go func() {
		if err := f.client.ImportReqWithRetries(r); err != nil {
			forwardLog.Error(err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK")) //nolint:errcheck
}

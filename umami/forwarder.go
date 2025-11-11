package umami

import (
	"net/http"

	logger "github.com/rtfmkiesel/kisslog"
)

type Forwarder struct {
	client *Client
	log    *logger.Logger
}

func (c *Client) Forward() *Forwarder {
	return &Forwarder{
		client: c,
		log:    logger.New("umami/forwarder"),
	}
}

func (f *Forwarder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	go func() {
		if err := f.client.ImportReqWithRetries(r); err != nil {
			f.log.Error(err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK")) //nolint:errcheck
}

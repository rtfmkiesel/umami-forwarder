package umami

import (
	"log"
	"net/http"
)

type Forwarder struct {
	client *Client
}

func (c *Client) Forward() *Forwarder {
	return &Forwarder{
		client: c,
	}
}

// Function for an HTTP handler, which accepts the mirrored requests
func (f *Forwarder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	go func() {
		if err := f.client.ImportReqWithRetries(r); err != nil {
			log.Printf("ERR: %s\n", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK")) //nolint:errcheck
}

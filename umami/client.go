package umami

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Client struct {
	config     *ClientConfig
	httpClient *http.Client
	sem        chan struct{} // To respect config.MaxRequests
}

// Creates a new client based on *ClientConfig
func NewClient(config *ClientConfig) (*Client, error) {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.IgnoreTLS,
		},
	}

	client := &Client{
		config: config,
		httpClient: &http.Client{
			Timeout:   time.Duration(config.Timeout) * time.Second,
			Transport: transportConfig,
		},
		sem: make(chan struct{}, config.MaxRequests),
	}

	return client, nil
}

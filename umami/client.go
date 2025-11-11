package umami

import (
	"crypto/tls"
	"net/http"
	"time"
)

type ClientConfig struct {
	WebsiteId     string // The website ID which is tagged onto the incoming requests (UUID)
	CollectionURL string // The absolute Umami collection URL, e.g. http://umami:3000/api/send
	IpHeader      string // Which header sent to the forwarder contains the real IP
	Timeout       int    // HTTP timeout in seconds when sending requests to Umami
	Retries       int    // HTTP retries when sending requests to Umami
	MaxRequests   int    // Max. concurrent HTTP requests to Umami
	IgnoreTLS     bool   // Ignore TLS errors when connecting to Umami
}

type Client struct {
	config     *ClientConfig
	httpClient *http.Client
	sem        chan struct{} // To respect MaxRequests
}

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

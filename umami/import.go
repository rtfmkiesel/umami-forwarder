package umami

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	logger "github.com/rtfmkiesel/kisslog"
)

func (c *Client) ImportReqWithRetries(r *http.Request) error {
	var log = logger.New("umami/import/ImportReqWithRetries")

	for attempt := 1; attempt <= c.config.Retries; attempt++ {
		if attempt > 1 {
			time.Sleep(time.Duration(attempt) * time.Second)
			log.Warning("Retrying import for '%s' (attempt %d/%d)", r.URL, attempt, c.config.Retries)
		}

		err := c.ImportReqOnce(r)
		if err == nil {
			return nil
		}
		log.Error(err)
	}

	return log.NewError("request for '%s' not imported: retries exhausted", r.URL)
}

func (c *Client) ImportReqOnce(r *http.Request) error {
	var log = logger.New("umami/import/ImportReqOnce")

	// Wait for a free slot
	c.sem <- struct{}{}
	defer func() {
		<-c.sem
	}()

	jsonBody, err := c.reqToUmamiRequestJsonBody(r)
	if err != nil {
		return log.NewError("failed to import event: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.config.CollectionURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return log.NewError("failed to import event: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", r.Header.Get("User-Agent")) // Set upstream user-agent as a fallback

	log.Debug("Sending request to Umami: payload='%s'", jsonBody)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return log.NewError("failed to import event: %v", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return log.NewError("failed to import event: request failed with status %d: %s", resp.StatusCode, body)
	}

	return nil
}

// Internal struct when sending Umami events
// https://umami.is/docs/api/sending-stats
type clientUmamiRequest struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

func (c *Client) reqToUmamiRequestJsonBody(r *http.Request) ([]byte, error) {
	ip := strings.TrimSpace(r.Header.Get(c.config.IpHeader))
	if ip == "" {
		return nil, fmt.Errorf("no '%s' header found, cannot create request for Umami", c.config.IpHeader)
	}

	ref := strings.TrimSpace(r.Header.Get("Referrer"))
	ua := strings.TrimSpace(r.Header.Get("User-Agent"))
	if ua == "" {
		return nil, fmt.Errorf("no 'User-Agent' header found, cannot create request for Umami")
	}

	payload := map[string]any{
		"website":    c.config.WebsiteId,
		"hostname":   r.Host,
		"url":        r.URL.String(),
		"referrer":   ref,
		"user-agent": ua,
		"ip":         ip,
	}

	req := clientUmamiRequest{
		Type:    "event",
		Payload: payload,
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	return jsonBody, nil
}

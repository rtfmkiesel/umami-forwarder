package umami

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Imports a request by forwarding it up to a maximum amount of times
// as defined by the config
func (c *Client) ImportReqWithRetries(r *http.Request) error {
	// Check if the request should be ignored
	if !c.shouldImportReq(r) {
		return nil
	}

	for attempt := 1; attempt <= c.config.Retries; attempt++ {
		if attempt > 1 {
			time.Sleep(time.Duration(attempt) * time.Second)
			log.Printf("WAR: Retrying import for '%s' (attempt %d/%d)\n", r.URL, attempt, c.config.Retries)
		}

		shouldRetry, err := c.ImportReqOnce(r)
		if err == nil {
			return nil
		}

		if !shouldRetry {
			return err
		}

		log.Printf("ERR: %s\n", err)
	}

	return fmt.Errorf("request for '%s' not imported: retries exhausted", r.URL)
}

// Imports a request by forwarding it once
func (c *Client) ImportReqOnce(r *http.Request) (bool, error) {
	// Wait for a free slot
	c.sem <- struct{}{}
	defer func() {
		<-c.sem
	}()

	jsonBody, err := c.reqToUmamiRequestJsonBody(r)
	if err != nil {
		return false, fmt.Errorf("failed to import event: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.config.CollectionURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, fmt.Errorf("failed to import event: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", r.Header.Get("User-Agent")) // Set upstream user-agent as a fallback

	log.Printf("Sending request to Umami: payload='%s'\n", jsonBody)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to import event: %v", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Errorf("failed to import event: request failed with status %d: %s", resp.StatusCode, body)

		shouldRetry := resp.StatusCode < 400 || resp.StatusCode >= 500
		return shouldRetry, errMsg
	}

	// Success
	return false, nil
}

// Internal struct when sending Umami events
// https://umami.is/docs/api/sending-stats
type clientUmamiRequest struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

// Formats a *http.Request to the JSON body needed from the Umami
// collection endpoint
func (c *Client) reqToUmamiRequestJsonBody(r *http.Request) ([]byte, error) {
	ip := strings.TrimSpace(r.Header.Get(c.config.IpHeader))
	if ip == "" {
		return nil, fmt.Errorf("no '%s' header found, cannot create request for Umami", c.config.IpHeader)
	}

	ref := strings.TrimSpace(r.Header.Get("Referer"))
	ua := strings.TrimSpace(r.Header.Get("User-Agent"))
	if ua == "" {
		return nil, fmt.Errorf("no 'User-Agent' header found, cannot create request for Umami")
	}

	payload := map[string]any{
		"website":    c.config.WebsiteId,
		"hostname":   r.Host,
		"url":        r.URL.String(),
		"referer":    ref,
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

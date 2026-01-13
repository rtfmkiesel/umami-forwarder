package umami

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type ClientConfig struct {
	WebsiteId        string          // The website ID which is tagged onto the incoming requests (UUID)
	CollectionURL    string          // The absolute Umami collection URL, e.g. http://umami:3000/api/send
	IgnoreMediaFiles bool            // Do not forward common media files
	IgnoreExtensions map[string]bool // File extension which are not going to be forwarded
	SkipFiltering    bool            // !ignoreMediaFiles && len(ignoreExt) == 0 && len(ignoreIps) == 0
	IpHeader         string          // Which header sent to the forwarder contains the real IP
	IgnoreIPs        map[string]bool // IPv4 addresses which are not going to be forwarded
	Timeout          int             // HTTP timeout in seconds when sending requests to Umami
	Retries          int             // HTTP retries when sending requests to Umami
	MaxRequests      int             // Max. concurrent HTTP requests to Umami
	IgnoreTLS        bool            // Ignore TLS errors when connecting to Umami
}

// Loads the config from ENV variables defined in README.md
func ConfigFromEnv() (*ClientConfig, error) {
	log.Printf("Loading config from environment variables\n")

	websiteID := strings.TrimSpace(os.Getenv("WEBSITE_ID"))
	if websiteID == "" {
		return nil, fmt.Errorf("environment variable WEBSITE_ID must be set")
	}

	collectionURL := strings.TrimSpace(os.Getenv("COLLECTION_URL"))
	if collectionURL == "" {
		return nil, fmt.Errorf("environment variable  COLLECTION_URL must be set")
	}

	ipHeader := strings.TrimSpace(os.Getenv("IP_HEADER"))
	if ipHeader == "" {
		return nil, fmt.Errorf("environment variable IP_HEADER must be set")
	}

	ignoreMediaFiles, err := parseBoolEnv("IGNORE_MEDIA", false)
	if err != nil {
		return nil, fmt.Errorf("invalid IGNORE_MEDIA value: %v", err)
	}

	ignoreTLS, err := parseBoolEnv("HTTP_IGNORE_TLS", false)
	if err != nil {
		return nil, fmt.Errorf("invalid IGNORE_TLS value: %v", err)
	}

	timeout, err := parseIntEnv("HTTP_TIMEOUT", 5)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_TIMEOUT value: %v", err)
	}

	retries, err := parseIntEnv("HTTP_RETRIES", 3)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_RETRIES value: %v", err)
	}

	maxRequests, err := parseIntEnv("HTTP_MAX_REQUESTS", 25)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_MAX_REQUESTS value: %v", err)
	}

	ignoreExt := parseExtensionsEnv("IGNORE_EXT")

	ignoreIps := parseIPsEnv("IGNORE_IPS")

	return &ClientConfig{
		WebsiteId:        websiteID,
		CollectionURL:    collectionURL,
		IgnoreMediaFiles: ignoreMediaFiles,
		IgnoreExtensions: ignoreExt,
		SkipFiltering:    !ignoreMediaFiles && len(ignoreExt) == 0 && len(ignoreIps) == 0,
		IpHeader:         ipHeader,
		IgnoreIPs:        ignoreIps,
		Timeout:          timeout,
		Retries:          retries,
		MaxRequests:      maxRequests,
		IgnoreTLS:        ignoreTLS,
	}, nil
}

func parseBoolEnv(key string, defaultVal bool) (bool, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}
	return strconv.ParseBool(val)
}

func parseIntEnv(key string, defaultVal int) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}

	intVal, err := strconv.Atoi(val)
	if intVal == 0 {
		return 0, fmt.Errorf("invalid %s value: cannot be 0", key)
	}

	return intVal, err
}

func parseExtensionsEnv(key string) map[string]bool {
	raw := os.Getenv(key)
	raw = strings.TrimSpace(raw)

	result := make(map[string]bool)
	if raw == "" {
		return result
	}

	for p := range strings.SplitSeq(raw, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if !strings.HasPrefix(p, ".") {
			p = "." + p
		}
		result[p] = true
	}
	return result
}

func parseIPsEnv(key string) map[string]bool {
	raw := os.Getenv(key)
	raw = strings.TrimSpace(raw)

	result := make(map[string]bool)
	if raw == "" {
		return result
	}

	for p := range strings.SplitSeq(raw, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if ip := net.ParseIP(p); ip != nil {
			result[p] = true
		}
	}
	return result
}

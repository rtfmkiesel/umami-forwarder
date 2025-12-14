package umami

import (
	"os"
	"strconv"
	"strings"

	logger "github.com/rtfmkiesel/kisslog"
)

func ConfigFromEnv() (*ClientConfig, error) {
	log := logger.New("umami/ConfigFromEnv")

	_, debug := os.LookupEnv("DEBUG")
	logger.FlagDebug = debug

	log.Debug("Loading config from environment variables")

	websiteID := os.Getenv("WEBSITE_ID")
	if websiteID == "" {
		return nil, log.NewError("WEBSITE_ID environment variable must be set")
	}

	collectionURL := os.Getenv("COLLECTION_URL")
	if collectionURL == "" {
		return nil, log.NewError("COLLECTION_URL environment variable must be set")
	}

	ignoreMediaFilesStr := os.Getenv("IGNORE_MEDIA")
	ignoreMediaFiles := false // default
	if ignoreMediaFilesStr != "" {
		b, err := strconv.ParseBool(ignoreMediaFilesStr)
		if err != nil {
			return nil, log.NewError("invalid IGNORE_MEDIA value: %v", err)
		}
		ignoreMediaFiles = b
	}

	ignoreExtRaw := os.Getenv("IGNORE_EXT")
	ignoreExt := make(map[string]bool)
	if ignoreExtRaw != "" {
		parts := strings.SplitSeq(ignoreExtRaw, ",")
		for p := range parts {
			p = strings.TrimSpace(p)
			if !strings.HasPrefix(p, ".") {
				p = "." + p // Make sure to have a dot prefix
			}
			ignoreExt[p] = true
		}
	}

	ipHeader := os.Getenv("IP_HEADER")
	if ipHeader == "" {
		return nil, log.NewError("IP_HEADER environment variable must be set")
	}

	timeoutStr := os.Getenv("HTTP_TIMEOUT")
	timeout := 5 // default
	if timeoutStr != "" {
		parsedTimeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			return nil, log.NewError("invalid TIMEOUT value: %v", err)
		}
		timeout = parsedTimeout
	}

	retriesStr := os.Getenv("HTTP_RETRIES")
	retries := 3 // default
	if retriesStr != "" {
		parsedRetries, err := strconv.Atoi(retriesStr)
		if err != nil {
			return nil, log.NewError("invalid RETRIES value: %v", err)
		}
		retries = parsedRetries
	}

	maxRequestsStr := os.Getenv("HTTP_MAX_REQUESTS")
	maxRequests := 25 // default
	if maxRequestsStr != "" {
		parsedMaxRequests, err := strconv.Atoi(maxRequestsStr)
		if err != nil {
			return nil, log.NewError("invalid MAX_REQUESTS value: %v", err)
		}
		maxRequests = parsedMaxRequests
	}

	ignoreTLSStr := os.Getenv("HTTP_IGNORE_TLS")
	ignoreTLS := false // default
	if ignoreTLSStr != "" {
		parsedIgnoreTLS, err := strconv.ParseBool(ignoreTLSStr)
		if err != nil {
			return nil, log.NewError("invalid IGNORE_TLS value: %v", err)
		}
		ignoreTLS = parsedIgnoreTLS
	}

	log.Info("Loaded config")

	return &ClientConfig{
		WebsiteId:        websiteID,
		CollectionURL:    collectionURL,
		IgnoreMediaFiles: ignoreMediaFiles,
		IgnoreExtensions: ignoreExt,
		SkipFiltering:    (!ignoreMediaFiles && len(ignoreExt) == 0),
		IpHeader:         ipHeader,
		Timeout:          timeout,
		Retries:          retries,
		MaxRequests:      maxRequests,
		IgnoreTLS:        ignoreTLS,
	}, nil
}

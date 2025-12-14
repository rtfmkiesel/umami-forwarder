package umami

import (
	"net/http"
	"path/filepath"

	logger "github.com/rtfmkiesel/kisslog"
)

var commonMediaFiles = map[string]bool{
	".jpg":   true,
	".jpeg":  true,
	".png":   true,
	".gif":   true,
	".webp":  true,
	".svg":   true,
	".bmp":   true,
	".ico":   true,
	".tiff":  true,
	".tif":   true,
	".mp4":   true,
	".webm":  true,
	".ogg":   true,
	".mov":   true,
	".avi":   true,
	".wmv":   true,
	".flv":   true,
	".mkv":   true,
	".m4v":   true,
	".mp3":   true,
	".wav":   true,
	".m4a":   true,
	".aac":   true,
	".flac":  true,
	".wma":   true,
	".opus":  true,
	".woff":  true,
	".woff2": true,
	".ttf":   true,
	".otf":   true,
	".eot":   true,
}

var filterLog = logger.New("umami/filter.go")

func (c *Client) shouldImportReq(r *http.Request) bool {
	if c.config.SkipFiltering {
		return true
	}

	ext := filepath.Ext(r.URL.Path)

	// Custom user filter
	if exists := c.config.IgnoreExtensions[ext]; exists {
		filterLog.Debug("Skipping request with extension %s (user filter)", ext)
		return false
	}

	// Common media files filter
	if exists := commonMediaFiles[ext]; exists {
		filterLog.Debug("Skipping request with extension %s (media filter)", ext)
		return false
	}

	return true
}

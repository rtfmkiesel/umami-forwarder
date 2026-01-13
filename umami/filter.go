package umami

import (
	"log"
	"net/http"
	"path/filepath"
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

// Checks based on IP filter, custom user media filter
// or common media filter if a request should be imported
func (c *Client) shouldImportReq(r *http.Request) bool {
	if c.config.SkipFiltering {
		return true
	}

	// IP filter
	ip := r.Header.Get(c.config.IpHeader)
	if exists := c.config.IgnoreIPs[ip]; exists {
		log.Printf("Skipping request from IP %s (user IP filter)\n", ip)
		return false
	}

	ext := filepath.Ext(r.URL.Path)

	// Custom user media filter
	if exists := c.config.IgnoreExtensions[ext]; exists {
		log.Printf("Skipping request with extension %s (user media filter)\n", ext)
		return false
	}

	// Common media files filter
	if exists := commonMediaFiles[ext]; exists {
		log.Printf("Skipping request with extension %s (media filter)\n", ext)
		return false
	}

	return true
}

package utils

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetHostName function returns hostname and port
func GetHostParts(uri string) (string, string) {
	tempURI := uri
	if !strings.HasPrefix(tempURI, "http://") && !strings.HasPrefix(tempURI, "https://") {
		tempURI = "https://" + tempURI
	}
	u, err := url.Parse(tempURI)
	if err != nil {
		return "localhost", "8080"
	}
	host := u.Hostname()
	port := u.Port()
	return host, port
}

// GetHost function to get host
func GetHost(c *gin.Context) string {
	currentURL := c.Request.Header.Get("x-api-url")
	if currentURL != "" {
		return strings.TrimSuffix(currentURL, "/")
	}
	scheme := c.Request.Header.Get("X-Forwarded-Proto")
	if scheme != "https" {
		scheme = "http"
	}
	host := c.Request.Host
	return strings.TrimSuffix(scheme+"://"+host, "/")
}

package models

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Stream struct {
	Name     string
	Path     string
	URL      string
	IsActive bool
}

func NewStream(basePath, baseURL string) *Stream {
	url := os.Getenv("AWS_STREAMING_URL")

	client := http.Client{Timeout: 3 * time.Second}

	isActive := false
	if resp, err := client.Get(url); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			if body, _ := io.ReadAll(resp.Body); strings.HasPrefix(string(body), "#EXTM3U") {
				isActive = true
			}
		}
	}

	return &Stream{
		Name:     "",
		Path:     "",
		URL:      url,
		IsActive: isActive,
	}
}

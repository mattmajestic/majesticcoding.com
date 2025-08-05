package models

import (
	"fmt"
	"os"
)

type Stream struct {
	Name     string
	Path     string
	URL      string
	IsActive bool
}

func NewStream(basePath, baseURL string) *Stream {
	name := os.Getenv("STREAMING_KEY")
	if name == "" {
		name = "test123"
	}
	path := basePath + name + ".m3u8"
	url := baseURL + name + ".m3u8"
	isActive := true
	if _, err := os.Stat(path); err == nil {
		isActive = true
	} else {
		fmt.Println("Stream inactive. Reason:", err)
	}
	return &Stream{
		Name:     name,
		Path:     path,
		URL:      url,
		IsActive: isActive,
	}
}

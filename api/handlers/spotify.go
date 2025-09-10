package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
)

// GET /api/spotify/current
func SpotifyCurrent(c *gin.Context) {
	if spClient == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not connected to Spotify; visit /api/spotify/login"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	log.Println("Fetching currently playing track from Spotify...")
	cp, err := spClient.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		log.Printf("Spotify API error: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	
	// Nothing playing or no device
	if cp == nil || cp.Item == nil {
		log.Println("No track currently playing or no active device")
		c.JSON(http.StatusOK, gin.H{
			"is_playing": false,
			"message":    "No track currently playing or no active device",
		})
		return
	}

	// replace your album handling + numeric fields with this:
	item := cp.Item

	artists := make([]string, 0, len(item.Artists))
	for _, a := range item.Artists {
		artists = append(artists, a.Name)
	}

	// Album is a value type in v2, so no nil check
	albumName := item.Album.Name
	albumImg := ""
	if len(item.Album.Images) > 0 {
		albumImg = item.Album.Images[0].URL
	}

	url := ""
	if item.ExternalURLs != nil {
		if v, ok := item.ExternalURLs["spotify"]; ok {
			url = v
		}
	}

	resp := models.CurrentTrack{
		Title:      item.Name,
		Artists:    artists,
		Album:      albumName,
		AlbumImage: albumImg,
		URL:        url,
		IsPlaying:  cp.Playing,
		ProgressMS: int(cp.Progress),   // cast from spotify.Numeric
		DurationMS: int(item.Duration), // cast from spotify.Numeric
	}
	c.JSON(http.StatusOK, resp)
}

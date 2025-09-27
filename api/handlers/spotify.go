package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
	"majesticcoding.com/api/services"
)

// GET /api/spotify/current
func SpotifyCurrent(c *gin.Context) {
	if spClient == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not connected to Spotify; visit /api/spotify/login"})
		return
	}

	cacheKey := "spotify:current:30s"

	// Try to get from Redis cache first (30 seconds TTL)
	var cachedTrack models.CurrentTrack
	err := services.RedisGetJSON(cacheKey, &cachedTrack)
	if err == nil {
		log.Printf("✅ Spotify current track cache HIT")
		c.JSON(http.StatusOK, cachedTrack)
		return
	}
	log.Printf("🔍 Spotify current track cache MISS, fetching from API")

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
		notPlayingResp := gin.H{
			"is_playing": false,
			"message":    "No track currently playing or no active device",
		}

		// Cache "not playing" status for shorter time (10 seconds)
		if err := services.RedisSetJSON("spotify:notplaying:10s", notPlayingResp, 10); err != nil {
			log.Printf("⚠️ Failed to cache Spotify not-playing status: %v", err)
		} else {
			log.Printf("💾 Cached Spotify not-playing status for 10 seconds")
		}

		c.JSON(http.StatusOK, notPlayingResp)
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

	// Cache the current track for 30 seconds
	if err := services.RedisSetJSON(cacheKey, resp, 30); err != nil {
		log.Printf("⚠️ Failed to cache Spotify current track: %v", err)
	} else {
		log.Printf("💾 Cached Spotify track '%s' for 30 seconds", resp.Title)
	}

	c.JSON(http.StatusOK, resp)
}

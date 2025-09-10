package models

type CurrentTrack struct {
	Title      string   `json:"title"`
	Artists    []string `json:"artists"`
	Album      string   `json:"album"`
	AlbumImage string   `json:"album_image"`
	URL        string   `json:"url"`
	IsPlaying  bool     `json:"is_playing"`
	ProgressMS int      `json:"progress_ms"`
	DurationMS int      `json:"duration_ms"`
}

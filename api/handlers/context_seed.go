package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
)

// SeedContext seeds bronze.website_context with site content, GitHub repos, and social stats.
// Protected — requires a valid Clerk token.
func SeedContext(c *gin.Context) {
	results := map[string]string{}

	// 1. Hardcoded site pages
	if err := seedSitePages(); err != nil {
		results["site_pages"] = err.Error()
	} else {
		results["site_pages"] = "ok"
	}

	// 2. GitHub repos
	repoCount, err := seedGitHubRepos("mattmajestic")
	if err != nil {
		results["github_repos"] = err.Error()
	} else {
		results["github_repos"] = fmt.Sprintf("seeded %d repos", repoCount)
	}

	// 3. Social stats from DB
	if err := services.StoreLatestSocialStatsContextFromDB(); err != nil {
		results["social_stats"] = err.Error()
	} else {
		results["social_stats"] = "ok"
	}

	c.JSON(http.StatusOK, gin.H{"seeded": results})
}

func seedSitePages() error {
	pages := []struct {
		contentType string
		title       string
		content     string
		url         string
		priority    int
	}{
		{
			"site_info", "About Majestic Coding",
			`Majestic Coding is a personal developer site built and maintained by Matt (mattmajestic).
The site showcases live coding streams, open source projects, and serves as a hub for the Majestic Coding community.
It features an AI chat assistant, real-time WebSocket chat during streams, GitHub/YouTube/Twitch/LeetCode stats,
a Spotify now-playing widget, geographic check-ins, and a live stream viewer powered by AWS IVS.`,
			"https://majesticcoding.com/about", 5,
		},
		{
			"tech_stack", "Majestic Coding Tech Stack",
			`Majestic Coding is built with: Go (Gin framework) for the backend, Tailwind CSS and HTMX for the frontend,
PostgreSQL on Neon for the database, Redis for caching, Clerk for authentication,
AWS IVS (Interactive Video Service) for RTMP ingest and HLS live stream playback,
Google Cloud Run for containerized deployment, Docker for containerization,
pgvector for vector embeddings and RAG, and Gemini AI for embeddings and chat responses.
The site also integrates OpenAI, Anthropic Claude, and Groq as alternative AI providers.`,
			"https://majesticcoding.com", 4,
		},
		{
			"social_presence", "Majestic Coding Socials",
			`Majestic Coding social channels:
GitHub at github.com/mattmajestic — open source projects, tools, and tutorials primarily in Go, TypeScript, Docker, and Kubernetes.
YouTube channel "Majestic Coding" — coding tutorials, live stream recordings, and tech walkthroughs.
Twitch channel MajesticCodingTwitch — live coding sessions streamed in real time, integrated into this site's /live page.
The Twitch chat and site WebSocket chat are shown side by side during streams.`,
			"https://majesticcoding.com", 5,
		},
		{
			"live_streaming", "Majestic Coding Live Streams",
			`Majestic Coding streams live coding sessions on Twitch (MajesticCodingTwitch) and YouTube.
Streams cover building this site live, Go backend development, DevOps, AI feature implementation, and viewer Q&A.
The /live page on majesticcoding.com embeds the stream via AWS IVS HLS playback and shows the Twitch chat alongside a site-native WebSocket chat.
Viewers can participate in chat from the site without needing a Twitch account.`,
			"https://majesticcoding.com/live", 4,
		},
		{
			"site_features", "Majestic Coding Site Features",
			`Key features of majesticcoding.com:
AI Chat (/ai) — multi-model AI assistant backed by pgvector RAG using site content as context.
Live Streaming (/live) — AWS IVS HLS player with real-time WebSocket chat and Twitch integration.
Stats dashboard — live GitHub, YouTube, Twitch, and LeetCode stats with Redis caching.
Check-ins map — geographic location check-ins displayed on an interactive globe.
Spotify widget — shows currently playing track via Spotify OAuth.
Infrastructure page — Kubernetes deployment info and cloud architecture overview.`,
			"https://majesticcoding.com", 3,
		},
	}

	for _, p := range pages {
		if err := services.StoreWebsiteContext(p.contentType, p.title, p.content, p.url, nil, p.priority); err != nil {
			return fmt.Errorf("failed to seed %q: %w", p.title, err)
		}
	}
	return nil
}

func seedGitHubRepos(username string) (int, error) {
	token := os.Getenv("GITHUB_TOKEN")
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&sort=updated", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var repos []struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Language    string   `json:"language"`
		Stars       int      `json:"stargazers_count"`
		Topics      []string `json:"topics"`
		HTMLURL     string   `json:"html_url"`
		Fork        bool     `json:"fork"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return 0, err
	}

	count := 0
	for _, r := range repos {
		if r.Fork {
			continue // skip forks, only original work
		}
		desc := r.Description
		if desc == "" {
			desc = "No description provided"
		}
		content := fmt.Sprintf(
			"GitHub repository: %s by %s. Description: %s. Primary language: %s. Stars: %d. Topics: %v. URL: %s",
			r.Name, username, desc, r.Language, r.Stars, r.Topics, r.HTMLURL,
		)
		priority := 2
		if r.Stars >= 5 {
			priority = 3
		}
		if err := services.StoreWebsiteContext("github_repo", r.Name, content, r.HTMLURL, nil, priority); err != nil {
			fmt.Printf("⚠️ Failed to seed repo %s: %v\n", r.Name, err)
			continue
		}
		count++
	}
	return count, nil
}

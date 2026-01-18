package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"majesticcoding.com/api/models"
	"majesticcoding.com/db"
)

type GeminiEmbeddingRequest struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}

type GeminiEmbeddingResponse struct {
	Embedding struct {
		Values []float64 `json:"values"`
	} `json:"embedding"`
}

// GenerateEmbedding creates an embedding for the given text using Gemini
func GenerateEmbedding(text string) ([]float64, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not found")
	}

	reqBody := GeminiEmbeddingRequest{
		Content: struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			Parts: []struct {
				Text string `json:"text"`
			}{
				{Text: text},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/embedding-001:embedContent?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API error: %d - %s", resp.StatusCode, string(body))
	}

	var embeddingResp GeminiEmbeddingResponse
	if err := json.Unmarshal(body, &embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return embeddingResp.Embedding.Values, nil
}

// ConvertStatsToText converts unified stats to readable text for embeddings
func ConvertStatsToText(stats *models.UnifiedStats) []string {
	var contexts []string

	// YouTube context
	if stats.YouTube != nil && stats.YouTube.Error == "" {
		text := fmt.Sprintf("YouTube Channel: %s has %d subscribers, %d total views across %d videos. This represents Majestic Coding content creation on YouTube platform.",
			stats.YouTube.ChannelName, stats.YouTube.Subscribers, stats.YouTube.Views, stats.YouTube.Videos)
		contexts = append(contexts, text)
	}

	// GitHub context
	if stats.GitHub != nil && stats.GitHub.Error == "" {
		text := fmt.Sprintf("GitHub Profile: %s has %d public repositories, %d followers, and has received %d total stars. This shows Majestic Coding open source development activity.",
			stats.GitHub.Username, stats.GitHub.PublicRepos, stats.GitHub.Followers, stats.GitHub.StarsReceived)
		contexts = append(contexts, text)
	}

	// Twitch context
	if stats.Twitch != nil && stats.Twitch.Error == "" {
		text := fmt.Sprintf("Twitch Channel: %s (%s) has %d followers. Channel type: %s. This represents Majestic Coding live streaming presence.",
			stats.Twitch.DisplayName, stats.Twitch.Description, stats.Twitch.Followers, stats.Twitch.BroadcasterType)
		contexts = append(contexts, text)
	}

	// LeetCode context
	if stats.LeetCode != nil && stats.LeetCode.Error == "" {
		text := fmt.Sprintf("LeetCode Profile: %s has solved %d problems, ranked #%d globally. Primary languages: %s. This shows Majestic Coding competitive programming skills.",
			stats.LeetCode.Username, stats.LeetCode.SolvedCount, stats.LeetCode.Ranking, stats.LeetCode.Languages)
		contexts = append(contexts, text)
	}

	return contexts
}

// StoreWebsiteContext stores general website context with embeddings
func StoreWebsiteContext(contentType, title, content, sourceURL string, metadata interface{}, priority int) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("database not available")
	}

	// Generate embedding
	embedding, err := GenerateEmbedding(content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Convert embedding to PostgreSQL array format
	embeddingJSON, err := json.Marshal(embedding)
	if err != nil {
		return fmt.Errorf("failed to marshal embedding: %w", err)
	}

	var metadataJSON []byte
	if metadata != nil {
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	// Insert or update context (upsert based on content_type and title)
	_, err = database.Exec(`
		INSERT INTO bronze.website_context (content_type, title, content_text, source_url, metadata, embedding, priority)
		VALUES ($1, $2, $3, $4, $5, $6::vector, $7)
		ON CONFLICT (content_type, title)
		DO UPDATE SET
			content_text = EXCLUDED.content_text,
			source_url = EXCLUDED.source_url,
			metadata = EXCLUDED.metadata,
			embedding = EXCLUDED.embedding,
			priority = EXCLUDED.priority,
			updated_at = CURRENT_TIMESTAMP
	`, contentType, title, content, sourceURL, string(metadataJSON), string(embeddingJSON), priority)

	if err != nil {
		return fmt.Errorf("failed to store website context: %w", err)
	}

	fmt.Printf("✅ Stored %s context: %s\n", contentType, title)
	return nil
}

// StoreSocialStatsContext stores social media stats as website context
func StoreSocialStatsContext(stats *models.UnifiedStats) error {
	contexts := ConvertStatsToText(stats)

	for i, contextText := range contexts {
		var contentType, title string
		var metadata interface{}

		switch i {
		case 0:
			if stats.YouTube != nil {
				contentType = "social_stats"
				title = "YouTube Channel Stats"
				metadata = stats.YouTube
			}
		case 1:
			if stats.GitHub != nil {
				contentType = "social_stats"
				title = "GitHub Profile Stats"
				metadata = stats.GitHub
			}
		case 2:
			if stats.Twitch != nil {
				contentType = "social_stats"
				title = "Twitch Channel Stats"
				metadata = stats.Twitch
			}
		case 3:
			if stats.LeetCode != nil {
				contentType = "social_stats"
				title = "LeetCode Profile Stats"
				metadata = stats.LeetCode
			}
		}

		if contentType == "" {
			continue
		}

		// Store with high priority (3) for current stats
		if err := StoreWebsiteContext(contentType, title, contextText, "", metadata, 3); err != nil {
			fmt.Printf("Failed to store %s: %v\n", title, err)
		}
	}

	return nil
}

// StoreLatestSocialStatsContextFromDB fetches latest stats rows and stores them as context.
func StoreLatestSocialStatsContextFromDB() error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("database not available")
	}

	var failures []string

	if stats, err := db.GetLatestYouTubeStats(database); err == nil {
		content := fmt.Sprintf(
			"YouTube Channel: %s has %d subscribers, %d total views across %d videos.",
			stats.ChannelName, stats.Subscribers, stats.Views, stats.Videos,
		)
		if err := StoreWebsiteContext("social_stats", "YouTube Channel Stats", content, "", stats, 3); err != nil {
			failures = append(failures, fmt.Sprintf("youtube: %v", err))
		}
	} else {
		failures = append(failures, fmt.Sprintf("youtube: %v", err))
	}

	if stats, err := db.GetLatestGitHubStats(database); err == nil {
		content := fmt.Sprintf(
			"GitHub Profile: %s has %d public repositories, %d followers, and %d total stars.",
			stats.Username, stats.PublicRepos, stats.Followers, stats.StarsReceived,
		)
		if err := StoreWebsiteContext("social_stats", "GitHub Profile Stats", content, "", stats, 3); err != nil {
			failures = append(failures, fmt.Sprintf("github: %v", err))
		}
	} else {
		failures = append(failures, fmt.Sprintf("github: %v", err))
	}

	if stats, err := db.GetLatestTwitchStats(database); err == nil {
		content := fmt.Sprintf(
			"Twitch Channel: %s has %d followers.",
			stats.DisplayName, stats.Followers,
		)
		if err := StoreWebsiteContext("social_stats", "Twitch Channel Stats", content, "", stats, 3); err != nil {
			failures = append(failures, fmt.Sprintf("twitch: %v", err))
		}
	} else {
		failures = append(failures, fmt.Sprintf("twitch: %v", err))
	}

	if stats, err := db.GetLatestLeetCodeStats(database, "mattmajestic"); err == nil {
		content := fmt.Sprintf(
			"LeetCode Profile: %s has solved %d problems and is ranked #%d. Primary languages: %s.",
			stats.Username, stats.SolvedCount, stats.Ranking, stats.Languages,
		)
		if err := StoreWebsiteContext("social_stats", "LeetCode Profile Stats", content, "", stats, 3); err != nil {
			failures = append(failures, fmt.Sprintf("leetcode: %v", err))
		}
	} else {
		failures = append(failures, fmt.Sprintf("leetcode: %v", err))
	}

	if len(failures) > 0 {
		return fmt.Errorf("failed to refresh contexts: %s", strings.Join(failures, ", "))
	}

	return nil
}

// RetrieveRelevantContext finds relevant context based on user query
func RetrieveRelevantContext(query string, limit int) ([]string, error) {
	database := db.GetDB()
	if database == nil {
		return nil, fmt.Errorf("database not available")
	}

	// Generate embedding for the query
	queryEmbedding, err := GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	queryEmbeddingJSON, err := json.Marshal(queryEmbedding)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query embedding: %w", err)
	}

	// Search for similar context using cosine similarity, prioritizing active content
	rows, err := database.Query(`
		SELECT title, content_text, content_type, priority,
		       (embedding <=> $1::vector) as distance
		FROM bronze.website_context
		WHERE embedding IS NOT NULL AND is_active = true
		ORDER BY priority DESC, distance ASC
		LIMIT $2
	`, string(queryEmbeddingJSON), limit)

	if err != nil {
		return nil, fmt.Errorf("failed to query context: %w", err)
	}
	defer rows.Close()

	var contexts []string
	for rows.Next() {
		var title, contentText, contentType string
		var priority int
		var distance float64

		if err := rows.Scan(&title, &contentText, &contentType, &priority, &distance); err != nil {
			continue
		}

		// Include relevant context (distance < 0.4 for Gemini embeddings)
		if distance < 0.4 {
			contexts = append(contexts, contentText)
		}
	}

	return contexts, nil
}

// CreatePersonalityContext creates a summary of the user's online presence
func CreatePersonalityContext() string {
	return `I am Matt, a software engineer and content creator. Here's my online presence:

	- I create coding content on YouTube
	- I maintain open source projects on GitHub
	- I live stream programming and tech content on Twitch
	- I solve competitive programming problems on LeetCode
	My content focuses on practical software development, DevOps, cloud technologies, and programming best practices.`
}

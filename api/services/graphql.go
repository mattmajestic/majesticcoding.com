package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"majesticcoding.com/api/models"
)

// GetUnifiedStats fetches all stats concurrently and returns unified response
func GetUnifiedStats(ctx context.Context) (*models.UnifiedStats, error) {
	stats := &models.UnifiedStats{}

	// Use WaitGroup for concurrent API calls
	var wg sync.WaitGroup

	// YouTube stats
	wg.Add(1)
	go func() {
		defer wg.Done()
		if youtubeStats, err := fetchYouTubeStatsGQL(ctx); err != nil {
			log.Printf("YouTube stats error: %v", err)
			stats.YouTube = &models.YouTubeStatsGQL{Error: err.Error()}
		} else {
			stats.YouTube = youtubeStats
		}
	}()

	// GitHub stats
	wg.Add(1)
	go func() {
		defer wg.Done()
		if githubStats, err := fetchGitHubStatsGQL(ctx); err != nil {
			log.Printf("GitHub stats error: %v", err)
			stats.GitHub = &models.GitHubStatsGQL{Error: err.Error()}
		} else {
			stats.GitHub = githubStats
		}
	}()

	// Twitch stats
	wg.Add(1)
	go func() {
		defer wg.Done()
		if twitchStats, err := fetchTwitchStatsGQL(ctx); err != nil {
			log.Printf("Twitch stats error: %v", err)
			stats.Twitch = &models.TwitchStatsGQL{Error: err.Error()}
		} else {
			stats.Twitch = twitchStats
		}
	}()

	// LeetCode stats
	wg.Add(1)
	go func() {
		defer wg.Done()
		if leetcodeStats, err := fetchLeetCodeStatsGQL(ctx); err != nil {
			log.Printf("LeetCode stats error: %v", err)
			stats.LeetCode = &models.LeetCodeStatsGQL{Error: err.Error()}
		} else {
			stats.LeetCode = leetcodeStats
		}
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	// Store the stats with embeddings for RAG context (async)
	go func() {
		if err := StoreSocialStatsContext(stats); err != nil {
			log.Printf("Failed to store social stats context: %v", err)
		}
	}()

	return stats, nil
}

// fetchYouTubeStatsGQL converts YouTube API response to GraphQL format
func fetchYouTubeStatsGQL(ctx context.Context) (*models.YouTubeStatsGQL, error) {
	// Use existing YouTube service logic
	stats, err := FetchYouTubeStats()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch YouTube stats: %w", err)
	}

	// Convert string values to integers
	subscribers := 0
	views := 0
	videos := 0
	channelName := ""

	if val, ok := stats["channelTitle"].(string); ok {
		channelName = val
	}
	if val, ok := stats["subscribers"].(string); ok {
		if converted, err := strconv.Atoi(val); err == nil {
			subscribers = converted
		}
	}
	if val, ok := stats["views"].(string); ok {
		if converted, err := strconv.Atoi(val); err == nil {
			views = converted
		}
	}
	if val, ok := stats["videos"].(string); ok {
		if converted, err := strconv.Atoi(val); err == nil {
			videos = converted
		}
	}

	return &models.YouTubeStatsGQL{
		ChannelName: channelName,
		Subscribers: subscribers,
		Views:       views,
		Videos:      videos,
	}, nil
}

// fetchGitHubStatsGQL converts GitHub API response to GraphQL format
func fetchGitHubStatsGQL(ctx context.Context) (*models.GitHubStatsGQL, error) {
	// Use existing GitHub service logic
	stats, err := FetchGitHubStats("mattmajestic")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub stats: %w", err)
	}

	return &models.GitHubStatsGQL{
		Username:      stats.Username,
		PublicRepos:   stats.PublicRepos,
		Followers:     stats.Followers,
		StarsReceived: stats.StarsReceived,
	}, nil
}

// fetchTwitchStatsGQL converts Twitch API response to GraphQL format
func fetchTwitchStatsGQL(ctx context.Context) (*models.TwitchStatsGQL, error) {
	// Use existing Twitch service logic
	stats, err := FetchTwitchStats("MajesticCodingTwitch")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Twitch stats: %w", err)
	}

	return &models.TwitchStatsGQL{
		DisplayName:     stats.DisplayName,
		Description:     stats.Description,
		BroadcasterType: stats.BroadcasterType,
		Followers:       stats.Followers,
	}, nil
}

// fetchLeetCodeStatsGQL converts LeetCode API response to GraphQL format
func fetchLeetCodeStatsGQL(ctx context.Context) (*models.LeetCodeStatsGQL, error) {
	// Use existing LeetCode service logic
	stats, err := FetchLeetCodeStats("mattmajestic")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch LeetCode stats: %w", err)
	}

	return &models.LeetCodeStatsGQL{
		Username:    stats.Username,
		Languages:   stats.Languages,
		SolvedCount: stats.SolvedCount,
		Ranking:     stats.Ranking,
	}, nil
}

// ExecuteGraphQLQuery handles GraphQL query execution
func ExecuteGraphQLQuery(ctx context.Context, query string, variables map[string]interface{}) (*models.GraphQLResponse, error) {
	// Simple query parser - check if query contains "unifiedStats"
	if strings.Contains(query, "unifiedStats") {
		// Parse which providers are requested
		includeYouTube := strings.Contains(query, "youtube")
		includeGitHub := strings.Contains(query, "github")
		includeTwitch := strings.Contains(query, "twitch")
		includeLeetCode := strings.Contains(query, "leetcode")

		stats, err := GetSelectiveStats(ctx, includeYouTube, includeGitHub, includeTwitch, includeLeetCode)
		if err != nil {
			return &models.GraphQLResponse{
				Errors: []models.GraphQLError{{Message: err.Error()}},
			}, nil
		}

		return &models.GraphQLResponse{
			Data: map[string]interface{}{
				"unifiedStats": stats,
			},
		}, nil
	}

	return &models.GraphQLResponse{
		Errors: []models.GraphQLError{{Message: "Unsupported query. Try: query { unifiedStats { youtube { channelName subscribers views videos } github { username publicRepos followers starsReceived } twitch { displayName description broadcasterType followers } leetcode { username languages solvedCount ranking } } }"}},
	}, nil
}

// GetSelectiveStats fetches only the requested stats
func GetSelectiveStats(ctx context.Context, includeYouTube, includeGitHub, includeTwitch, includeLeetCode bool) (*models.UnifiedStats, error) {
	stats := &models.UnifiedStats{}

	// Use WaitGroup for concurrent API calls
	var wg sync.WaitGroup

	// YouTube stats
	if includeYouTube {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if youtubeStats, err := fetchYouTubeStatsGQL(ctx); err != nil {
				log.Printf("YouTube stats error: %v", err)
				stats.YouTube = &models.YouTubeStatsGQL{Error: err.Error()}
			} else {
				stats.YouTube = youtubeStats
			}
		}()
	}

	// GitHub stats
	if includeGitHub {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if githubStats, err := fetchGitHubStatsGQL(ctx); err != nil {
				log.Printf("GitHub stats error: %v", err)
				stats.GitHub = &models.GitHubStatsGQL{Error: err.Error()}
			} else {
				stats.GitHub = githubStats
			}
		}()
	}

	// Twitch stats
	if includeTwitch {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if twitchStats, err := fetchTwitchStatsGQL(ctx); err != nil {
				log.Printf("Twitch stats error: %v", err)
				stats.Twitch = &models.TwitchStatsGQL{Error: err.Error()}
			} else {
				stats.Twitch = twitchStats
			}
		}()
	}

	// LeetCode stats
	if includeLeetCode {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if leetcodeStats, err := fetchLeetCodeStatsGQL(ctx); err != nil {
				log.Printf("LeetCode stats error: %v", err)
				stats.LeetCode = &models.LeetCodeStatsGQL{Error: err.Error()}
			} else {
				stats.LeetCode = leetcodeStats
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return stats, nil
}

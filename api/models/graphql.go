package models

// GraphQL unified stats response
type UnifiedStats struct {
	YouTube  *YouTubeStatsGQL  `json:"youtube,omitempty"`
	GitHub   *GitHubStatsGQL   `json:"github,omitempty"`
	Twitch   *TwitchStatsGQL   `json:"twitch,omitempty"`
	LeetCode *LeetCodeStatsGQL `json:"leetcode,omitempty"`
	Error    *string           `json:"error,omitempty"`
}

// GraphQL-specific stats models (cleaner field names for GraphQL)
type YouTubeStatsGQL struct {
	ChannelName string `json:"channelName"`
	Subscribers int    `json:"subscribers"`
	Views       int    `json:"views"`
	Videos      int    `json:"videos"`
	Error       string `json:"error,omitempty"`
}

type GitHubStatsGQL struct {
	Username      string `json:"username"`
	PublicRepos   int    `json:"publicRepos"`
	Followers     int    `json:"followers"`
	StarsReceived int    `json:"starsReceived"`
	Error         string `json:"error,omitempty"`
}

type TwitchStatsGQL struct {
	DisplayName     string `json:"displayName"`
	Description     string `json:"description"`
	BroadcasterType string `json:"broadcasterType"`
	Followers       int    `json:"followers"`
	Error           string `json:"error,omitempty"`
}

type LeetCodeStatsGQL struct {
	Username    string `json:"username"`
	Languages   string `json:"languages"`
	SolvedCount int    `json:"solvedCount"`
	Ranking     int    `json:"ranking"`
	Error       string `json:"error,omitempty"`
}

// GraphQL request/response models
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   interface{}    `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

type GraphQLError struct {
	Message string        `json:"message"`
	Path    []interface{} `json:"path,omitempty"`
}

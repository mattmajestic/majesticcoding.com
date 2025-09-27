package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type UpstashRedisClient struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

type RedisResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
}

var upstashClient *UpstashRedisClient

// InitRedis initializes the Redis client
func InitRedis() error {
	redisURL := os.Getenv("UPSTASH_REDIS_REST_URL")
	redisToken := os.Getenv("UPSTASH_REDIS_REST_TOKEN")

	if redisURL == "" || redisToken == "" {
		return fmt.Errorf("Redis credentials not found in environment")
	}

	upstashClient = &UpstashRedisClient{
		BaseURL: redisURL,
		Token:   redisToken,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// Test the connection with a simple ping
	_, err := upstashClient.executeCommand([]interface{}{"PING"})
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}

// executeCommand executes a Redis command via REST API
func (c *UpstashRedisClient) executeCommand(command []interface{}) (interface{}, error) {
	if c == nil {
		return nil, fmt.Errorf("Upstash Redis client not initialized")
	}

	// Convert command to JSON
	commandJSON, err := json.Marshal(command)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal command: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(commandJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response
	var redisResp RedisResponse
	if err := json.Unmarshal(body, &redisResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if redisResp.Error != "" {
		return nil, fmt.Errorf("Redis error: %s", redisResp.Error)
	}

	return redisResp.Result, nil
}

// RedisSet stores a key-value pair with optional TTL (in seconds)
func RedisSet(key string, value interface{}, ttlSeconds int) error {
	if upstashClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	var command []interface{}
	if ttlSeconds > 0 {
		command = []interface{}{"SETEX", key, ttlSeconds, value}
	} else {
		command = []interface{}{"SET", key, value}
	}

	_, err := upstashClient.executeCommand(command)
	if err != nil {
		log.Printf("Redis SET error for key %s: %v", key, err)
		return err
	}

	return nil
}

// RedisGet retrieves a value by key
func RedisGet(key string) (string, error) {
	if upstashClient == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}

	result, err := upstashClient.executeCommand([]interface{}{"GET", key})
	if err != nil {
		log.Printf("Redis GET error for key %s: %v", key, err)
		return "", err
	}

	if result == nil {
		return "", nil // Key doesn't exist
	}

	return fmt.Sprintf("%v", result), nil
}

// RedisDelete removes a key
func RedisDelete(key string) error {
	if upstashClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	_, err := upstashClient.executeCommand([]interface{}{"DEL", key})
	if err != nil {
		log.Printf("Redis DELETE error for key %s: %v", key, err)
		return err
	}

	return nil
}

// RedisSetJSON stores a JSON object with optional TTL
func RedisSetJSON(key string, value interface{}, ttlSeconds int) error {
	if upstashClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	// Marshal value to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value to JSON: %v", err)
	}

	return RedisSet(key, string(jsonValue), ttlSeconds)
}

// RedisGetJSON retrieves and unmarshals a JSON object
func RedisGetJSON(key string, dest interface{}) error {
	if upstashClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	jsonStr, err := RedisGet(key)
	if err != nil {
		return err
	}

	if jsonStr == "" {
		return fmt.Errorf("key not found")
	}

	return json.Unmarshal([]byte(jsonStr), dest)
}

// RedisGetRawJSON retrieves raw JSON string without unmarshaling
func RedisGetRawJSON(key string) (string, error) {
	if upstashClient == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}

	return RedisGet(key)
}

// RedisClearPattern deletes keys matching a pattern (for cache clearing)
func RedisClearPattern(pattern string) error {
	if upstashClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	// Get all keys matching pattern
	result, err := upstashClient.executeCommand([]interface{}{"KEYS", pattern})
	if err != nil {
		log.Printf("Redis KEYS error for pattern %s: %v", pattern, err)
		return err
	}

	// Redis KEYS returns an array of strings
	if keys, ok := result.([]interface{}); ok && len(keys) > 0 {
		// Delete each key
		for _, key := range keys {
			if keyStr, ok := key.(string); ok {
				err := RedisDelete(keyStr)
				if err != nil {
					log.Printf("Failed to delete key %s: %v", keyStr, err)
				}
			}
		}
		log.Printf("Cleared %d cache keys matching pattern: %s", len(keys), pattern)
	}

	return nil
}

// RedisGetKeys returns all keys matching a pattern
func RedisGetKeys(pattern string) ([]string, error) {
	if upstashClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}

	// Get all keys matching pattern
	result, err := upstashClient.executeCommand([]interface{}{"KEYS", pattern})
	if err != nil {
		log.Printf("Redis KEYS error for pattern %s: %v", pattern, err)
		return nil, err
	}

	var keys []string
	// Redis KEYS returns an array of strings
	if keyList, ok := result.([]interface{}); ok {
		for _, key := range keyList {
			if keyStr, ok := key.(string); ok {
				keys = append(keys, keyStr)
			}
		}
	}

	return keys, nil
}

// RedisSetAdd adds a member to a Redis set with optional TTL
func RedisSetAdd(setKey, member string, ttlSeconds int) error {
	if upstashClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	// Add member to set
	_, err := upstashClient.executeCommand([]interface{}{"SADD", setKey, member})
	if err != nil {
		log.Printf("Redis SADD error for set %s: %v", setKey, err)
		return err
	}

	// Set TTL if provided
	if ttlSeconds > 0 {
		_, err = upstashClient.executeCommand([]interface{}{"EXPIRE", setKey, ttlSeconds})
		if err != nil {
			log.Printf("Redis EXPIRE error for set %s: %v", setKey, err)
			return err
		}
	}

	return nil
}

// RedisSetCount returns the number of members in a Redis set
func RedisSetCount(setKey string) (int, error) {
	if upstashClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}

	result, err := upstashClient.executeCommand([]interface{}{"SCARD", setKey})
	if err != nil {
		log.Printf("Redis SCARD error for set %s: %v", setKey, err)
		return 0, err
	}

	// Convert result to int
	if count, ok := result.(float64); ok {
		return int(count), nil
	}

	return 0, fmt.Errorf("unexpected result type for SCARD")
}

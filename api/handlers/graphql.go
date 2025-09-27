package handlers

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
	"majesticcoding.com/api/services"
)

// GraphQLHandler handles GraphQL queries
// @Summary Execute GraphQL query
// @Description Execute a GraphQL query for unified stats
// @Tags GraphQL
// @Accept json
// @Produce json
// @Param request body models.GraphQLRequest true "GraphQL Query"
// @Success 200 {object} models.GraphQLResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/graphql [post]
func GraphQLHandler(c *gin.Context) {
	var request models.GraphQLRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Create cache key based on query hash
	queryHash := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%v", request.Query, request.Variables))))
	cacheKey := fmt.Sprintf("graphql:query:%s", queryHash)

	// Try to get from Redis cache first (10 minutes TTL = 600 seconds)
	var cachedResponse models.GraphQLResponse
	err := services.RedisGetJSON(cacheKey, &cachedResponse)
	if err == nil {
		log.Printf("✅ GraphQL query cache HIT for hash: %s", queryHash[:8])
		c.JSON(http.StatusOK, cachedResponse)
		return
	}
	log.Printf("🔍 GraphQL query cache MISS for hash: %s", queryHash[:8])

	// Execute the GraphQL query
	response, err := services.ExecuteGraphQLQuery(c.Request.Context(), request.Query, request.Variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cache successful responses for 10 minutes (600 seconds)
	if response != nil && len(response.Errors) == 0 {
		if err := services.RedisSetJSON(cacheKey, response, 600); err != nil {
			log.Printf("⚠️ Failed to cache GraphQL response: %v", err)
		} else {
			log.Printf("💾 Cached GraphQL response for 10 minutes (hash: %s)", queryHash[:8])
		}
	}

	c.JSON(http.StatusOK, response)
}

// GraphQLPlaygroundHandler serves the GraphQL playground template
// @Summary GraphQL Playground
// @Description Interactive GraphQL query interface
// @Tags GraphQL
// @Produce html
// @Success 200 {string} string "HTML page"
// @Router /api/graphql/playground [get]
func GraphQLPlaygroundHandler(c *gin.Context) {
	c.HTML(200, "graphql.tmpl", gin.H{
		"Title": "GraphQL Playground",
	})
}

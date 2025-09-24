package handlers

import (
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

	// Execute the GraphQL query
	response, err := services.ExecuteGraphQLQuery(c.Request.Context(), request.Query, request.Variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

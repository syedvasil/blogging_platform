package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go.mongodb.org/mongo-driver/mongo"
	"strconv"

	"blog-platform/config"
	dbmongo "blog-platform/database/mongo"
	"blog-platform/internal/middleware"
)

// @title Blog Platform API
// @version 1.0
// @description API documentation for the Blog Platform.

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg := config.Constcfg

	// Initialize MongoDB connection
	dbConn := dbmongo.InitDB(cfg.DB.URI)

	// Create a new Gin router
	server := gin.Default()

	// Apply authentication middleware
	server.Use(middleware.AuthMiddleware(dbConn))

	// Define API routes
	router := server.RouterGroup
	v1 := router.Group("/api/v1")

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup API routes for users and posts
	setupV1UserRoutes(dbConn, v1)
	setupV1PostRoutes(dbConn, v1)

	// Start the server
	err := server.Run(":" + strconv.Itoa(int(cfg.App.Port)))
	if err != nil {
		return
	}
}

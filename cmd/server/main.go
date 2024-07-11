package main

import (
	dbmongo "blog-platform/database/mongo"
	"blog-platform/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"blog-platform/config"
	_ "blog-platform/docs"
	"strconv"
)

// @title Blog Platform API
// @version 1.0
// @description API documentation for the Blog Platform.

// @contact.name   Syed
// @contact.url    https://www.linkedin.com/in/syed-vasil/
// @contact.email  syedvasil@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	// Load configuration
	cfg := config.Constcfg

	// Initialize MongoDB connection
	dbConn := dbmongo.InitDB(cfg.DB.URI)

	// Create a new Gin router
	server := gin.Default()

	router := server.RouterGroup
	// Swagger documentation endpoint
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Apply authentication middleware
	server.Use(middleware.AuthMiddleware(dbConn))

	// Define API routes
	v1 := router.Group("/api/v1")
	// Setup API routes for users and posts
	setupV1UserRoutes(dbConn, v1)
	setupV1PostRoutes(dbConn, v1)

	// Start the server
	err := server.Run(":" + strconv.Itoa(int(cfg.App.Port)))
	if err != nil {
		return
	}
}

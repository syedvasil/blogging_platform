package main

import (
	"blog-platform/config"
	dbmongo "blog-platform/database/mongo"
	"blog-platform/internal/middleware"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {

	// commenting this out as use const would be simpler in current situation
	cfg := config.Constcfg

	dbConn := dbmongo.InitDB(cfg.DB.URI)

	server := gin.Default()
	server.Use(middleware.AuthMiddleware(dbConn))

	router := server.RouterGroup
	v1 := router.Group("/api/v1")
	setupV1UserRoutes(dbConn, v1)
	setupV1PostRoutes(dbConn, v1)

	err := server.Run(":" + strconv.Itoa(int(cfg.App.Port)))
	if err != nil {
		return
	}
}

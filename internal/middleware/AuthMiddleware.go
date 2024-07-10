package middleware

import (
	repoModels "blog-platform/internal/app/repositories/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
)

func AuthMiddleware(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check username and password in MongoDB
		var user repoModels.User
		filter := bson.M{"username": username}
		err := db.Collection("users").FindOne(context.Background(), filter).Decode(&user)
		if err != nil || user.Password != password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// User authenticated
		c.Set("Username", user.Username)
		c.Set("ID", user.ID.Hex())
		c.Set("Role", user.Role)
		c.Next()
	}
}

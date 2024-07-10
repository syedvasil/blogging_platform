package main

import (
	ctrlPost "blog-platform/internal/app/controller/post"
	ctrlUser "blog-platform/internal/app/controller/user"
	"blog-platform/internal/app/repositories/post"
	"blog-platform/internal/app/repositories/user"
	srvPost "blog-platform/internal/app/service/post"
	srvUser "blog-platform/internal/app/service/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupV1UserRoutes(db *mongo.Database, routerGroup *gin.RouterGroup) {
	userCtrl := ctrlUser.New(srvUser.New(user.New(db)))
	userGroup := routerGroup.Group("/user")
	{
		userGroup.POST("", userCtrl.CreateUser)
		userGroup.GET("", userCtrl.GetUsers)
		userGroup.GET("/:id", userCtrl.GetUser)
		userGroup.PUT("/:id", userCtrl.UpdateUser)
		userGroup.DELETE("/:id", userCtrl.DeleteUser)
	}
}
func setupV1PostRoutes(db *mongo.Database, routerGroup *gin.RouterGroup) {
	postController := ctrlPost.New(srvPost.New(post.New(db)))
	postGroup := routerGroup.Group("/posts")
	{
		postGroup.POST("", postController.CreatePost)
		postGroup.GET("", postController.GetPosts)
		postGroup.GET("/:id", postController.GetPost)
		postGroup.PUT("/:id", postController.UpdatePost)
		postGroup.DELETE("/:id", postController.DeletePost)
	}
}

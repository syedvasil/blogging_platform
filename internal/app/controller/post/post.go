package post

import (
	"blog-platform/internal/app/controller/models"
	repoModels "blog-platform/internal/app/repositories/models"
	"blog-platform/internal/utils"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

type Service interface {
	CreatePost(post repoModels.Post) error
	GetPosts(ctx context.Context, author, date string, page, limit int) ([]repoModels.Post, *repoModels.ListMetaData, error)
	GetPostByID(id primitive.ObjectID) (repoModels.Post, error)
	UpdatePost(post repoModels.Post, access models.UserAccess) error
	DeletePost(id primitive.ObjectID, access models.UserAccess) error
}

type Controller struct {
	service Service
}

func New(service Service) *Controller {
	return &Controller{service}
}

func (c *Controller) CreatePost(ctx *gin.Context) {
	var req models.PostReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAccess := models.UserAccess{}
	err := userAccess.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"err": "resource cannot be accessed reason:" + err.Error()})
		return
	}

	pModel := models.CreatePostFromReq(req, userAccess)
	if err = c.service.CreatePost(pModel); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, pModel)
}

func (c *Controller) GetPosts(ctx *gin.Context) {
	username := ctx.Query("username")
	date := ctx.Query("date")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	posts, pagi, err := c.service.GetPosts(ctx, username, date, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.ListPostReq{Data: posts, Metadata: *pagi})
}

func (c *Controller) GetPost(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	resPost, err := c.service.GetPostByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resPost)
}

func (c *Controller) UpdatePost(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req models.PostReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAccess := models.UserAccess{}
	err = userAccess.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"err": "resource cannot be accessed reason:" + err.Error()})
		return
	}

	var updatePost repoModels.Post
	updatePost.ID = id
	if err = c.service.UpdatePost(updatePost, userAccess); err != nil {
		utils.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, updatePost)
}

func (c *Controller) DeletePost(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	userAccess := models.UserAccess{}
	err = userAccess.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"err": "resource cannot be accessed reason:" + err.Error()})
		return
	}

	if err = c.service.DeletePost(id, userAccess); err != nil {
		utils.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

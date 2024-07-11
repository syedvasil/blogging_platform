package post

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"

	"blog-platform/internal/app/controller/models"
	repoModels "blog-platform/internal/app/repositories/models"
	"blog-platform/internal/utils"
)

//go:generate mockery --name=Service --case underscore
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

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post with the input payload
// @Tags posts
// @Accept json
// @Produce json
// @Param post body models.PostReq true "Post"
// @Success 201 {object} repoModels.Post
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /posts [post]
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

// GetPosts godoc
// @Summary Get all posts
// @Description Get a list of all posts with optional filters
// @Tags posts
// @Produce json
// @Param username query string false "Author username"
// @Param date query string false "Creation date"
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} models.ListPostReq
// @Failure 500 {object} gin.H
// @Router /posts [get]
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

// GetPost godoc
// @Summary Get a post by ID
// @Description Get details of a post by ID
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} repoModels.Post
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /posts/{id} [get]
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

// UpdatePost godoc
// @Summary Update a post
// @Description Update post details by ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body models.PostReq true "Post"
// @Success 200 {object} repoModels.Post
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /posts/{id} [put]
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

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a post by ID
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /posts/{id} [delete]
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

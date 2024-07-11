package user

import (
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
	CreateUser(user repoModels.User) error
	GetUsers(page, limit int) ([]repoModels.User, error)
	GetUserByID(id primitive.ObjectID) (repoModels.User, error)
	UpdateUser(user repoModels.User, access models.UserAccess) error
	DeleteUser(id primitive.ObjectID, access models.UserAccess) error
}

type Controller struct {
	service Service
}

func New(service Service) *Controller {
	return &Controller{service}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserReq true "User"
// @Success 201 {object} repoModels.User
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	var userReq models.UserReq
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := models.CreateUserFromReq(userReq)
	if err := c.service.CreateUser(newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, newUser)
}

// GetUsers godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {array} repoModels.User
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users [get]
func (c *Controller) GetUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	access := models.UserAccess{}
	err := access.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"err": "resource cannot be accessed reason:" + err.Error()})
		return
	}

	users, err := c.service.GetUsers(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get details of a user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} repoModels.User
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/{id} [get]
func (c *Controller) GetUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	resUser, err := c.service.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resUser)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body repoModels.User true "User"
// @Success 200 {object} repoModels.User
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/{id} [put]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var userUpdate repoModels.User
	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access := models.UserAccess{}
	err = access.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"err": "resource cannot be accessed reason:" + err.Error()})
		return
	}

	userUpdate.ID = id
	if err = c.service.UpdateUser(userUpdate, access); err != nil {
		utils.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, userUpdate)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Soft delete a user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/{id} [delete]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	access := models.UserAccess{}
	err = access.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"err": "resource cannot be accessed reason:" + err.Error()})
		return
	}

	if err = c.service.DeleteUser(id, access); err != nil {
		utils.HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

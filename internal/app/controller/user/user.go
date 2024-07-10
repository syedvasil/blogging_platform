package user

import (
	"blog-platform/internal/app/controller/models"
	repoModels "blog-platform/internal/app/repositories/models"
	"blog-platform/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

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

func (c *Controller) GetUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	access := models.UserAccess{}
	err := access.GetUserFromCtx(ctx)
	if err != nil {
		//|| *access.Role != "admin"
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

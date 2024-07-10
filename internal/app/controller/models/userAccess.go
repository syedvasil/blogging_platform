package models

import (
	"errors"
	"time"

	repoModels "blog-platform/internal/app/repositories/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserAccess struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
	Role *string            `json:"role"`
}

type PostReq struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title   string             `bson:"title" json:"title"`
	Content string             `bson:"content" json:"content"`
}

type ListPostReq struct {
	Data     []repoModels.Post
	Metadata repoModels.ListMetaData
}

type UserReq struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"-"`
}

func (userA *UserAccess) GetUserFromCtx(ctx *gin.Context) error {
	username, ok := ctx.Get("Username")
	if !ok {
		return errors.New("user not found")
	}
	if nw, err := username.(string); !err {
		return errors.New("user not found")
	} else {
		userA.Name = nw
	}

	var idString string
	id, ok := ctx.Get("ID")
	if !ok {
		return errors.New("id not found")
	}

	if idS, err := id.(string); !err {
		return errors.New("id not found")
	} else {
		idString = idS
	}

	if nID, err := primitive.ObjectIDFromHex(idString); err != nil {
		return errors.New("id not found")
	} else {
		userA.ID = nID
	}

	role, ok := ctx.Get("Role")
	if !ok {
		return errors.New("role not found")
	}
	if nRole, err := role.(string); err {
		userA.Role = nil
	} else {
		userA.Role = &nRole
	}

	return nil
}

func CreatePostFromReq(req PostReq, userAccess UserAccess) repoModels.Post {
	now := time.Now()
	return repoModels.Post{
		ID:        primitive.NewObjectID(),
		Title:     req.Title,
		Content:   req.Content,
		UpdatedAt: now,
		CreatedAt: now,
		Author: repoModels.BasicUser{
			ID:       userAccess.ID,
			Username: userAccess.Name,
		},
	}
}

func CreateUserFromReq(req UserReq) repoModels.User {
	now := time.Now()
	return repoModels.User{
		ID:        primitive.NewObjectID(),
		Username:  req.Username,
		Password:  req.Password,
		Role:      "user", //admin user has to be created directly on DB or should have separate API non-public facing
		CreatedAt: now,
	}
}

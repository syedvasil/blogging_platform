package post

import (
	"blog-platform/internal/app/controller/models"
	repoModels "blog-platform/internal/app/repositories/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//go:generate mockery --name=Repository --case underscore
type Repository interface {
	CreatePost(post repoModels.Post) error
	GetPosts(ctx context.Context, filter interface{}, offset, limit int) ([]repoModels.Post, *repoModels.ListMetaData, error)
	GetPostByID(id primitive.ObjectID) (repoModels.Post, error)
	UpdatePost(post repoModels.Post) error
	DeletePost(id primitive.ObjectID) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) CreatePost(post repoModels.Post) error {
	return s.repo.CreatePost(post)
}

func (s *Service) GetPosts(ctx context.Context, username, date string, page, limit int) ([]repoModels.Post, *repoModels.ListMetaData, error) {
	offset := (page - 1) * limit
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}

	if username != "" {
		filter["author.username"] = username
	}
	if date != "" {
		startDate, _ := time.Parse("2006-01-02", date)
		endDate := startDate.AddDate(0, 0, 1)
		filter["created_at"] = bson.M{"$gte": startDate, "$lt": endDate}
	}

	return s.repo.GetPosts(ctx, filter, offset, limit)
}

func (s *Service) GetPostByID(id primitive.ObjectID) (repoModels.Post, error) {
	return s.repo.GetPostByID(id)
}

func (s *Service) UpdatePost(post repoModels.Post, access models.UserAccess) error {
	err := s.GetPostAndAuthorise(post.ID, access)
	if err != nil {
		return err
	}

	return s.repo.UpdatePost(post)
}

func (s *Service) DeletePost(id primitive.ObjectID, access models.UserAccess) error {
	err := s.GetPostAndAuthorise(id, access)
	if err != nil {
		return err
	}

	return s.repo.DeletePost(id)
}

func (s *Service) GetPostAndAuthorise(id primitive.ObjectID, access models.UserAccess) error {
	post, err := s.repo.GetPostByID(id)
	if err != nil {
		return err
	}

	if post.Author.ID == access.ID || *access.Role == "admin" {
		return nil
	}

	return errors.New("not allowed")
}

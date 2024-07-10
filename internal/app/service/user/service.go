package user

import (
	"blog-platform/internal/app/controller/models"
	repoModels "blog-platform/internal/app/repositories/models"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	CreateUser(user repoModels.User) error
	GetUsers(filter interface{}, offset, limit int) ([]repoModels.User, error)
	GetUserByID(id primitive.ObjectID) (repoModels.User, error)
	UpdateUser(user repoModels.User) error
	DeleteUser(id primitive.ObjectID) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) CreateUser(user repoModels.User) error {
	return s.repo.CreateUser(user)
}

func (s *Service) GetUsers(page, limit int) ([]repoModels.User, error) {
	offset := (page - 1) * limit
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	return s.repo.GetUsers(filter, offset, limit)
}

func (s *Service) GetUserByID(id primitive.ObjectID) (repoModels.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *Service) UpdateUser(user repoModels.User, access models.UserAccess) error {
	err := s.GetUserAndAuthorise(user.ID, access)
	if err != nil {
		return err
	}

	return s.repo.UpdateUser(user)
}

func (s *Service) DeleteUser(id primitive.ObjectID, access models.UserAccess) error {
	err := s.GetUserAndAuthorise(id, access)
	if err != nil {
		return err
	}

	return s.repo.DeleteUser(id)
}

func (s *Service) GetUserAndAuthorise(id primitive.ObjectID, access models.UserAccess) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	if user.ID == access.ID || *access.Role == "admin" {
		return nil
	}

	return errors.New("not allowed")
}

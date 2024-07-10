package user

import (
	repoModels "blog-platform/internal/app/repositories/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "users"

type Repository struct {
	db *mongo.Collection
}

func New(db *mongo.Database) *Repository {
	return &Repository{db: db.Collection(collectionName)}
}

func (r *Repository) CreateUser(user repoModels.User) error {
	_, err := r.db.InsertOne(context.Background(), user)
	return err
}

func (r *Repository) GetUsers(filter interface{}, offset, limit int) ([]repoModels.User, error) {
	var users []repoModels.User
	ctx := context.Background()

	cursor, err := r.db.Find(ctx, filter, options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var user repoModels.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, cursor.Err()
}

func (r *Repository) GetUserByID(id primitive.ObjectID) (repoModels.User, error) {
	var user repoModels.User
	err := r.db.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	return user, err
}

func (r *Repository) UpdateUser(user repoModels.User) error {
	_, err := r.db.ReplaceOne(context.Background(), bson.M{"_id": user.ID}, user)
	return err
}

func (r *Repository) DeleteUser(id primitive.ObjectID) error {
	now := time.Now()
	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": &now}})
	return err
}

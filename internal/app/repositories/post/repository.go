package post

import (
	repoModels "blog-platform/internal/app/repositories/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const CollName = "posts"

type Repository struct {
	db *mongo.Collection
}

func New(db *mongo.Database) *Repository {
	return &Repository{db: db.Collection(CollName)}
}

func (r *Repository) CreatePost(post repoModels.Post) error {
	_, err := r.db.InsertOne(context.Background(), post)
	return err
}

func (r *Repository) GetPosts(ctx context.Context, filter interface{}, offset, limit int) ([]repoModels.Post, *repoModels.ListMetaData, error) {
	var posts []repoModels.Post

	cursor, err := r.db.Find(ctx, filter, options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)))
	if err != nil {
		return nil, nil, err
	}

	for cursor.Next(ctx) {
		var post repoModels.Post
		err := cursor.Decode(&post)
		if err != nil {
			return nil, nil, err
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, nil, err
	}

	total, err := r.db.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	pagination := repoModels.ListMetaData{
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}

	return posts, &pagination, nil

}

func (r *Repository) GetPostByID(id primitive.ObjectID) (repoModels.Post, error) {
	var post repoModels.Post
	err := r.db.FindOne(context.Background(), bson.M{"_id": id}).Decode(&post)
	return post, err
}

func (r *Repository) UpdatePost(post repoModels.Post) error {
	_, err := r.db.ReplaceOne(context.Background(), bson.M{"_id": post.ID}, post)
	return err
}

func (r *Repository) DeletePost(id primitive.ObjectID) error {
	now := time.Now()
	_, err := r.db.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": &now}})
	return err
}

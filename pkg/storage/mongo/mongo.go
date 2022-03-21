package mongo

import (
	"context"
	"errors"

	"github.com/MarySmirnova/go_news/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ storage.Interface = &Store{}
var ctx context.Context = context.Background()

const (
	NameDB         = "news"
	NameCollection = "tasks"
)

type Store struct {
	db *mongo.Client
}

func New(connstr string) (*Store, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connstr))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Store{
		db: client,
	}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.db.Database(NameDB).Collection(NameCollection)
	filter := bson.D{{}}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	allPosts := []storage.Post{}
	for cur.Next(ctx) {
		p := storage.Post{}
		err = cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		allPosts = append(allPosts, p)
	}

	return allPosts, nil
}

func (s *Store) AddPost(p storage.Post) error {
	collection := s.db.Database(NameDB).Collection(NameCollection)
	_, err := collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdatePost(p storage.Post) error {
	collection := s.db.Database(NameDB).Collection(NameCollection)
	filter := bson.D{primitive.E{Key: "id", Value: p.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "title", Value: p.Title},
		primitive.E{Key: "content", Value: p.Content},
		primitive.E{Key: "author_id", Value: p.AuthorID},
		primitive.E{Key: "author_name", Value: p.AuthorName},
	}}}

	res := collection.FindOneAndUpdate(ctx, filter, update)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (s *Store) DeletePost(p storage.Post) error {
	collection := s.db.Database(NameDB).Collection(NameCollection)
	filter := bson.D{primitive.E{Key: "id", Value: p.ID}}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no posts were deleted")
	}

	return nil
}

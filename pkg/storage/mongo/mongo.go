package mongo

import (
	"context"

	"github.com/MarySmirnova/go_news/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
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
	filter := bson.D{}

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

func (s *Store) UpdatePost(storage.Post) error {

	return nil
}

func (s *Store) DeletePost(storage.Post) error {

	return nil
}

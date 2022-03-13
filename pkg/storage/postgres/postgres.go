package postgres

import (
	"context"
	"time"

	"github.com/MarySmirnova/go_news/pkg/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ storage.Interface = &Store{}
var ctx context.Context = context.Background()

type Store struct {
	db *pgxpool.Pool
}

func New(connstr string) (*Store, error) {
	db, err := pgxpool.Connect(ctx, connstr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(ctx); err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	qu := `
	SELECT 
		posts.id,
    	authors.id,
    	authors.name,
    	posts.title,
    	posts.content,
    	posts.created_at
	FROM
		posts,
		authors
	WHERE posts.author_id = authors.id;`

	rows, err := s.db.Query(ctx, qu)
	if err != nil {
		return nil, err
	}

	allPosts := []storage.Post{}
	for rows.Next() {
		p := storage.Post{}
		err = rows.Scan(&p.ID, &p.AuthorID, &p.AuthorName, &p.Title, &p.Content, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		allPosts = append(allPosts, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allPosts, nil
}

func (s *Store) AddPost(p storage.Post) error {
	qu := `
	INSERT INTO posts (
		author_id,
		title,
		content,
		created_at)
	VALUES ($1, $2, $3, $4);`

	_, err := s.db.Exec(ctx, qu, p.AuthorID, p.Title, p.Content, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdatePost(p storage.Post) error {
	qu := `
	UPDATE posts
	SET 
		author_id = $1,
		title = $2,
		content = $3
	WHERE 
		id = $4;`

	_, err := s.db.Exec(ctx, qu, p.AuthorID, p.Title, p.Content, p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeletePost(p storage.Post) error {
	qu := `
	DELETE FROM posts
	WHERE id = $1`

	_, err := s.db.Exec(ctx, qu, p.ID)
	if err != nil {
		return err
	}

	return nil
}

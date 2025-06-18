package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserId    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	Version   int       `json:"version"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserId, pq.Array(post.Tags)).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostsStore) GetByID(ctx context.Context, postID int64) (*Post, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `
		SELECT id, content, title, user_id, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`

	post := Post{}

	err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserId,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostsStore) UpdateByID(ctx context.Context, post *Post) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		UPDATE posts
		SET title = $1, content = $2, version = version + 1
		WHERE id = $3  AND version = $4 
		RETURNING version
	`

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostsStore) DeleteByID(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	sqlResult, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	affected, err := sqlResult.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostsStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `
		SELECT  
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		INNER JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		LEFT JOIN comments c ON c.post_id  = p.id
		LEFT JOIN users u ON u.id = p.user_id
		WHERE 
			(f.user_id = $1 OR p.user_id = $1)
			AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
			AND (p.tags @> $5 OR $5 = '{}')
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2
		OFFSET $3;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Term, pq.Array(fq.Tags))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata = []PostWithMetadata{}

	for rows.Next() {
		var post PostWithMetadata
		err := rows.Scan(
			&post.ID,
			&post.UserId,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentsCount,
		)

		if err != nil {
			return nil, err
		}

		feed = append(feed, post)
	}

	return feed, nil
}

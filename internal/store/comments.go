package store

import (
	"context"
	"database/sql"
)

type CommentsStore struct {
	db *sql.DB
}

type Comment struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	PostID    int64  `json:"post_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `
		INSERT INTO comments (user_id, post_id, content)
		VALUES ($1, $2, $3) RETURNING id, user_id, post_id, content, created_at
	`
	err := s.db.QueryRowContext(ctx, query, comment.UserID, comment.PostID, comment.Content).Scan(
		&comment.ID,
		&comment.UserID,
		&comment.PostID,
		&comment.Content,
		&comment.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *CommentsStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `
		SELECT 
			c.id, c.post_id, c.user_id, c.content, c.created_at, 
			u.username, u.id
		FROM comments c
		INNER JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC;
	`

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		c := Comment{}
		c.User = User{}
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.ID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

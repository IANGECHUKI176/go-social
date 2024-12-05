package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}
type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID int64, userID int64) error {
	query := `
	INSERT INTO followers (follower_id,user_id) VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followerID, userID)
	if err != nil {

		if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
			return ErrConflict
		}
	}
	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID int64, userID int64) error {
	query := `
	DELETE FROM followers
	WHERE follower_id = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followerID, userID)

	return err
}

func (s *FollowerStore) ExistsFollow(ctx context.Context, followerID int64, userID int64) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1 FROM followers WHERE follower_id = $1 AND user_id = $2)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var exists bool
	err := s.db.QueryRowContext(ctx, query, followerID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil

}

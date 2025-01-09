package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gopher_social/internal/store"
	"time"

	"github.com/go-redis/redis/v8"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	keyCache := fmt.Sprintf("user-%v", userID)
	data, err := s.rdb.Get(ctx, keyCache).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	keyCache := fmt.Sprintf("user-%v", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, keyCache, json, UserExpTime).Err()
}

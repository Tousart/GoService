package redis

import (
	"context"
	"encoding/json"
	"errors"
	"httpServer/server/config"
	"httpServer/server/domain"
	"httpServer/server/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionsRepository struct {
	client      *redis.Client
	timeSession time.Duration
}

func NewSessionsRepository(cfg config.Redis) (*SessionsRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.DefaultPassword,
		DB:       cfg.DBNumber,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &SessionsRepository{
		client:      client,
		timeSession: cfg.Duration,
	}, nil
}

func (rs *SessionsRepository) GetSessionId(sessionId string) (string, error) {
	value, err := rs.client.Get(context.Background(), sessionId).Bytes()
	if errors.Is(err, redis.Nil) {
		return "", repository.ErrNotFoundSessionId
	} else if err != nil {
		return "", err
	}

	var session domain.Session
	if err := json.Unmarshal(value, &session); err != nil {
		return "", err
	}

	return session.SessionId, nil
}

func (rs *SessionsRepository) PostSessionId(session *domain.Session) (string, error) {
	sessionId := session.SessionId
	sessionJson, err := json.Marshal(session)
	if err != nil {
		return "", nil
	}

	if err := rs.client.Set(context.Background(), sessionId, sessionJson, rs.timeSession).Err(); err != nil {
		return "", err
	}

	return sessionId, nil
}

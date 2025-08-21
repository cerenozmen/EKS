package service

import (
	"context"
	"encoding/json"
	"event-service/model"
	"event-service/repository"
	"fmt"
	
	"time"

	"github.com/redis/go-redis/v9"
)

type EventService struct {
	repo        *repository.EventRepository
	redisClient *redis.Client
}

func NewEventService(r *repository.EventRepository, redisClient *redis.Client) *EventService {
	return &EventService{repo: r, redisClient: redisClient}
}


func (s *EventService) CreateEvent(e model.Event) (model.Event, error) {
	err := s.repo.CreateEvent(&e)
	if err != nil {
		return e, err
	}


	s.redisClient.Del(context.Background(), "events")
	return e, nil
}


func (s *EventService) GetEvents(isActive *bool, page, limit int) ([]model.Event, error) {
	events, err := s.repo.GetEvents(isActive, page, limit)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ID’ye göre etkinlik getir
func (s *EventService) GetEventByID(id int) (*model.Event, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("event_%d", id)

	// 1. Redis'ten oku
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var e model.Event
		if err := json.Unmarshal([]byte(cached), &e); err == nil {
			return &e, nil
		}
	}
	// 2. DB'den oku
	e, err := s.repo.GetEventByID(id)
	if err != nil {
		return nil, err
	}
	// 3. Redis'e kaydet (go routine)
	go func() {
		data, _ := json.Marshal(e)
		s.redisClient.Set(ctx, cacheKey, data, 10*time.Minute)
	}()

	return e, nil
}

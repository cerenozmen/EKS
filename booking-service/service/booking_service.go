package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"booking-service/model"
	"booking-service/repository"
)

type BookingService struct {
	Repo            *repository.BookingRepository
	EventServiceURL string
}

func NewBookingService(repo *repository.BookingRepository, eventServiceURL string) *BookingService {
	return &BookingService{
		Repo:            repo,
		EventServiceURL: eventServiceURL,
	}
}

func (s *BookingService) Register(userID, eventID int) (string, error) {

	event, err := s.getEvent(eventID)
	if err != nil {
		return "", fmt.Errorf("etkinlik bilgisi alınamadı: %w", err)
	}

	if !event.IsActive {
		return "", errors.New("bu etkinlik kayıt için aktif değil")
	}

	bookedCount, err := s.Repo.CountByEventID(eventID)
	if err != nil {
		return "", fmt.Errorf("kayıt sayısı alınırken hata oluştu: %w", err)
	}

	if bookedCount >= int64(event.Capacity) {
		return "", errors.New("etkinliğin kontenjanı doldu")
	}

	booking := &model.Booking{
		UserId:    userID,
		EventId:   eventID,
		CreatedAt: time.Now(),
	}
	if err := s.Repo.Create(booking); err != nil {
		return "", err
	}

	return event.Name, nil

}

func (s *BookingService) getEvent(eventID int) (*model.Event, error) {
	url := fmt.Sprintf("%s/events/%d", s.EventServiceURL, eventID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP isteği sırasında hata oluştu: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("yanıt gövdesi okunamadı: %w", err)
	}

	// Gelen JSON'u logla
	fmt.Println("Event Service'den dönen JSON:", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Event Service hata döndürdü. Durum Kodu: %d, Mesaj: %s", resp.StatusCode, string(bodyBytes))
	}

	var event model.Event
	if err := json.Unmarshal(bodyBytes, &event); err != nil {
		return nil, fmt.Errorf("JSON çözümleme hatası: %w", err)
	}

	return &event, nil
}
func (s *BookingService) CancelByIDs(userID, eventID int) (string, error) {
	deletedEventID, err := s.Repo.DeleteByUserAndEvent(userID, eventID)
	if err != nil {
		return "", err
	}

	event, err := s.getEvent(deletedEventID)
	if err != nil {
		return "", fmt.Errorf("etkinlik bilgisi alınamadı: %w", err)
	}

	return event.Name, nil
}

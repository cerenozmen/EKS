package repository

import (
	"event-service/model"

	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}


func (r *EventRepository) CreateEvent(e *model.Event) error {
	if err := r.db.Create(e).Error; err != nil {
		return err
	}
	return nil
}


func (r *EventRepository) GetEvents(isActive *bool, page, limit int) ([]model.Event, error) {
	var events []model.Event
	query := r.db

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	offset := (page - 1) * limit

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}



func (r *EventRepository) GetEventByID(id int) (*model.Event, error) {
	var e model.Event
	if err := r.db.First(&e, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

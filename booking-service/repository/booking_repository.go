package repository

import (
	"booking-service/model"

	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) CountByEventID(eventID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.Booking{}).Where("event_id = ?", eventID).Count(&count).Error
	return count, err
}
func (r *BookingRepository) Delete(id int) (int, error) {
	var booking model.Booking
	result := r.db.First(&booking, id)
	if result.Error != nil {
		return 0, result.Error
	}
	result = r.db.Delete(&booking)
	return booking.EventId, result.Error // EventId yerine EventID kullanılmalı
}

func (r *BookingRepository) DeleteByUserAndEvent(userID, eventID int) (int, error) {
	var booking model.Booking
	result := r.db.Where("user_id = ? AND event_id = ?", userID, eventID).First(&booking)
	if result.Error != nil {
		return 0, result.Error
	}
	result = r.db.Delete(&booking)
	return booking.EventId, result.Error // EventId yerine EventID kullanılmalı
}

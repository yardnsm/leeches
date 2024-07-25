package model

import (
	"fmt"
	"time"

	"github.com/yardnsm/gohever"
	"gorm.io/gorm"
)

// Charge request states
type ChargeRequestState int

const (
	StateCreated ChargeRequestState = iota
	StateExpired

	StatePending
	StateAborted

	StateRejected
	StateApproved // This is not used though

	StateCharged
	StateFailed
)

type ChargeRequest struct {
	gorm.Model

	Amount int32
	Reason string

	CardType gohever.CardType
	State    ChargeRequestState

	RequesterID uint `gorm:"foreignKey:RequesterID"`
	Requester   User

	CachedCardStatus   gohever.CardStatus   `gorm:"serializer:json"`
	CachedCardEstimate gohever.CardEstimate `gorm:"serializer:json"`
	CachedAt           time.Time

	// We'll keep track on the messages being sent to admins / users per request
	ChargeMessages []ChargeMessage
}

type ChargeRequestsRepository struct {
	db *gorm.DB
}

func NewChargeRequestsRepository(db *gorm.DB) *ChargeRequestsRepository {
	return &ChargeRequestsRepository{db}
}

func (r *ChargeRequestsRepository) GetByID(id uint) (*ChargeRequest, error) {
	var request *ChargeRequest

	err := r.db.Where("id = ?", id).
		Preload("ChargeMessages.User").
		Preload("ChargeMessages").
		Preload("Requester").
		First(&request).Error
	if err != nil {
		return nil, fmt.Errorf("Cannot find charge request: %v", err)
	}

	return request, nil
}

func (r *ChargeRequestsRepository) Create(request *ChargeRequest) error {
	err := r.db.Create(&request).Error

	if err != nil {
		return fmt.Errorf("Cannot create charge request: %v", err)
	}

	return nil
}

func (r *ChargeRequestsRepository) Save(request *ChargeRequest) error {
	err := r.db.Save(&request).Error

	if err != nil {
		return fmt.Errorf("Unable to save charge request: %v", err)
	}

	return nil
}

package model

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	TelegramID  int64
	DisplayName string

	IsApproved bool
	IsAdmin    bool

	ChargeRequests []ChargeRequest `gorm:"foreignKey:RequesterID"`
}

// Implements tele.Recipient
func (user *User) Recipient() string {
	return fmt.Sprint(user.TelegramID)
}

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{db}
}

func (r *UsersRepository) GetByID(id uint) (*User, error) {
	var user *User

	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("Cannot find user: %v", err)
	}

	return user, nil
}

func (r *UsersRepository) GetByTelegramID(telegramID int64) (*User, error) {
	var user *User

	err := r.db.Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("Cannot find user: %v", err)
	}

	return user, nil
}

func (r *UsersRepository) GetAll() ([]User, error) {
	var users []User

	if err := r.db.Find(&users).Error; err != nil {
		return users, fmt.Errorf("Table is empty: %v", err)
	}

	return users, nil
}

func (r *UsersRepository) Create(user *User) error {
	err := r.db.Create(&user).Error

	if err != nil {
		return fmt.Errorf("Cannot create user: %v", err)
	}

	return nil
}

func (r *UsersRepository) Save(user *User) error {
	err := r.db.Save(&user).Error

	if err != nil {
		return fmt.Errorf("Unable to save user: %v", err)
	}

	return nil
}

func (r *UsersRepository) Delete(id uint) error {
	err := r.db.Delete(&User{}, id).Error

	if err != nil {
		return fmt.Errorf("Unable to delete user: %v", err)
	}

	return nil
}

func (r *UsersRepository) SetUserAdmin(id uint, isAdmin bool) (*User, error) {
	user, err := r.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Unable set user's admin state, user not found: %v", err)
	}

	user.IsAdmin = isAdmin

	err = r.Save(user)
	if err != nil {
		return nil, fmt.Errorf("Unable set user's admin state: %v", err)
	}

	return user, nil
}

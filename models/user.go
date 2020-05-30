package models

import (
	"gallery-api/rand"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const cost = 12

type User struct {
	gorm.Model
	Email    string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Token    string `gorm:"unique_index"`
}

type UserService interface {
	Create(user *User) error
	Login(user *User) (string, error)
	GetByToken(token string) (*User, error)
	Logout(token string) error
}

func NewUserService(db *gorm.DB) UserService {
	return &userGorm{db}
}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) Create(user *User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), cost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	token, err := rand.GetToken()
	if err != nil {
		return err
	}
	user.Token = token
	return ug.db.Create(user).Error
}

func (ug *userGorm) Login(user *User) (string, error) {
	found := new(User)
	err := ug.db.Where("email = ?", user.Email).First(&found).Error
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	token, err := rand.GetToken()
	if err != nil {
		return "", err
	}
	err = ug.db.Model(&User{}).
		Where("id = ?", found.ID).
		Update("token", token).Error
	if err != nil {
		return "", err
	}
	return token, nil
}

func (ug *userGorm) Logout(token string) error {
	user, err := ug.GetByToken(token)
	if err != nil {
		return err
	}
	return ug.db.Model(&User{}).
		Where("id = ?", user.ID).
		Update("token", "").Error
}

func (ug *userGorm) GetByToken(token string) (*User, error) {
	user := new(User)
	err := ug.db.Where("token = ?", token).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

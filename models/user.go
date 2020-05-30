package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"gallery-api/rand"
	"hash"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const cost = 12
const hmacKey = "secret"

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
	mac := hmac.New(sha256.New, []byte(hmacKey))
	return &userGorm{db, mac}
}

type userGorm struct {
	db   *gorm.DB
	hmac hash.Hash
}

func (ug *userGorm) Create(temp *User) error {
	user := new(User)
	user.Email = temp.Email
	user.Password = temp.Password

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), cost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	token, err := rand.GetToken()
	if err != nil {
		return err
	}

	fmt.Println("token ===> ", token)
	ug.hmac.Write([]byte(token))
	tokenHash := ug.hmac.Sum(nil)
	ug.hmac.Reset()
	tokenHashStr := base64.URLEncoding.EncodeToString(tokenHash)
	fmt.Println("tokenHashStr ===> ", tokenHashStr)

	user.Token = tokenHashStr
	temp.Token = token

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

	fmt.Println("token ===> ", token)

	ug.hmac.Write([]byte(token))
	tokenHash := ug.hmac.Sum(nil)
	ug.hmac.Reset()
	tokenHashStr := base64.URLEncoding.EncodeToString(tokenHash)
	fmt.Println("tokenHashStr ===> ", tokenHashStr)

	err = ug.db.Model(&User{}).
		Where("id = ?", found.ID).
		Update("token", tokenHashStr).Error
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
	ug.hmac.Write([]byte(token))
	tokenHash := ug.hmac.Sum(nil)
	ug.hmac.Reset()
	tokenHashStr := base64.URLEncoding.EncodeToString(tokenHash)

	user := new(User)
	err := ug.db.Where("token = ?", tokenHashStr).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

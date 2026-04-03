package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id       int     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserUuid string  `json:"user_uuid" gorm:"type:uuid"`
	Password *string `json:"password"`
}

var UserAlreadyExistsErr = errors.New("user already exist")

var UserNotFoundErr = errors.New("user not found")

func passwordHash(password string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(passHash), nil
}

func CreateUser(db *gorm.DB, u *User) error {
	var existingUser User
	err := db.Where("user_uuid = ?", u.UserUuid).First(&existingUser).Error
	if err == nil && existingUser.UserUuid == u.UserUuid {
		return UserAlreadyExistsErr
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	password, err := passwordHash(*u.Password)
	if err != nil {
		return err
	}

	u.Password = &password
	return db.Create(u).Error
}

func GetUserByUuidAndPass(db *gorm.DB, uuid, password string) (User, error) {
	var user User

	result := db.Where("user_uuid = ?", uuid).First(&user)
	err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, UserNotFoundErr
	}

	return user, result.Error
}

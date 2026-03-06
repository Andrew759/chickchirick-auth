package user

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Phone      string    `json:"phone" gorm:"type:varchar(30)"`
	Name       string    `json:"name" gorm:"type:varchar(256);not null"`
	Surname    string    `json:"surname" gorm:"type:varchar(256);not null"`
	Login      string    `json:"login" gorm:"type:varchar(256);unique; not null"`
	Property   *Property `json:"property" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt  time.TimestampWithTimeZoneMicro
	UpdatedAt  time.TimestampWithTimeZoneMicro
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type UserAlreadyExistErr struct {
	Field string
	Value string
}

func (e *UserAlreadyExistErr) Error() string {
	return fmt.Sprintf("user with %s '%s' already exists", e.Field, e.Value)
}

var UserNotFoundErr = errors.New("user not found")

func CreateUser(db *gorm.DB, u *User) error {
	var existing User
	err := db.Where("login = ? OR phone = ?", u.Login, u.Phone).First(&existing).Error

	if err == nil {
		if existing.Login == u.Login {
			return &UserAlreadyExistErr{Field: "login", Value: u.Login}
		}
		if existing.Phone == u.Phone {
			return &UserAlreadyExistErr{Field: "phone", Value: u.Phone}
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return db.Create(u).Error

}

func UpdateUserById(db *gorm.DB, u *User, id int) error {
	var user User
	result := db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return UserNotFoundErr
	}

	return db.Save(u).Error
}

func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)

	return users, result.Error
}

func GetUserById(db *gorm.DB, id int) (User, error) {
	var user User
	result := db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, UserNotFoundErr
	}

	return user, result.Error
}

func DeleteUserById(db *gorm.DB, id int) error {
	var user User
	result := db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return UserNotFoundErr
	}

	return db.Delete(&User{}, id).Error
}

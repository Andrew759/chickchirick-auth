package user

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	FileUuid   uuid.UUID `gorm:"type:uuid, not null"`
	UserId     int       `json:"user_id" gorm:"type:int;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User       User      `json:"user" gorm:"foreignKey:UserId;references:Id"`
	CreatedAt  time.TimestampWithTimeZoneMicro
	UpdatedAt  time.TimestampWithTimeZoneMicro
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

var PhotoNotFoundErr = errors.New("user not found")

func CreatePhoto(db *gorm.DB, b *Photo) error {
	return db.Create(b).Error
}

func UpdatePhotoById(db *gorm.DB, p *Photo, id int) error {
	var photo Photo
	result := db.First(&photo, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return PhotoNotFoundErr
	}

	return db.Save(p).Error
}

func GetPhotos(db *gorm.DB) ([]Photo, error) {
	var photos []Photo
	result := db.Find(&photos)

	return photos, result.Error
}

func GetPhotoById(db *gorm.DB, id int) (Photo, error) {
	var photo Photo
	result := db.First(&photo, id)

	return photo, result.Error
}

func GetPhotosByUserId(db *gorm.DB, userId int) ([]Photo, error) {
	var photos []Photo
	result := db.Where("user_id = ?", userId).Find(&photos)

	return photos, result.Error
}

func DeletePhotoById(db *gorm.DB, id int) error {
	var photo Photo
	result := db.First(&photo, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return PhotoNotFoundErr
	}

	return db.Delete(&Photo{}, id).Error
}

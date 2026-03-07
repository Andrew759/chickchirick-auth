package auth

import (
	"chickchirick-auth/pkg/chirik_gorm_tweaks/time"
	"errors"

	"gorm.io/gorm"
)

type Code struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int                             `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Code       int8                            `json:"code" gorm:"type:smallint"`
	SessionId  int                             `json:"session_id" gorm:"type:int"`
	Session    Session                         `json:"session" gorm:"references:SessionId"`
	ExpiresAt  time.TimestampWithTimeZoneMicro `json:"expires_at" gorm:"type:timestamp without time zone"`
	CreatedAt  time.TimestampWithTimeZoneMicro
	UpdatedAt  time.TimestampWithTimeZoneMicro
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

var CodeNotFoundErr = errors.New("code not found")

func CreateCode(db *gorm.DB, c *Code) error {
	return db.Create(c).Error
}

func UpdateCodeById(db *gorm.DB, c *Code, id int) error {
	var code Code
	result := db.First(&code, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return CodeNotFoundErr
	}

	return db.Save(c).Error
}

func GetCodes(db *gorm.DB) ([]Code, error) {
	var codes []Code
	result := db.Find(&codes)

	return codes, result.Error
}

func GetCodeById(db *gorm.DB, id int) (Code, error) {
	var codes Code
	result := db.First(&codes, id)

	return codes, result.Error
}

func DeleteCodeById(db *gorm.DB, id int) error {
	var code Code
	result := db.First(&code, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return CodeNotFoundErr
	}

	return db.Delete(&Code{}, id).Error
}

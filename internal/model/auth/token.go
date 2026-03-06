package auth

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"errors"
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int                             `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	SessionId  int                             `json:"session_id" gorm:"type:int"`
	Session    Session                         `json:"session" gorm:"references:SessionId"`
	Token      string                          `json:"token" gorm:"type:varchar(256)"`
	ExpiresAt  time.TimestampWithTimeZoneMicro `json:"expires_at" gorm:"type:timestamp without time zone"`
	CreatedAt  time.TimestampWithTimeZoneMicro
	UpdatedAt  time.TimestampWithTimeZoneMicro
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

var TokenNotFoundErr = errors.New("token not found")

func CreateToken(db *gorm.DB, t *Token) error {
	return db.Create(t).Error
}

func UpdateTokenById(db *gorm.DB, t *Token, id int) error {
	var token Token
	result := db.First(&token, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return TokenNotFoundErr
	}

	return db.Save(t).Error
}

func GetTokens(db *gorm.DB) ([]Token, error) {
	var tokens []Token
	result := db.Find(&tokens)

	return tokens, result.Error
}

func GetTokenById(db *gorm.DB, id int) (Token, error) {
	var token Token
	result := db.First(&token, id)

	return token, result.Error
}

func DeleteTokenById(db *gorm.DB, id int) error {
	var token Token
	result := db.First(&token, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return TokenNotFoundErr
	}

	return db.Delete(&Token{}, id).Error
}

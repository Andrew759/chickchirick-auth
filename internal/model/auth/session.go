package auth

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int                             `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserUuid   uuid.UUID                       `json:"user_uuid" gorm:"type:uuid"`
	Status     bool                            `json:"status" gorm:"type:boolean"`
	StartDate  time.TimestampWithTimeZoneMicro `json:"start_date" gorm:"type:timestamp without time zone"`
	CreatedAt  time.TimestampWithTimeZoneMicro
	UpdatedAt  time.TimestampWithTimeZoneMicro
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

var SessionNotFoundErr = errors.New("session not found")

func CreateSession(db *gorm.DB, s *Session) error {
	return db.Create(s).Error
}

func UpdateSessionById(db *gorm.DB, s *Session, id int) error {
	var session Session
	result := db.First(&session, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return SessionNotFoundErr
	}

	return db.Save(s).Error
}

func GetSessions(db *gorm.DB) ([]Session, error) {
	var sessions []Session
	result := db.Find(&sessions)

	return sessions, result.Error
}

func GetSessionById(db *gorm.DB, id int) (Session, error) {
	var session Session
	result := db.First(&session, id)

	return session, result.Error
}

func DeleteSessionById(db *gorm.DB, id int) error {
	var session Session
	result := db.First(&session, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return SessionNotFoundErr
	}

	return db.Delete(&Session{}, id).Error
}

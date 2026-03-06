package user

import (
	"chickChirick/internal/model/c_model"
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"errors"

	"gorm.io/gorm"
)

type Property struct {
	//TODO: нужно указывать c_migrator_t_name сразу в двух местах при использовании GORM. Переделать
	c_model.Model `c_migrator:"enabled" c_migrator_t_name:"properties"`
	UserId        int `json:"user_id" gorm:"primaryKey;unique;not null"`
	//TODO: структура пользователя не должна попадать в ответы методов properties
	User      User    `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Timezone  int16   `json:"timezone" gorm:"type:smallint;default:3"`
	Email     string  `json:"email" gorm:"unique;type:varchar(256)"`
	Password  *string `json:"password" gorm:"type:varchar(1024)"`
	CreatedAt time.TimestampWithTimeZoneMicro
	UpdatedAt time.TimestampWithTimeZoneMicro
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

var PropertyNotFoundErr = errors.New("property not found")

var PropertyForUserAlreadyExistsErr = errors.New("property for user already exists")

func (p *Property) TableName() string {
	return "properties"
}

func (p *Property) SetPassword(password string) {
	p.Password = &password
}

func CreateProperty(db *gorm.DB, p *Property) error {
	//TODO: тут баг, потому что у версии таблицы properties из самописного мигратора - нет ненужного поля id
	return db.Create(p).Error
}

func UpdatePropertyByUserId(db *gorm.DB, p *Property, userId int) error {
	//TODO: тут баг. Из-за того, что емейл проверяется на стадии валидации
	var property Property
	result := db.Model(&Property{}).Where("user_id = ?", userId).Take(&property)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return PropertyNotFoundErr
	}

	return db.Save(p).Error
}

func GetProperties(db *gorm.DB) ([]Property, error) {
	var properties []Property
	result := db.Find(&properties)

	return properties, result.Error
}

func GetPropertyByUserId(db *gorm.DB, userId int) (Property, error) {
	var property Property
	result := db.Model(&Property{}).Where("user_id = ?", userId).Take(&property)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return property, PropertyNotFoundErr
	}

	return property, result.Error
}

func GetPropertyToAnotherUserByEmail(db *gorm.DB, userId int, email string) (Property, error) {
	var property Property
	result := db.Model(&Property{}).Where("user_id != ? AND email = ?", userId, email).Take(&property)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return property, PropertyNotFoundErr
	}

	return property, result.Error
}

func HasProperty(db *gorm.DB, p Property) (bool, error) {
	var count int64
	err := db.Model(&Property{}).Where("user_id = ? OR email = ?", p.UserId, p.Email).Count(&count).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) && count == 0 {
		return false, nil
	}
	if count > 0 {
		return true, PropertyForUserAlreadyExistsErr
	}

	return false, err
}

func DeletePropertyByUserId(db *gorm.DB, userId int) error {
	result := db.Where("user_id = ?", userId).Delete(&Property{})

	if result.Error != nil {
		return result.Error
	}

	// Если ни одна запись не была затронута (удалена),
	// значит свойств у этого пользователя не было
	if result.RowsAffected == 0 {
		return PropertyNotFoundErr
	}

	return nil
}

package service

import (
	"chickchirick-auth/cmd/config/dto"
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

// DBDecorator TODO: отрефакторить. Выбрать имя без постфикса интерфейс. Сделать также во всём приложении
type DBDecorator struct {
	GormInterface   *gorm.DB
	NativeInterface *sql.DB
}

func InitORM(config dto.DataBaseConfigInterface) DBDecorator {
	dsn := dsn(config)

	ORM, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("db connect failed: %w", err))
	}

	nativeDB, err := ORM.DB()
	if err != nil {
		panic(fmt.Errorf("error receiving the native interface: %w", err))
	}

	dbd := DBDecorator{
		GormInterface:   ORM,
		NativeInterface: nativeDB,
	}

	return dbd
}

func dsn(config dto.DataBaseConfigInterface) string {
	dsn := []string{
		"host=" + config.Host(),
		"user=" + config.User(),
		"password=" + config.Password(),
		"dbname=" + config.Name(),
		"port=" + strconv.Itoa(config.Port()),
	}
	if config.Timezone() != "" {
		dsn = append(dsn, "TimeZone="+config.Timezone())
	}

	return strings.Join(dsn, " ")
}

func (dbd DBDecorator) CloseDB() {
	err := dbd.NativeInterface.Close()
	if err != nil {
		panic(fmt.Errorf("db close error: %w", err))
	}
}

func (dbd DBDecorator) GDB() *gorm.DB {
	return dbd.GormInterface
}

func (dbd DBDecorator) NativeDB() *sql.DB {
	return dbd.NativeInterface
}

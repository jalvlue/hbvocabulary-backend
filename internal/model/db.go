package model

import (
	"HBVocabulary/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	dsn := config.Conf.DBSource
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

package model

import "gorm.io/gorm"

type TestRecord struct {
	gorm.Model
	UserID string `json:"user_id"`
	Socre  int    `json:"score"`
}

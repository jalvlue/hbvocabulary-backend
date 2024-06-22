package model

import (
	"time"
)

type User struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	TestCount int    `json:"test_count"`
	MaxScore  int    `json:"max_score"`
	CreatedAt time.Time
}

func (User) TableName() string {
	return "t_users"
}

func (s *Store) CreateUser(u *User) error {
	result := s.db.Model(&User{}).Create(u)
	return result.Error
}

func (s *Store) GetUserByUsername(username string) (*User, error) {
	u := &User{}
	result := s.db.Model(&User{}).Where("username = ?", username).First(u)
	return u, result.Error
}

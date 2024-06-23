package model

type Vocabulary struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Word string `json:"word"`
}

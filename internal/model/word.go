package model

type Vocabulary struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Word string `json:"word"`
}

func (Vocabulary) TableName() string {
	return "t_vocabularies"
}

func (s *Store) GetWordListById(wordIDs []int) ([]Vocabulary, error) {
	words := []Vocabulary{}
	result := s.db.Model(&Vocabulary{}).Select("id, word").Where("id in ?", wordIDs).Find(&words)
	return words, result.Error
}

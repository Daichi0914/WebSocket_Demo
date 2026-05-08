package infrastructure

import (
	"backend/domain"

	"gorm.io/gorm"
)

type mysqlMessageRepository struct {
	db *gorm.DB
}

// NewMysqlMessageRepository creates a new mysql repository
func NewMysqlMessageRepository(db *gorm.DB) domain.MessageRepository {
	// Auto migrate the domain model
	db.AutoMigrate(&domain.Message{})
	return &mysqlMessageRepository{db: db}
}

func (m *mysqlMessageRepository) Store(message *domain.Message) error {
	return m.db.Create(message).Error
}

func (m *mysqlMessageRepository) FetchRecent(limit int) ([]domain.Message, error) {
	var messages []domain.Message
	err := m.db.Order("created_at asc").Limit(limit).Find(&messages).Error
	return messages, err
}

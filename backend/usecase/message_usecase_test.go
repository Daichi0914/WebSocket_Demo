package usecase_test

import (
	"errors"
	"testing"
	"backend/domain"
	"backend/usecase"

	"github.com/stretchr/testify/assert"
)

// --- Mocks ---

type mockMessageRepository struct {
	StoreFunc       func(message *domain.Message) error
	FetchRecentFunc func(limit int) ([]domain.Message, error)
}

func (m *mockMessageRepository) Store(message *domain.Message) error {
	if m.StoreFunc != nil {
		return m.StoreFunc(message)
	}
	return nil
}

func (m *mockMessageRepository) FetchRecent(limit int) ([]domain.Message, error) {
	if m.FetchRecentFunc != nil {
		return m.FetchRecentFunc(limit)
	}
	return nil, nil
}

type mockPubSub struct {
	PublishFunc   func(channel string, message []byte) error
	SubscribeFunc func(channel string, handler func(payload []byte)) error
}

func (m *mockPubSub) Publish(channel string, message []byte) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(channel, message)
	}
	return nil
}

func (m *mockPubSub) Subscribe(channel string, handler func(payload []byte)) error {
	if m.SubscribeFunc != nil {
		return m.SubscribeFunc(channel, handler)
	}
	return nil
}

// --- Tests ---

func TestMessageUsecase_SaveAndPublishMessage(t *testing.T) {
	t.Run("正常系: メッセージが正しく保存され、Pub/Subに配信されること", func(t *testing.T) {
		// Given: 正常に動作するモックリポジトリとPubSubを用意する
		storeCalled := false
		publishCalled := false

		mockRepo := &mockMessageRepository{
			StoreFunc: func(message *domain.Message) error {
				storeCalled = true
				assert.Equal(t, "Alice", message.Sender, "送信者が一致しません")
				assert.Equal(t, "Hello", message.Content, "コンテンツが一致しません")
				return nil
			},
		}

		mockPub := &mockPubSub{
			PublishFunc: func(channel string, message []byte) error {
				publishCalled = true
				assert.Equal(t, "chat_channel", channel, "チャンネル名が一致しません")
				return nil
			},
		}

		u := usecase.NewMessageUsecase(mockRepo, mockPub)

		// When: メッセージの保存と公開を実行する
		err := u.SaveAndPublishMessage("Alice", "Hello")

		// Then: エラーなく終了し、StoreとPublishが呼ばれていることを確認する
		assert.NoError(t, err)
		assert.True(t, storeCalled, "Storeが呼ばれていません")
		assert.True(t, publishCalled, "Publishが呼ばれていません")
	})

	t.Run("異常系: リポジトリの保存に失敗した場合、エラーが返ること", func(t *testing.T) {
		// Given: データベース保存時にエラーを返すモックリポジトリを用意する
		expectedErr := errors.New("db error")
		mockRepo := &mockMessageRepository{
			StoreFunc: func(message *domain.Message) error {
				return expectedErr
			},
		}
		mockPub := &mockPubSub{
			PublishFunc: func(channel string, message []byte) error {
				assert.Fail(t, "Storeが失敗した場合はPublishは呼ばれないはずです")
				return nil
			},
		}

		u := usecase.NewMessageUsecase(mockRepo, mockPub)

		// When: メッセージの保存と公開を実行する
		err := u.SaveAndPublishMessage("Bob", "Fail message")

		// Then: 期待したエラーが返却されることを確認する
		assert.ErrorIs(t, err, expectedErr)
	})
}

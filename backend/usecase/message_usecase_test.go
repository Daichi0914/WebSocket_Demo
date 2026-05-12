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

	t.Run("異常系: Publishに失敗した場合、エラーが返ること", func(t *testing.T) {
		expectedErr := errors.New("redis publish error")
		mockRepo := &mockMessageRepository{
			StoreFunc: func(message *domain.Message) error { return nil },
		}
		mockPub := &mockPubSub{
			PublishFunc: func(channel string, message []byte) error {
				return expectedErr
			},
		}

		u := usecase.NewMessageUsecase(mockRepo, mockPub)
		err := u.SaveAndPublishMessage("Alice", "Hello")

		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestMessageUsecase_GetRecentMessages(t *testing.T) {
	t.Run("正常系: リポジトリからメッセージが取得できること", func(t *testing.T) {
		expectedMsgs := []domain.Message{
			{Sender: "Alice", Content: "Hello"},
			{Sender: "Bob", Content: "Hi"},
		}
		mockRepo := &mockMessageRepository{
			FetchRecentFunc: func(limit int) ([]domain.Message, error) {
				assert.Equal(t, 10, limit)
				return expectedMsgs, nil
			},
		}

		u := usecase.NewMessageUsecase(mockRepo, nil)
		msgs, err := u.GetRecentMessages(10)

		assert.NoError(t, err)
		assert.Equal(t, expectedMsgs, msgs)
	})

	t.Run("異常系: リポジトリ取得に失敗した場合、エラーが返ること", func(t *testing.T) {
		expectedErr := errors.New("db fetch error")
		mockRepo := &mockMessageRepository{
			FetchRecentFunc: func(limit int) ([]domain.Message, error) {
				return nil, expectedErr
			},
		}

		u := usecase.NewMessageUsecase(mockRepo, nil)
		_, err := u.GetRecentMessages(10)

		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestMessageUsecase_ListenToMessages(t *testing.T) {
	t.Run("正常系: 購読に成功し、メッセージがハンドラーに渡されること", func(t *testing.T) {
		mockPub := &mockPubSub{
			SubscribeFunc: func(channel string, handler func(payload []byte)) error {
				// モック内でハンドラーを即座に実行する
				payload := `{"sender":"Alice","content":"Hello"}`
				handler([]byte(payload))
				return nil
			},
		}

		u := usecase.NewMessageUsecase(nil, mockPub)
		
		called := false
		err := u.ListenToMessages(func(msg domain.Message) {
			called = true
			assert.Equal(t, "Alice", msg.Sender)
			assert.Equal(t, "Hello", msg.Content)
		})

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("異常系: 購読の開始自体に失敗した場合、エラーが返ること", func(t *testing.T) {
		expectedErr := errors.New("subscribe error")
		mockPub := &mockPubSub{
			SubscribeFunc: func(channel string, handler func(payload []byte)) error {
				return expectedErr
			},
		}

		u := usecase.NewMessageUsecase(nil, mockPub)
		err := u.ListenToMessages(func(msg domain.Message) {})

		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("準正常系: 届いたデータが不正なJSONの場合、ハンドラーが呼ばれないこと", func(t *testing.T) {
		mockPub := &mockPubSub{
			SubscribeFunc: func(channel string, handler func(payload []byte)) error {
				handler([]byte("invalid json"))
				return nil
			},
		}

		u := usecase.NewMessageUsecase(nil, mockPub)
		
		called := false
		err := u.ListenToMessages(func(msg domain.Message) {
			called = true
		})

		assert.NoError(t, err)
		assert.False(t, called, "不正なデータでハンドラーが呼ばれてはいけません")
	})
}

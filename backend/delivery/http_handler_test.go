package delivery_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"backend/delivery"
	"backend/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mocks ---
type mockMessageUsecase struct {
	GetRecentMessagesFunc     func(limit int) ([]domain.Message, error)
	SaveAndPublishMessageFunc func(sender, content string) error
	ListenToMessagesFunc      func(handler func(msg domain.Message)) error
}

func (m *mockMessageUsecase) GetRecentMessages(limit int) ([]domain.Message, error) {
	if m.GetRecentMessagesFunc != nil {
		return m.GetRecentMessagesFunc(limit)
	}
	return nil, nil
}

func (m *mockMessageUsecase) SaveAndPublishMessage(sender, content string) error {
	return nil
}

func (m *mockMessageUsecase) ListenToMessages(handler func(msg domain.Message)) error {
	return nil
}

// --- Tests ---
func TestHTTPHandler_GetMessages(t *testing.T) {
	t.Run("正常系: メッセージ一覧が正しく取得できること", func(t *testing.T) {
		// Given: 2件のメッセージを返すUsecaseのモックを作成する
		mockMessages := []domain.Message{
			{Sender: "Alice", Content: "Hello"},
			{Sender: "Bob", Content: "Hi"},
		}
		mockUsecase := &mockMessageUsecase{
			GetRecentMessagesFunc: func(limit int) ([]domain.Message, error) {
				return mockMessages, nil
			},
		}

		handler := delivery.NewHTTPHandler(mockUsecase)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux)

		req, err := http.NewRequest(http.MethodGet, "/api/messages", nil)
		require.NoError(t, err)

		// When: エンドポイントにGETリクエストを送信する
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		// Then: ステータスコードが200で、期待したJSONデータが返却されること
		assert.Equal(t, http.StatusOK, rr.Code)

		var response []domain.Message
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err, "JSONのデコードに失敗しました")

		assert.Len(t, response, 2, "メッセージ数が一致しません")
		assert.Equal(t, "Alice", response[0].Sender, "送信者が一致しません")
	})

	t.Run("異常系: GET以外のメソッドを許可しないこと", func(t *testing.T) {
		// Given: Usecaseのモックを作成し、ハンドラーを準備する
		mockUsecase := &mockMessageUsecase{}
		handler := delivery.NewHTTPHandler(mockUsecase)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux)

		req, err := http.NewRequest(http.MethodPost, "/api/messages", nil)
		require.NoError(t, err)

		// When: POSTリクエストを送信する
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		// Then: 405 Method Not Allowed が返却されること
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})
}

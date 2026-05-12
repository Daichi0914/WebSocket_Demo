//go:build integration

package infrastructure_test

import (
	"context"
	"log"
	"os"
	"testing"

	"backend/domain"
	"backend/infrastructure"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	// Podman環境でのエラー（bridgeネットワークがない問題）を回避するためRyukを無効化
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	ctx := context.Background()

	log.Println("Starting MySQL test container...")
	
	// コンテナの起動
	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithDatabase("test_db"),
		mysql.WithUsername("test_user"),
		mysql.WithPassword("test_password"),
	)
	if err != nil {
		log.Fatalf("Failed to start container: %s", err)
	}

	// テスト終了時にコンテナを破棄
	defer func() {
		log.Println("Terminating MySQL test container...")
		if err := mysqlContainer.Terminate(ctx); err != nil {
			log.Fatalf("Failed to terminate container: %s", err)
		}
	}()

	// 接続文字列を取得
	connStr, err := mysqlContainer.ConnectionString(ctx, "parseTime=true&loc=Local")
	if err != nil {
		log.Fatalf("Failed to get connection string: %s", err)
	}

	// GORMで接続
	db, err = gorm.Open(gorm_mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to db: %s", err)
	}

	// テスト群を実行
	m.Run()
}

func TestMysqlMessageRepository(t *testing.T) {
	// リポジトリの初期化（AutoMigrateが実行される）
	repo := infrastructure.NewMysqlMessageRepository(db)

	t.Run("正常系: メッセージをDBに保存できること", func(t *testing.T) {
		// 事前にテーブルをクリーンアップ
		db.Exec("DELETE FROM messages")

		msg := &domain.Message{
			Sender:  "Alice",
			Content: "Test message",
		}

		// When: 保存を実行
		err := repo.Store(msg)
		require.NoError(t, err)

		// Then: DBに1件保存されていること
		var count int64
		db.Model(&domain.Message{}).Count(&count)
		assert.Equal(t, int64(1), count)

		// Then: 内容が正しく保存されていること
		var savedMsg domain.Message
		db.First(&savedMsg)
		assert.Equal(t, "Alice", savedMsg.Sender)
		assert.Equal(t, "Test message", savedMsg.Content)
	})

	t.Run("正常系: 最新のメッセージを上限付きで取得できること", func(t *testing.T) {
		// 事前にテーブルをクリーンアップ
		db.Exec("DELETE FROM messages")

		// Given: 3件のメッセージを保存
		repo.Store(&domain.Message{Sender: "User1", Content: "First"})
		repo.Store(&domain.Message{Sender: "User2", Content: "Second"})
		repo.Store(&domain.Message{Sender: "User3", Content: "Third"})

		// When: 上限2で取得
		messages, err := repo.FetchRecent(2)
		require.NoError(t, err)

		// Then: 2件だけ取得できること（古い順）
		assert.Len(t, messages, 2)
		assert.Equal(t, "First", messages[0].Content)
		assert.Equal(t, "Second", messages[1].Content)
	})
}

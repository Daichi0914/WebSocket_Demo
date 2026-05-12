.PHONY: up down build rebuild logs ps restart clean e2e-up e2e-down test-backend test-integration test-frontend test-e2e stg-up stg-down stg-build prod-up prod-down prod-build

# 起動 (バックグラウンド)
up:
	podman compose --env-file .env.dev up -d

# 停止
down:
	podman compose --env-file .env.dev down

# コンテナのビルド
build:
	podman compose --env-file .env.dev build

# 再ビルドして起動
rebuild:
	podman compose --env-file .env.dev up -d --build

# ログを表示
logs:
	podman compose --env-file .env.dev logs -f

# コンテナの状態を確認
ps:
	podman compose --env-file .env.dev ps

# 再起動
restart:
	podman compose --env-file .env.dev restart

# 停止してボリュームごと削除
clean:
	podman compose --env-file .env.dev down -v

# --- Testing ---

# バックエンドの単体テスト
test-backend:
	cd backend && go test -v ./...

# バックエンドの結合テスト（DBコンテナを使用）
test-integration:
	cd backend && go test -v -tags=integration ./...

# フロントエンドの単体テスト
test-frontend:
	cd frontend && npm run test

# E2Eテスト（Playwright）
# ※バックエンドが起動している必要があります（make up）
test-e2e:
	$(MAKE) e2e-up
	@echo "Waiting for E2E backend to be ready..."
	@for i in $$(seq 1 30); do \
		if curl -s http://localhost:8081/api/messages > /dev/null 2>&1; then \
			echo "E2E backend is ready!"; \
			break; \
		fi; \
		echo "Waiting... ($$i/30)"; \
		sleep 2; \
		if [ $$i -eq 30 ]; then echo "Timeout waiting for E2E backend"; exit 1; fi; \
	done
	@cd frontend && E2E_ENV=true NEXT_PUBLIC_API_URL=http://localhost:8081 NEXT_PUBLIC_WS_URL=ws://localhost:8081/ws npm run test:e2e; \
	EXIT_CODE=$$?; $(MAKE) -C .. e2e-down; exit $$EXIT_CODE

# --- E2E Environment ---

# E2Eテスト用のバックエンド環境を起動（開発環境と別ポートで起動）
e2e-up:
	podman compose -p websocket-chat-e2e --env-file .env.e2e --profile e2e up -d

# E2Eテスト用のバックエンド環境を停止・ボリュームごと削除
e2e-down:
	podman compose -p websocket-chat-e2e --env-file .env.e2e --profile e2e down -v
# --- Staging Environment ---

# ステージング環境の起動
stg-up:
	podman compose -p chat-stg -f compose.prod.yaml --env-file .env.stg up -d

# ステージング環境の停止
stg-down:
	podman compose -p chat-stg -f compose.prod.yaml --env-file .env.stg down

# ステージング環境のビルド
stg-build:
	podman compose -p chat-stg -f compose.prod.yaml --env-file .env.stg build

# --- Production Environment ---

# 本番環境の起動 (おうちサーバー向け)
prod-up:
	podman compose -p chat-prod -f compose.prod.yaml --env-file .env.prod up -d

# 本番環境の停止
prod-down:
	podman compose -p chat-prod -f compose.prod.yaml --env-file .env.prod down

# 本番環境のビルド
prod-build:
	podman compose -p chat-prod -f compose.prod.yaml --env-file .env.prod build

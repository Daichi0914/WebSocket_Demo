.PHONY: up down build logs ps restart clean

# 起動 (バックグラウンド)
up:
	podman compose up -d

# 停止
down:
	podman compose down

# コンテナのビルド
build:
	podman compose build

# 再ビルドして起動
rebuild:
	podman compose up -d --build

# ログを表示
logs:
	podman compose logs -f

# コンテナの状態を確認
ps:
	podman compose ps

# 再起動
restart:
	podman compose restart

# 停止してボリュームごと削除
clean:
	podman compose down -v

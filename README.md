# Cosmic Chat (WebSocket Demo)

リアルタイムなチャット体験を提供する、WebSocketを活用したデモアプリケーションです。

## 概要・特徴

- **リアルタイムメッセージング**: WebSocketを利用し、遅延のないチャットを実現しています。
- **アカウント登録不要**: 画面を開いて好きな名前（表示名）を入力するだけですぐにチャットに参加できます。名前はブラウザの一時メモリ（React State）に保持される「仮のアイデンティティ」として扱われます。
- **履歴の永続化**: 過去のメッセージはMySQLに保存されるため、新しく参加した際にもこれまでのチャット履歴を閲覧できます。
- **スケールを意識した設計**: RedisのPub/Sub機能を利用してメッセージをブロードキャストしているため、バックエンドのGoサーバーが複数台にスケールアウトしてもチャットが成立するアーキテクチャになっています。

## 技術スタック

### Frontend
- **Framework**: Next.js (App Router)
- **Library**: React 19, TypeScript
- **Styling**: Vanilla CSS

### Backend
- **Language**: Go 1.24
- **Libraries**: Gorilla WebSocket, GORM
- **Architecture**: クリーンアーキテクチャを意識した層次構造

### Infrastructure (Docker / Podman)
- **Database**: MySQL 8.0
- **Cache / Message Broker**: Redis 7 (Alpine)

## ローカルでの動かし方

プロジェクトルートにある `Makefile` を使うことで、簡単にPodman（またはDocker）コンテナを操作できます。

```bash
# コンテナの起動（バックグラウンド）
make up

# コンテナの停止
make down

# コンテナの再ビルドと起動（設定変更時など）
make rebuild

# ログの確認
make logs

# コンテナの状態確認
make ps

# 完全にリセットする（DBのデータ等ボリュームも削除）
make clean
```

## メッセージの流れ（内部動作）

1. ユーザーがフロントエンド（Next.js）で名前を入力し、バックエンド（Go）とWebSocket接続を確立。
2. ユーザーがメッセージを送信すると、WebSocketを通じてバックエンドへJSONが送られる。
3. バックエンドはGORMを使ってメッセージをMySQLに永続化する。
4. 保存後、バックエンドはRedisのPub/SubにメッセージをPublishする。
5. Redisからメッセージを受信（Subscribe）した全バックエンドインスタンスは、自身に接続されている全WebSocketクライアントへメッセージをブロードキャストする。
6. フロントエンド側でメッセージを受け取り、ReactのStateが更新され画面にチャットが描画される。

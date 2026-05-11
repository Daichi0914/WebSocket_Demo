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


# WebSocket 学習ガイド

今回のチャットアプリでWebSocketがどのように動いているのか、基礎概念から実際のコードまで事細かに解説します。

---

## 1. WebSocketとは？（HTTPとの違い）

普段のWebサイト（HTTP）は、**「クライアント（ブラウザ）が質問して、サーバーが答える」** という一問一答の方式です。サーバーから突然「新しい情報があるよ！」と話しかけることはできません。

しかし、チャットのように「誰かが発言したら、**リロードしなくても**全員の画面に即座に表示させたい」場合、HTTPでは不便です。
そこで登場するのが **WebSocket** です。

WebSocketは、一度電話を繋いだら **「つなぎっぱなし（常時接続）」** にして、ブラウザとサーバーの **どちらからでも、いつでも** データを送り合える（双方向通信）技術です。

---

## 2. フロントエンドの実装（React側）

フロントエンド（`frontend/src/app/page.tsx`）がWebSocketとどうやり取りしているかを見てみましょう。

### ① 接続する（電話をかける）
```typescript
// 接続先のURL（例: ws://localhost:8080/ws）を指定して接続
const socket = new WebSocket(WS_URL);

// 接続が成功した時に呼ばれる処理
socket.onopen = () => {
  setIsConnected(true);
  console.log("Connected to WS");
};
```
`new WebSocket()` を呼ぶだけで、サーバーに対して「WebSocketで繋ごう！」と要求（ハンドシェイク）を送ります。

### ② メッセージを受け取る（受信）
```typescript
// サーバーからデータが送られてきた時に呼ばれる処理
socket.onmessage = (event) => {
  // 送られてきた文字列（JSON）をオブジェクトに変換
  const msg: Message = JSON.parse(event.data);
  // チャット履歴（ReactのState）の末尾に追加する
  setMessages(prev => [...prev, msg]);
};
```
つなぎっぱなしになっている間、サーバーからデータが降ってくると自動的に `onmessage` が発火します。ここで画面を更新します。

### ③ メッセージを送る（送信）
```typescript
const msg = {
  sender: "あなたの名前",
  content: "こんにちは！"
};
// JSON文字列に変換してサーバーに送信
ws.current.send(JSON.stringify(msg));
```
ユーザーが送信ボタンを押した時、単に `send()` を使ってサーバーにデータを投げます。

---

## 3. バックエンドの実装（Go側）

バックエンド（`backend/delivery/ws_handler.go`）は、少し複雑です。「何人ものクライアントとの電話を同時に繋いでおく」必要があるからです。

### ① HTTPからWebSocketへの「アップグレード」
最初にクライアントが接続してくる時は、実は普通の「HTTPアクセス」として来ます。それをWebSocketに切り替える（アップグレードする）必要があります。

```go
func (h *WebSocketHandler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
    // HTTPリクエストをWebSocket接続に昇格させる！
    conn, err := h.upgrader.Upgrade(w, r, nil)
    // ...
}
```

### ② クライアント（接続者）の管理
サーバーは、「現在誰が接続しているか」を覚えておき、誰かが発言したら全員に配る必要があります。

```go
type WebSocketHandler struct {
    // 繋がっているすべてのクライアントを記録する名簿
    clients   map[*websocket.Conn]bool
    clientsMu sync.Mutex // 複数人同時アクセス時の競合を防ぐためのロック
}

// 接続されたら名簿に追加
h.clientsMu.Lock()
h.clients[conn] = true
h.clientsMu.Unlock()
```
`map` というデータ構造を使って、現在の接続（`conn`）を名簿に登録しています。

### ③ 無限ループで待ち受ける
電話が繋がっている間は、相手がいつ喋り出すかわかりません。そのため、**無限ループ**でずっと耳を傾け続けます。

```go
for {
    // メッセージが来るまでここで待機する
    _, p, err := conn.ReadMessage()
    if err != nil {
        break // 切断されたらループを抜ける
    }
    
    // 届いたデータを処理する（DBに保存して全員に配る）
    var input struct { /* ... */ }
    json.Unmarshal(p, &input)
    h.usecase.SaveAndPublishMessage(input.Sender, input.Content)
}
```

### ④ 全員に配る（ブロードキャスト）
誰かが発言したメッセージを、名簿に載っている全員に送信します。

```go
func (h *WebSocketHandler) broadcast(payload []byte) {
    h.clientsMu.Lock()
    defer h.clientsMu.Unlock()

    // 名簿（clients）に載っている全員に対してループ処理
    for client := range h.clients {
        // メッセージを送信
        err := client.WriteMessage(websocket.TextMessage, payload)
        if err != nil {
            // 送信失敗した場合は、名簿から削除して切断
            client.Close()
            delete(h.clients, client)
        }
    }
}
```

---

## 4. プロ級の工夫：Redis Pub/Sub なぜ必要？

今回のアプリはただのWebSocketではなく、裏側に `Redis` が挟まっています。これはなぜでしょうか？

もしGoのサーバーが1台だけなら、上記の「名簿（clients）」全員に配るだけで済みます。しかし、ユーザーが1万人になり、**Goのサーバーを3台に増やした（スケールアウトした）**とします。

* Aさんはサーバー1に接続
* Bさんはサーバー2に接続

この時、Aさんが発言しても、サーバー1の「名簿」にBさんは載っていないので、Bさんにメッセージが届きません。
これを解決するのが **Redis Pub/Sub** です。

1. Aさんが発言する。
2. サーバー1は自分の名簿に配ると同時に、**Redis（全サーバー共通の掲示板）に「Aが発言したよ！」と書き込む（Publish）**。
3. サーバー2や3は、常にRedisを監視（Subscribe）している。
4. Redisから通知を受け取ったサーバー2や3は、自分の名簿に載っているクライアント（Bさんなど）に配る。

これにより、サーバーが何台に増えても、全員にリアルタイムでメッセージが届く**「スケーラブルな構成」**が実現されています。今回のデモアプリでは、学習のためにあえてこの本格的な構成を採用しています。

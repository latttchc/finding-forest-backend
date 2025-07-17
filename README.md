# Finding Forest Backend

就活中の学生向けの匿名掲示板アプリケーションのバックエンドAPI。企業ごとの面接情報、ES対策、雰囲気共有などを行うためのREST APIを提供します。

## 🛠 技術スタック

- **Go 1.22+**
- **Echo Framework** - HTTP Webフレームワーク
- **GORM** - ORM（Object-Relational Mapping）
- **PostgreSQL** - データベース
- **go-playground/validator** - バリデーション
- **Docker** - コンテナ化

## 📁 プロジェクト構成

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # アプリケーションエントリーポイント
├── internal/
│   ├── config/
│   │   └── config.go            # 設定管理
│   ├── handlers/
│   │   ├── post.go              # 投稿ハンドラー
│   │   └── comment.go           # コメントハンドラー
│   ├── models/
│   │   ├── post.go              # 投稿モデル
│   │   └── comment.go           # コメントモデル
│   ├── repositories/
│   │   ├── post.go              # 投稿データアクセス層
│   │   └── comment.go           # コメントデータアクセス層
│   ├── services/
│   │   ├── post.go              # 投稿ビジネスロジック
│   │   └── comment.go           # コメントビジネスロジック
│   └── validators/
│       └── validator.go         # カスタムバリデーター
├── pkg/
│   └── database/
│       └── database.go          # データベース接続
├── go.mod
├── go.sum
├── Dockerfile
├── .env.example
└── README.md
```

## 🗄️ データベース設計

### Post（投稿）テーブル
```sql
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(20) NOT NULL,
    company_name VARCHAR(50) NOT NULL,
    job_type VARCHAR(30),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### Comment（コメント）テーブル
```sql
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL REFERENCES posts(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## 🔧 API エンドポイント

### ヘルスチェック
- `GET /health` - サーバーの稼働状況確認

### 投稿関連
- `GET /api/posts` - 投稿一覧取得（検索・フィルタ・ページネーション対応）
- `GET /api/posts/:id` - 投稿詳細取得
- `POST /api/posts` - 新規投稿作成

### コメント関連
- `POST /api/comments` - コメント作成
- `GET /api/posts/:post_id/comments` - 特定投稿のコメント一覧取得

## 🚀 セットアップ

### 1. 環境変数の設定

`.env.example`を参考に`.env`ファイルを作成してください：

```bash
cp .env.example .env
```

### 2. PostgreSQLの起動

Dockerを使用してPostgreSQLを起動：

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_DB=jobboard \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  postgres:15
```

### 3. 依存関係のインストール

```bash
go mod download
```

### 4. アプリケーションの起動

```bash
go run cmd/server/main.go
```

サーバーは`http://localhost:8080`で起動します。

## 🐳 Docker での実行

### ビルド

```bash
docker build -t finding-forest-backend .
```

### 実行

```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=password \
  -e DB_NAME=jobboard \
  -e DB_SSLMODE=disable \
  finding-forest-backend
```

## 📝 API使用例

### 投稿作成

```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Google面接体験談",
    "content": "Googleの面接を受けた際の体験談です...",
    "category": "面接",
    "company_name": "Google",
    "job_type": "エンジニア"
  }'
```

### 投稿一覧取得

```bash
curl "http://localhost:8080/api/posts?page=1&limit=10&category=面接&company_name=Google"
```

### コメント作成

```bash
curl -X POST http://localhost:8080/api/comments \
  -H "Content-Type: application/json" \
  -d '{
    "post_id": 1,
    "content": "参考になります！"
  }'
```

## 🔒 バリデーション

### 投稿
- `title`: 必須、1-100文字
- `content`: 必須、1-2000文字
- `category`: 必須、「面接」「ES」「企業情報」「その他」のいずれか
- `company_name`: 必須、1-50文字
- `job_type`: 任意、最大30文字

### コメント
- `post_id`: 必須、存在する投稿ID
- `content`: 必須、1-300文字

## 🧪 テスト

```bash
go test ./...
```

## 📊 ヘルスチェック

```bash
curl http://localhost:8080/health
```

レスポンス：
```json
{
  "status": "ok"
}
```

## 🔄 マイグレーション

アプリケーション起動時にGORMのAutoMigrate機能により自動的にテーブルが作成されます。

## 🌐 CORS設定

すべてのオリジンからのアクセスを許可しています。本番環境では適切に制限してください。

## 📈 今後の拡張予定

- 認証機能の追加
- 投稿の編集・削除機能
- いいね機能
- 通報機能
- 管理者機能
- キャッシュ機能（Redis）
- 全文検索機能
- API レート制限

## 🤝 コントリビューション

1. このリポジトリをフォーク
2. 新しいブランチを作成 (`git checkout -b feature/new-feature`)
3. 変更をコミット (`git commit -am 'Add new feature'`)
4. ブランチにプッシュ (`git push origin feature/new-feature`)
5. Pull Requestを作成

## 📄 ライセンス

MIT License

## 📞 サポート

質問や問題がある場合は、GitHubのIssueを作成してください。

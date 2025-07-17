package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/latttchc/finding-forest-backend/internal/config"
	"github.com/latttchc/finding-forest-backend/internal/handlers"
	"github.com/latttchc/finding-forest-backend/internal/repositories"
	"github.com/latttchc/finding-forest-backend/internal/services"
	"github.com/latttchc/finding-forest-backend/internal/validators"
	"github.com/latttchc/finding-forest-backend/pkg/database"
)

func main() {
	// 設定読み込み
	cfg := config.Load()

	// データベース接続
	db, err := database.Connect(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// データベースマイグレーション
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// バリデーター初期化
	validate := validator.New()

	// リポジトリ初期化
	postRepo := repositories.NewPostRepository(db)
	commentRepo := repositories.NewCommentRepository(db)

	// サービス初期化
	postService := services.NewPostService(postRepo, commentRepo, validate)
	commentService := services.NewCommentService(commentRepo, postRepo, validate)

	// ハンドラー初期化
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)

	// Echo インスタンス作成
	e := echo.New()

	// カスタムバリデーター設定
	e.Validator = validators.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ヘルスチェック
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// API ルート設定
	api := e.Group("/api")

	// 投稿関連のルート
	api.GET("/posts", postHandler.GetPosts)
	api.GET("/posts/:id", postHandler.GetPost)
	api.POST("/posts", postHandler.CreatePost)

	// コメント関連のルート
	api.POST("/comments", commentHandler.CreateComment)
	api.GET("/posts/:post_id/comments", commentHandler.GetCommentsByPostID)

	// サーバー起動
	log.Printf("Server starting on port %s", cfg.Server.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))
}

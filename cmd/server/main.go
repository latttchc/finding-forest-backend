package server

import (
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/latttchc/finding-forest-backend/internal/config"
	"github.com/latttchc/finding-forest-backend/pkg/database"
	"google.golang.org/grpc/profiling/service"
	"honnef.co/go/tools/config"
)

func main() {
	// 設定読み込み
	cfg := config.Load()

	// データベース接続
	db, err := database.Connect(cfg.GetDSN())
	if err != nil {
		log.Fatal("ロードに失敗しました。: %v", err)
	}

	// データベースマイグレーション
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database: %v", err)
	}

	// バリデータ初期化
	validator := validator.New()

	// リポジトリ初期化
	postRepo := repositories.NewPostRepository(db)
	commentRepo := repositories.NewCommentRepository(db)

	// サービス初期化
	postService := service.NewPostService(postRepo, validator)
	commentService := service.NewCommentService(commentRepo, validator)

	// ハンドラー初期化
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := hanlers.NewCommentHandler(commentService)

	// Echo作成
	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// API Route
	api := e.Group("/api")

	// 投稿関連のルート
	api.GET("/posts", postHandler.GetPosts)
	api.GET("/posts/:id", postHandler.GetPost)
	api.POST("/posts", postHandler.CreatePost)

	// コメント関連のルート
	api.POST("/commments", commentHandler.CreateComment)

	// サーバー起動
	log.Printf("Server starting on port %s", &cfg.Server.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))
}

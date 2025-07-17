package services

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/latttchc/finding-forest-backend/internal/models"
	"github.com/latttchc/finding-forest-backend/internal/repositories"
)

// PostService は投稿に関するビジネスロジックを定義するインターフェースです
type PostService interface {
	CreatePost(req *models.PostCreateRequest) (*models.PostResponse, error)
	GetPost(id uint) (*models.PostDetailResponse, error)
	GetPosts(page, limit int, category, companyName string) (*PostListResult, error)
}

// PostListResult は投稿一覧取得の結果を表す構造体です
type PostListResult struct {
	Posts      []models.PostListResponse `json:"posts"`       // 投稿一覧
	Total      int64                     `json:"total"`       // 総件数
	Page       int                       `json:"page"`        // 現在のページ
	Limit      int                       `json:"limit"`       // 1ページあたりの件数
	TotalPages int                       `json:"total_pages"` // 総ページ数
}

// postService は PostService インターフェースの実装です
type postService struct {
	postRepo    repositories.PostRepository    // 投稿データアクセス層
	commentRepo repositories.CommentRepository // コメントデータアクセス層
	validator   *validator.Validate            // バリデーター
}

// NewPostService は新しい PostService インスタンスを作成します
func NewPostService(postRepo repositories.PostRepository, commentRepo repositories.CommentRepository, validator *validator.Validate) PostService {
	return &postService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		validator:   validator,
	}
}

// CreatePost は新しい投稿を作成します
// バリデーションを実行後、データベースに保存し、レスポンスを返します
func (s *postService) CreatePost(req *models.PostCreateRequest) (*models.PostResponse, error) {
	// バリデーション
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// リクエストをモデルに変換
	post := &models.Post{
		Title:       req.Title,
		Content:     req.Content,
		Category:    req.Category,
		CompanyName: req.CompanyName,
		JobType:     req.JobType,
	}

	// データベースに保存
	if err := s.postRepo.Create(post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// レスポンスに変換
	response := &models.PostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		Category:    post.Category,
		CompanyName: post.CompanyName,
		JobType:     post.JobType,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}

	return response, nil
}

// GetPost は指定されたIDの投稿詳細を取得します
// 投稿に関連するコメントも含めて取得します
func (s *postService) GetPost(id uint) (*models.PostDetailResponse, error) {
	// 投稿とコメントを取得
	post, err := s.postRepo.GetWithComments(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// コメントをレスポンス形式に変換
	comments := make([]models.CommentResponse, len(post.Comments))
	for i, comment := range post.Comments {
		comments[i] = models.CommentResponse{
			ID:        comment.ID,
			PostID:    comment.PostID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
	}

	// レスポンスに変換
	response := &models.PostDetailResponse{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		Category:    post.Category,
		CompanyName: post.CompanyName,
		JobType:     post.JobType,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Comments:    comments,
	}

	return response, nil
}

// GetPosts は投稿一覧を取得します
// ページネーション、カテゴリフィルタ、企業名検索に対応しています
func (s *postService) GetPosts(page, limit int, category, companyName string) (*PostListResult, error) {
	// ページネーション設定のバリデーション
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// 投稿一覧を取得
	posts, total, err := s.postRepo.GetAll(limit, offset, category, companyName)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	// レスポンス形式に変換
	postResponses := make([]models.PostListResponse, len(posts))
	for i, post := range posts {
		// 各投稿のコメント数を取得
		commentCount, err := s.commentRepo.CountByPostID(post.ID)
		if err != nil {
			commentCount = 0 // エラーの場合は0とする
		}

		postResponses[i] = models.PostListResponse{
			ID:           post.ID,
			Title:        post.Title,
			Category:     post.Category,
			CompanyName:  post.CompanyName,
			JobType:      post.JobType,
			CreatedAt:    post.CreatedAt,
			CommentCount: commentCount,
		}
	}

	// 総ページ数を計算
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &PostListResult{
		Posts:      postResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

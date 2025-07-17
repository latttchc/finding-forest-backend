package services

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/latttchc/finding-forest-backend/internal/models"
	"github.com/latttchc/finding-forest-backend/internal/repositories"
)

// CommentService はコメントに関するビジネスロジックを定義するインターフェースです
type CommentService interface {
	CreateComment(req *models.CommentCreateRequest) (*models.CommentResponse, error)
	GetCommentsByPostID(postID uint) ([]models.CommentResponse, error)
}

// commentService は CommentService インターフェースの実装です
type commentService struct {
	commentRepo repositories.CommentRepository // コメントデータアクセス層
	postRepo    repositories.PostRepository    // 投稿データアクセス層
	validator   *validator.Validate            // バリデーター
}

// NewCommentService は新しい CommentService インスタンスを作成します
func NewCommentService(commentRepo repositories.CommentRepository, postRepo repositories.PostRepository, validator *validator.Validate) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		validator:   validator,
	}
}

// CreateComment は新しいコメントを作成します
// バリデーションと投稿の存在確認を行った後、データベースに保存します
func (s *commentService) CreateComment(req *models.CommentCreateRequest) (*models.CommentResponse, error) {
	// バリデーション
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 投稿が存在するかチェック
	_, err := s.postRepo.GetByID(req.PostID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	// リクエストをモデルに変換
	comment := &models.Comment{
		PostID:  req.PostID,
		Content: req.Content,
	}

	// データベースに保存
	if err := s.commentRepo.Create(comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// レスポンスに変換
	response := &models.CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	return response, nil
}

// GetCommentsByPostID は指定された投稿のコメント一覧を取得します
// 投稿の存在確認を行った後、コメントを取得します
func (s *commentService) GetCommentsByPostID(postID uint) ([]models.CommentResponse, error) {
	// 投稿が存在するかチェック
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}

	// コメントを取得
	comments, err := s.commentRepo.GetByPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	// レスポンス形式に変換
	responses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = models.CommentResponse{
			ID:        comment.ID,
			PostID:    comment.PostID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
	}

	return responses, nil
}

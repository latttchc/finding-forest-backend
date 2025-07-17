package services

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/latttchc/finding-forest-backend/internal/models"
	"github.com/latttchc/finding-forest-backend/internal/repositories"
)

type PostService interface {
	CreatePost(req *models.PostCreateRequest) (*models.PostResponse, error)
	GetPost(id uint) (*models.PostDetailResponse, error)
	GetPosts(page, limit int, category, companyName string) (*PostListResult, error)
}

type PostListResult struct {
	Posts      []models.PostListResponse `json:"posts"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

type postService struct {
	postRepo    repositories.PostRepository
	commentRepo repositories.CommentRepository
	validator   *validator.Validate
}

func NewPostService(postRepo repositories.PostRepository, commentRepo repositories.CommentRepository, validator *validator.Validate) PostService {
	return &postService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		validator:   validator,
	}
}

func (s *postService) CreatePost(req *models.PostCreateRequest) (*models.PostResponse, error) {
	// バリデーション
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// モデルに変換
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

func (s *postService) GetPosts(page, limit int, category, companyName string) (*PostListResult, error) {
	// ページネーション設定
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
		// コメント数を取得
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

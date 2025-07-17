package repositories

import (
	"github.com/latttchc/finding-forest-backend/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id uint) (*models.Post, error)
	GetAll(limit, offset int, category, companyName string) ([]models.Post, int64, error)
	GetWithComments(id uint) (*models.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetAll(limit, offset int, category, companyName string) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := r.db.Model(&models.Post{})

	// フィルタリング
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if companyName != "" {
		query = query.Where("company_name ILIKE ?", "%"+companyName+"%")
	}

	// 総数を取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// データを取得
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) GetWithComments(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("Comments").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

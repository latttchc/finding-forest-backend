package repositories

import (
	"github.com/latttchc/finding-forest-backend/internal/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	GetByPostID(postID uint) ([]models.Comment, error)
	GetByID(id uint) (*models.Comment, error)
	CountByPostID(postID uint) (int64, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) GetByPostID(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("post_id = ?", postID).
		Order("create_at DESC").
		Find(&comments).Error
	return comments, err
}

func (r *commentRepository) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) CountByPostID(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

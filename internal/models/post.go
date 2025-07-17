package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null" validate:"required,min=1,max=100"`
	Content     string         `json:"content" gorm:"type:text;not null" validate:"required,min=1,max=2000"`
	Category    string         `json:"category" gorm:"not null" validate:"required,oneof=面接 ES 企業情報 その他"`
	CompanyName string         `json:"company_name" gorm:"not null" validate:"required,min=1,max=50"`
	JobType     string         `json:"job_type" validate:"max=30"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
}

// PostCreateRequest は投稿作成リクエストの構造体
type PostCreateRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=100"`
	Content     string `json:"content" validate:"required,min=1,max=2000"`
	Category    string `json:"category" validate:"required,oneof=面接 ES 企業情報 その他"`
	CompanyName string `json:"company_name" validate:"required,min=1,max=50"`
	JobType     string `json:"job_type" validate:"max=30"`
}

// PostResponse は投稿レスポンスの構造体
type PostResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Category    string    `json:"category"`
	CompanyName string    `json:"company_name"`
	JobType     string    `json:"job_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PostListResponse は投稿一覧レスポンスの構造体
type PostListResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Category     string    `json:"category"`
	CompanyName  string    `json:"company_name"`
	JobType      string    `json:"job_type"`
	CreatedAt    time.Time `json:"created_at"`
	CommentCount int64     `json:"comment_count"`
}

// PostDetailResponse は投稿詳細レスポンスの構造体
type PostDetailResponse struct {
	ID          uint              `json:"id"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Category    string            `json:"category"`
	CompanyName string            `json:"company_name"`
	JobType     string            `json:"job_type"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Comments    []CommentResponse `json:"comments"`
}

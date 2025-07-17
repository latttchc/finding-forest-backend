package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	PostID    uint           `json:"post_id" gorm:"not null"`
	Content   string         `json:"content" gorm:"type:text;not null" validate:"required,min=1,max=300"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Post Post `json:"-" gorm:"foreignKey:PostID"`
}

// CommentCreateRequest はコメント作成リクエストの構造体
type CommentCreateRequest struct {
	PostID  uint   `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required,min=1,max=300"`
}

// CommentResponse はコメントレスポンスの構造体
type CommentResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

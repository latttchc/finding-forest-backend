package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/latttchc/finding-forest-backend/internal/models"
	"github.com/latttchc/finding-forest-backend/internal/services"
)

type CommentHandler struct {
	commentService services.CommentService
}

func NewCommentHandler(commentService services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment は新しいコメントを作成します
func (h *CommentHandler) CreateComment(c echo.Context) error {
	var req models.CommentCreateRequest

	// リクエストボディをバインド
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// サービス層を呼び出し
	response, err := h.commentService.CreateComment(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response)
}

// GetCommentsByPostID は指定された投稿のコメント一覧を取得します
func (h *CommentHandler) GetCommentsByPostID(c echo.Context) error {
	// パスパラメータからpost_idを取得
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid post ID",
		})
	}

	// サービス層を呼び出し
	response, err := h.commentService.GetCommentsByPostID(uint(postID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"comments": response,
	})
}

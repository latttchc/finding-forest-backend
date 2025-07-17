package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/latttchc/finding-forest-backend/internal/models"
	"github.com/latttchc/finding-forest-backend/internal/services"
)

// PostHandler は投稿に関するHTTPリクエストを処理するハンドラーです
type PostHandler struct {
	postService services.PostService
}

// NewPostHandler は新しい PostHandler インスタンスを作成します
func NewPostHandler(postService services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost は新しい投稿を作成するHTTPハンドラーです
// POST /api/posts
func (h *PostHandler) CreatePost(c echo.Context) error {
	var req models.PostCreateRequest

	// リクエストボディをバインド
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// サービス層を呼び出し
	response, err := h.postService.CreatePost(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response)
}

// GetPost は指定されたIDの投稿詳細を取得するHTTPハンドラーです
// GET /api/posts/:id
func (h *PostHandler) GetPost(c echo.Context) error {
	// パスパラメータからIDを取得
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid post ID",
		})
	}

	// サービス層を呼び出し
	response, err := h.postService.GetPost(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// GetPosts は投稿一覧を取得するHTTPハンドラーです
// GET /api/posts?page=1&limit=20&category=面接&company_name=Google
func (h *PostHandler) GetPosts(c echo.Context) error {
	// クエリパラメータを取得
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")
	category := c.QueryParam("category")
	companyName := c.QueryParam("company_name")

	// ページネーション設定
	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// サービス層を呼び出し
	response, err := h.postService.GetPosts(page, limit, category, companyName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

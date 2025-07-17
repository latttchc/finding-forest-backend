package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator は Echo フレームワーク用のカスタムバリデーターです
type CustomValidator struct {
	validator *validator.Validate
}

// Validate は構造体のバリデーションを実行します
// Echo の c.Validate() メソッドから呼び出されます
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// New は新しいカスタムバリデーターのインスタンスを作成します
// Echo のバリデーターインターフェースを実装しています
func New() echo.Validator {
	return &CustomValidator{validator: validator.New()}
}

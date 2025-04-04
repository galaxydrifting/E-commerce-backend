package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"
	"e-commerce/repository"
	"e-commerce/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// 創建必要的依賴
	mockUserRepo := repository.NewMockUserRepository()
	authService := services.NewAuthService(mockUserRepo)
	authController := controllers.NewAuthController(authService)
	authMiddleware := middlewares.NewAuthMiddleware(nil, authService)

	// 設置路由
	SetupAuthRoutes(r, authController, authMiddleware)

	// 定義要測試的路由
	routes := []struct {
		name         string
		method       string
		path         string
		expectedPath string
	}{
		{"Register", "POST", "/api/v1/auth/register", "/api/v1/auth/register"},
		{"Login", "POST", "/api/v1/auth/login", "/api/v1/auth/login"},
		{"Logout", "POST", "/api/v1/auth/logout", "/api/v1/auth/logout"},
		{"Get Profile", "GET", "/api/v1/auth/profile", "/api/v1/auth/profile"},
		{"Update Profile", "PUT", "/api/v1/auth/profile", "/api/v1/auth/profile"},
	}

	for _, route := range routes {
		t.Run(route.name, func(t *testing.T) {
			// 創建測試請求
			req := httptest.NewRequest(route.method, route.path, nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			// 驗證路徑是否正確註冊
			// 注意：未授權的請求會返回 401，這是正常的
			// 我們只是驗證路由是否正確註冊
			assert.NotEqual(t, http.StatusNotFound, resp.Code,
				"Route %s should exist but got 404", route.path)
		})
	}
}

//go:build wireinject
// +build wireinject

package di

import (
	"e-commerce/configs"
	"e-commerce/controllers"
	"e-commerce/middlewares"
	"e-commerce/repository"
	"e-commerce/services"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// Container 定義應用程式的依賴注入容器
type Container struct {
	DB             *gorm.DB
	AuthController *controllers.AuthController
	AuthMiddleware *middlewares.AuthMiddleware
}

// provideGormDB 提供原始的 gorm.DB 實例
func provideGormDB(envFile string) *gorm.DB {
	return configs.ConnectDB(envFile).DB
}

// provideConfigsDatabase 提供包裝後的 configs.Database 實例
func provideConfigsDatabase(db *gorm.DB) *configs.Database {
	return &configs.Database{DB: db}
}

// Initialize 初始化應用程式依賴
func Initialize(envFile string) (*Container, error) {
	wire.Build(
		// Database
		provideGormDB,
		provideConfigsDatabase,

		// Repository
		repository.NewGormUserRepository,
		wire.Bind(new(repository.UserRepository), new(*repository.GormUserRepository)),

		// Service
		services.NewAuthService,

		// Controller
		controllers.NewAuthController,

		// Middleware
		middlewares.NewAuthMiddleware,

		// Container
		wire.Struct(new(Container), "DB", "AuthController", "AuthMiddleware"),
	)
	return nil, nil
}

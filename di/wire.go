//go:build wireinject
// +build wireinject

package di

import (
	"e-commerce/configs"
	"e-commerce/controllers"
	"e-commerce/repository"
	"e-commerce/services"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// Container 定义应用程序的依赖注入容器
type Container struct {
	DB             *gorm.DB
	AuthController *controllers.AuthController
}

// provideDB 提供数据库实例
func provideDB(envFile string) *gorm.DB {
	return configs.ConnectDB(envFile).DB
}

// InitializeAPI 初始化所有依赖
var InitializeSet = wire.NewSet(
	provideDB,
	repository.NewGormUserRepository,
	services.NewAuthService,
	controllers.NewAuthController,
	wire.Struct(new(Container), "*"),
)

func Initialize(envFile string) (*Container, error) {
	wire.Build(InitializeSet)
	return nil, nil
}

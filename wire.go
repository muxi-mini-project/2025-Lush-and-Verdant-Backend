//go:build wireinject
// +build wireinject

package main

import (
	"2025-Lush-and-Verdant-Backend/client"
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/middleware"
	"2025-Lush-and-Verdant-Backend/route"
	"2025-Lush-and-Verdant-Backend/service"
	"2025-Lush-and-Verdant-Backend/tool"
	"github.com/google/wire"
)

func InitApp(ConfigPath string) (*route.App, error) {
	wire.Build(
		route.ProviderSet,
		controller.ProviderSet,
		service.ProviderSet,
		client.ProviderSet,
		middleware.ProviderSet,
		tool.ProviderSet,
		config.ProviderSet,
		dao.ProviderSet,
		wire.Bind(new(service.GoalService), new(*service.GoalServiceImpl)),
		wire.Bind(new(service.UserService), new(*service.UserServiceImpl)),
		wire.Bind(new(service.SloganService), new(*service.SloganServiceImpl)),
	)
	return &route.App{}, nil
}

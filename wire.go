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
		wire.Bind(new(service.ImageService), new(*service.ImageServiceImpl)),
		wire.Bind(new(service.ChatService), new(*service.ChatServiceImpl)),
		wire.Bind(new(service.GroupService), new(*service.GroupServiceImpl)),
		wire.Bind(new(service.LikeService), new(*service.LikeServiceImpl)),
		wire.Bind(new(dao.UserDAO), new(*dao.UserDAOImpl)),
		wire.Bind(new(dao.GoalDAO), new(*dao.GoalDAOImpl)),
		wire.Bind(new(dao.SloganDAO), new(*dao.SloganDAOImpl)),
		wire.Bind(new(dao.ImageDAO), new(*dao.ImageDAOImpl)),
		wire.Bind(new(dao.EmailCodeDAO), new(*dao.EmailCodeDAOImpl)),
		wire.Bind(new(dao.GroupDAO), new(*dao.GroupDAOImpl)),
		wire.Bind(new(dao.ChatDAO), new(*dao.ChatDAOImpl)),
		wire.Bind(new(dao.LikeDAO), new(*dao.LikeDAOImpl)),
	)
	return &route.App{}, nil
}

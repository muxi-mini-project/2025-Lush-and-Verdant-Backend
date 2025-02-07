package main

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/route"
	"2025-Lush-and-Verdant-Backend/service"
	"github.com/google/wire"
	"log"
)

func InitApp() (*route.App, error) {
	wire.Build(
		dao.DAOSet,
		service.ServiceSet,
		controller.ControllerSet,
		route.RouteSet,
	)
	return nil, nil
}

func main() {
	app, err := InitApp()
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}

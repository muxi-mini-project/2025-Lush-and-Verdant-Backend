package main

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/route"
)

func main() {

	uc := controller.NewUserController()
	us := route.NewUserSvc(uc)
	app := route.NewApp(us)
	app.Run()
}

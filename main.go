package main

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/route"
)

func main() {
	//controller.CreateSlogan()

	uc := controller.NewUserController()
	us := route.NewUserSvc(uc)

	sc := controller.NewSloganController()
	ss := route.NewSloganSvc(sc)

	gc := controller.NewGoalController()
	gs := route.NewGoalSvc(gc)

	app := route.NewApp(us, ss, gs)
	app.Run()

}

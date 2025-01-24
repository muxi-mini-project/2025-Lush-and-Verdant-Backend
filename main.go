package main

import (
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/route"
	"2025-Lush-and-Verdant-Backend/service"
)

func main() {
	//controller.CreateSlogan()
	db := dao.NewDB(config.Dsn)
	usr := service.NewUserService(db)
	uc := controller.NewUserController(usr)
	us := route.NewUserSvc(uc)

	ssr := service.NewSloganService(db)
	sc := controller.NewSloganController(ssr)
	ss := route.NewSloganSvc(sc)

	gsr := service.NewGoalService(db)
	gc := controller.NewGoalController(gsr)
	gs := route.NewGoalSvc(gc)

	app := route.NewApp(us, ss, gs)
	app.Run()

}

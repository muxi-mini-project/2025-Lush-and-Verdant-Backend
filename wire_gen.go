// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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
)

// Injectors from wire.go:

func InitApp(ConfigPath string) (*route.App, error) {
	viperSetting := config.NewViperSetting(ConfigPath)
	mySQLConfig := config.NewMySQLConfig(viperSetting)
	db, err := dao.MySQLDB(mySQLConfig)
	if err != nil {
		return nil, err
	}
	userDAOImpl := dao.NewUserDAO(db)
	redisConfig := config.NewRedisConfig(viperSetting)
	redisClient, err := dao.RedisDB(redisConfig)
	if err != nil {
		return nil, err
	}
	emailCodeDAOImpl := dao.NewEmailCodeDAOImpl(redisClient)
	jwtConfig := config.NewJwtConfig(viperSetting)
	jwtClient := middleware.NewJwtClient(jwtConfig)
	qqConfig := config.NewQQConfig(viperSetting)
	mail := tool.NewMail(redisClient, qqConfig)
	priConfig := config.NewPriConfig(viperSetting)
	userServiceImpl := service.NewUserServiceImpl(userDAOImpl, emailCodeDAOImpl, jwtClient, mail, priConfig)
	userController := controller.NewUserController(userServiceImpl)
	userSvc := route.NewUserSvc(userController)
	sloganDAOImpl := dao.NewSloganDAOImpl(db)
	sloganServiceImpl := service.NewSloganServiceImpl(sloganDAOImpl, userDAOImpl)
	sloganController := controller.NewSloganController(sloganServiceImpl)
	sloganSvc := route.NewSloganSvc(sloganController, jwtClient)
	goalDAOImpl := dao.NewGoalDAOImpl(db)
	goalServiceImpl := service.NewGoalServiceImpl(goalDAOImpl)
	chatGptConfig := config.NewChatGptConfig(viperSetting)
	chatGptClient := client.NewChatGptClient(chatGptConfig)
	goalController := controller.NewGoalController(goalServiceImpl, chatGptClient)
	goalSvc := route.NewGoalSvc(goalController, jwtClient)
	qiNiuYunConfig := config.NewQNYConfig(viperSetting)
	imageDAOImpl := dao.NewImageDAO(db)
	imageServiceImpl := service.NewImageServiceImpl(qiNiuYunConfig, imageDAOImpl)
	imageController := controller.NewImageController(imageServiceImpl)
	imageSvc := route.NewImageSvc(imageController, jwtClient)
	groupDAOImpl := dao.NewGroupDAOImpl(db, redisClient)
	groupServiceImpl := service.NewGroupServiceImpl(groupDAOImpl)
	groupController := controller.NewGroupController(groupServiceImpl)
	groupSve := route.NewGroupSve(groupController, jwtClient)
	chatDAOImpl := dao.NewChatDAOImpl(redisClient)
	chatServiceImpl := service.NewChatServiceImpl(chatDAOImpl, groupDAOImpl)
	chatController := controller.NewChatController(chatServiceImpl, jwtClient)
	chatSve := route.NewChatSve(chatController, jwtClient)
	kafkaConfig := config.NewKafkaConfig(viperSetting)
	likeDAOImpl := dao.NewLikeDAOImpl(redisClient)
	likeServiceImpl := service.NewLikeServiceImpl(kafkaConfig, likeDAOImpl)
	likeController := controller.NewLikeController(likeServiceImpl)
	likeSvc := route.NewLikeSvc(likeController, jwtClient)
	app := route.NewApp(userSvc, sloganSvc, goalSvc, imageSvc, groupSve, chatSve, likeSvc)
	return app, nil
}

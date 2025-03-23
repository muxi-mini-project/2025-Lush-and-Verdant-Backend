package service

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewGoalServiceImpl,
	NewSloganServiceImpl,
	NewUserServiceImpl,
	NewImageServiceImpl,
	NewChatServiceImpl,
	NewGroupServiceImpl,
	NewLikeServiceImpl,
)

package controller

import "github.com/google/wire"

// ProviderSet 就是简单的wire的provider
var ProviderSet = wire.NewSet(
	NewUserController,
	NewSloganController,
	NewGoalController,
	NewImageController,
	NewGroupController,
	NewChatController,
	NewLikeController,
)

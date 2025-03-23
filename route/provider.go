package route

import "github.com/google/wire"

// ProviderSet 也是wire的provider
var ProviderSet = wire.NewSet(
	NewSloganSvc,
	NewUserSvc,
	NewGoalSvc,
	NewImageSvc,
	NewChatSve,
	NewGroupSve,
	NewApp,
	NewLikeSvc,
)

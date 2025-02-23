package dao

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewDB,
	NewUserDAO,
	NewGoalDAOImpl,
	NewSloganDAOImpl,
	NewEmailDAOImpl,
	NewImageDAO,
)

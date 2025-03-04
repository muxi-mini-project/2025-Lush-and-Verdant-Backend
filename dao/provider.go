package dao

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	MySQLDB,
	RedisDB,
	NewUserDAO,
	NewGoalDAOImpl,
	NewSloganDAOImpl,
	NewEmailDAOImpl,
	NewEmailCodeDAOImpl,
	NewImageDAO,
)

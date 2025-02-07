package dao

import (
	"2025-Lush-and-Verdant-Backend/config"
	"github.com/google/wire"
)

var DAOSet = wire.NewSet(NewDB, wire.Value(config.Dsn))

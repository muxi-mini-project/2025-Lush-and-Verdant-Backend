package config

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewViperSetting,
	NewDatabaseConfig,
	NewJwtConfig,
	NewPriConfig,
	NewTimeLayoutConfig,
	NewChatGptConfig,
	NewQQConfig,
)

type DatabaseConfig struct {
	Addr string `yaml:"addr"`
}
type JwtConfig struct {
	SecretKey string `yaml:"secretkey"`
}
type PriConfig struct {
	Name string `yaml:"name"`
}
type TimeLayoutConfig struct {
	Template string `yaml:"template"`
}
type ChatGptConfig struct {
	Sdk string `yaml:"sdk"`
}
type QQConfig struct {
	Email string `yaml:"email"`
	Key   string `yaml:"key"`
}

func NewDatabaseConfig(vs *ViperSetting) *DatabaseConfig {
	var databaseConfig = &DatabaseConfig{}
	vs.ReadSection("database", &databaseConfig)
	return databaseConfig
}

func NewJwtConfig(vs *ViperSetting) *JwtConfig {
	var jwtConfig = &JwtConfig{}
	vs.ReadSection("jwt", &jwtConfig)
	return jwtConfig
}
func NewPriConfig(vs *ViperSetting) *PriConfig {
	var priConfig = &PriConfig{}
	vs.ReadSection("pri", &priConfig)
	return priConfig
}
func NewTimeLayoutConfig(vs *ViperSetting) *TimeLayoutConfig {
	var timeLayoutConfig = &TimeLayoutConfig{}
	vs.ReadSection("timelayout", &timeLayoutConfig)
	return timeLayoutConfig
}
func NewChatGptConfig(vs *ViperSetting) *ChatGptConfig {
	var chatGptConfig = &ChatGptConfig{}
	vs.ReadSection("chatgpt", &chatGptConfig)
	return chatGptConfig
}
func NewQQConfig(vs *ViperSetting) *QQConfig {
	var qqConfig = &QQConfig{}
	vs.ReadSection("qq", &qqConfig)
	return qqConfig
}

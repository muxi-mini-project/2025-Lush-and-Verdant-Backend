package config

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewViperSetting,
	NewMySQLConfig,
	NewRedisConfig,
	NewJwtConfig,
	NewPriConfig,
	NewTimeLayoutConfig,
	NewChatGptConfig,
	NewQQConfig,
	NewQNYConfig,
	NewKafkaConfig,
)

type MySQLConfig struct {
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

type QiNiuYunConfig struct {
	AccessKey  string `yaml:"accessKey"`  // 七牛云提供的公钥
	SecretKey  string `yaml:"secretKey"`  // 七牛云提供的密钥
	BucketName string `yaml:"bucketName"` // 所创建七牛云对象存储的名称
	DomainName string `yaml:"domainName"` // 对象存储所绑定的七牛云的域名
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type KafkaConfig struct {
	Addr string `yaml:"addr"`
}

func NewMySQLConfig(vs *ViperSetting) *MySQLConfig {
	var databaseConfig = &MySQLConfig{}
	vs.ReadSection("database", &databaseConfig)
	return databaseConfig
}

func NewRedisConfig(vs *ViperSetting) *RedisConfig {
	var redisConfig = &RedisConfig{}
	vs.ReadSection("redis", &redisConfig)
	return redisConfig
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

func NewQNYConfig(vs *ViperSetting) *QiNiuYunConfig {
	var qiNiuYunConfig = &QiNiuYunConfig{}
	vs.ReadSection("qiniuyun", &qiNiuYunConfig)
	return qiNiuYunConfig
}

func NewKafkaConfig(vs *ViperSetting) *KafkaConfig {
	var kafkaConfig = &KafkaConfig{}
	vs.ReadSection("kafka", &kafkaConfig)
	return kafkaConfig
}

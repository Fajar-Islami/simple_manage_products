package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type RedisConf struct {
	Password          string `mapstructure:"redis_password"`
	Port              int    `mapstructure:"redis_port"`
	Host              string `mapstructure:"redis_host"`
	DB                int    `mapstructure:"redis_db"`
	DefaultDB         int    `mapstructure:"redis_Defaultdb"`
	RedisMinIdleConns int    `mapstructure:"redis_MinIdleConns"`
	RedisPoolSize     int    `mapstructure:"redis_PoolSize"`
	RedisPoolTimeout  int    `mapstructure:"redis_PoolTimeout"`
}

const currentfilepath = "internal/infrastructure/redis/redis.go"

func NewRedisClient(v *viper.Viper) *redis.Client {
	var redisConfig RedisConf
	err := v.Unmarshal(&redisConfig)
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelPanic, fmt.Sprintf("failed init database redis : %s", err.Error()), nil)
	}

	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		MinIdleConns: redisConfig.RedisMinIdleConns,
		PoolSize:     redisConfig.RedisPoolSize,
		PoolTimeout:  time.Duration(redisConfig.RedisPoolTimeout) * time.Second,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelPanic, "", fmt.Errorf("Cannot conenct to redis : %s", err.Error()))
		panic(err)
	}
	helper.Logger(currentfilepath, helper.LoggerLevelInfo, fmt.Sprintf("Redis ping : %s", pong), nil)

	return client
}

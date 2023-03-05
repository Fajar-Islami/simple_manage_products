package container

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Fajar-Islami/simple_manage_products/internal/helper"
	"github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/mysql"
	redisclient "github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/redis"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var v *viper.Viper

const currentfilepath = "internal/infrastructure/container/container.go"

type (
	Container struct {
		Mysqldb *gorm.DB
		Apps    *Apps
		Logger  *Logger
		Redis   *redis.Client
	}

	Logger struct {
		Log zerolog.Logger
	}

	Apps struct {
		Name      string `mapstructure:"name"`
		Host      string `mapstructure:"host"`
		Version   string `mapstructure:"version"`
		Address   string `mapstructure:"address"`
		HttpPort  int    `mapstructure:"httpport"`
		SecretJwt string `mapstructure:"secretJwt"`
	}
)

func loadEnv() {
	projectDirName := "simple_manage_products"
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	v.SetConfigFile(string(rootPath) + `/.env`)
}

func init() {
	v = viper.New()

	v.AutomaticEnv()
	loadEnv()

	path, err := os.Executable()
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelPanic, "", fmt.Errorf("os.Executable panic : %s", err.Error()))
	}

	dir := filepath.Dir(path)
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelPanic, "", fmt.Errorf("failed read config : %s", err.Error()))
	}

	err = v.ReadInConfig()
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelPanic, "", fmt.Errorf("failed init config : %s", err.Error()))
	}

	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "Succeed read configuration file", nil)
}

func AppsInit(v *viper.Viper) (apps Apps) {
	err := v.Unmarshal(&apps)
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelError, "", fmt.Errorf("error when unmarshal configuration file : ", err.Error()))
	}
	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "Succeed when unmarshal configuration file", nil)
	return
}

func LoggerInit() *Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return &Logger{
		Log: zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
	}
}

func InitContainer() (cont *Container) {
	apps := AppsInit(v)
	mysqldb := mysql.DatabaseInit(v)
	logger := LoggerInit()
	redisClient := redisclient.NewRedisClient(v)

	return &Container{
		Apps:    &apps,
		Mysqldb: mysqldb,
		Logger:  logger,
		Redis:   redisClient,
	}

}

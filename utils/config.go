package utils

import (
	"time"

	"github.com/spf13/viper"
)

// COnfig stores all configuration of the application
// The value are read by viper from a config file or environment variables
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	Enviroment string `mapstructure:"ENVIROMENT"`
}

func LoadConfig(path string) (config Config, err error) {
	// location of config file
	viper.AddConfigPath(path)
	// File name
	viper.SetConfigName("app")
	// File type
	viper.SetConfigType("env")

	// Automatic overide value from config CLI
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return	
	}

	err = viper.Unmarshal(&config)
	return
}

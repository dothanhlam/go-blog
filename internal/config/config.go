package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config stores all configuration for the application.
type Config struct {
	ServerPort   string `mapstructure:"SERVER_PORT"`
	DatabaseURL  string `mapstructure:"DATABASE_URL"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	StorageType  string `mapstructure:"STORAGE_TYPE"` // "local" or "s3"
	S3Bucket     string `mapstructure:"S3_BUCKET"`
	S3Region     string `mapstructure:"S3_REGION"`
	AWSAccessKey string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	TokenExpiresInHours int `mapstructure:"TOKEN_EXPIRES_IN_HOURS"`
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	// Tell viper to look for a .env file in the current directory
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Find and read the config file.
	// We can ignore the error if the file is not found,
	// as configuration can still be provided by system environment variables.
	_ = viper.ReadInConfig()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("STORAGE_TYPE", "local")
	viper.SetDefault("JWT_SECRET", "a-very-secret-key-that-should-be-changed")
    viper.SetDefault("TOKEN_EXPIRES_IN_HOURS", 72)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
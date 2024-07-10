package config

import (
	"github.com/spf13/viper"
)

const (
	defaultMongoDBURI = "mongodb://localhost:27017"
)

type AppConfig struct {
	App struct {
		Env  string
		Port uint16
	}

	DB struct {
		URI string
	}
}

var cfg *AppConfig

func Config() *AppConfig {
	if cfg == nil {
		loadConfig()
	}

	return cfg
}

func loadConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Ignore config file not found, perhaps environment variables will be used if present.
	_ = viper.ReadInConfig()

	setDefaultValues()

	cfg = &AppConfig{}

	// App.
	cfg.App.Env = viper.GetString("APP_ENV")
	cfg.App.Port = viper.GetUint16("APP_PORT")

	// Gin.
	//cfg.Gin.Mode = viper.GetString("GIN_MODE")

	//db
	cfg.DB.URI = viper.GetString("DB_URI")
}

func setDefaultValues() {
	viper.SetDefault("DB_URI", defaultMongoDBURI)
}

var Constcfg = AppConfig{
	App: struct {
		Env  string
		Port uint16
	}{
		Env:  "dev",
		Port: 8080,
	},
	DB: struct{ URI string }{
		URI: "mongodb://localhost:27017",
	},
}

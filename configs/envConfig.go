package configs

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Initilize this variable to access the env values
var EnvConfigs *envConfigs

// We will call this in main.go to load the env variables
func InitEnvConfigs() {
	EnvConfigs = loadEnvVariables()
}

// struct to map env values
type envConfigs struct {
	APP_ENV        string `mapstructure:"APP_ENV"`
	SERVER_PORT    int    `mapstructure:"PORT"`
	DB_HOST        string `mapstructure:"DB_HOST"`
	DB_NAME        string `mapstructure:"DB_NAME"`
	DB_USER        string `mapstructure:"DB_USER"`
	DB_PASSWORD    string `mapstructure:"DB_PASSWORD"`
	DB_PORT        int    `mapstructure:"DB_PORT"`
	JWT_EXPIRED    string `mapstructure:"JWT_EXPIRED"`
	JWT_SECRET_KEY string `mapstructure:"JWT_SECRET_KEY"`
	REDIS_HOST     string `mapstructure:"REDIS_HOST"`
	REDIS_PORT     int    `mapstructure:"REDIS_PORT"`
	REDIS_PASSWORD string `mapstructure:"REDIS_PASSWORD"`
	REDIS_DB       int    `mapstructure:"REDIS_DB"`
}

func StartBindEnvs() {
	viper.BindEnv("APP_ENV")
	viper.BindEnv("PORT")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("JWT_EXPIRED")
	viper.BindEnv("JWT_SECRET_KEY")
}

// Call to load the variables from env
func loadEnvVariables() (config *envConfigs) {
	// Tell viper the path/location of your env file. If it is root just add "."
	APP_ENV := os.Getenv("APP_ENV")
	if APP_ENV == "" {
		APP_ENV = "development"
		viper.Set("APP_ENV", APP_ENV)
	}

	viper.AutomaticEnv()
	viper.SetConfigType("env")
	if APP_ENV == "development" {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			StartBindEnvs()
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Error().Err(err).Msgf("Config file not found %s", err)
			} else {
				log.Fatal().Err(err).Msgf("Error reading config file from %s", err)
			}
		}
	} else {
		StartBindEnvs()
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("Error unmarshalling config")
	}

	return
}

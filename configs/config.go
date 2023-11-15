package configs

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var cfg Config
var doOnce sync.Once

type Config struct {
	Application struct {
		Log struct {
			Path string `mapstructure:"PATH"`
		}
		Graceful struct {
			MaxSecond time.Duration `mapstructure:"MAX_SECOND"`
		} `mapstructure:"GRACEFUL"`
		Cache struct {
			Redis struct {
				Host string `mapstructure:"HOST"`
				Port string `mapstructure:"PORT"`
			} `mapstructure:"REDIS"`
		} `mapstructure:"CACHE"`
		DB struct {
			Host     string `mapstructure:"HOST"`
			Port     string `mapstructure:"PORT"`
			Name     string `mapstructure:"NAME"`
			User     string `mapstructure:"USER"`
			Password string `mapstructure:"PASSWORD"`
		} `mapstructure:"DB"`
	} `mapstructure:"APPLICATION"`
}

func Get() Config {
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {

		log.Error().Err(err).Msgf("Error reading config file from %s", err)
	}

	doOnce.Do(func() {
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Error().Err(err).Msgf("Error unmarshal config file from %s", err)
		}
	})

	return cfg
}

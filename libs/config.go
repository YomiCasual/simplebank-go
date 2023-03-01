package lib

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)


type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	TestDBSource string `mapstructure:"TEST_DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	SymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig from file or environmentable varialbse
func LoadConfig(path string) (config Config, err error) {


	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err !=nil {
		return
	}

	err = viper.Unmarshal(&config);

	if HasError(err) {
		log.Fatal("Cannot load env", err)
	}

	return

}
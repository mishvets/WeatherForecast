package util

import "github.com/spf13/viper"

// Config stores all configuration of the application.
// The values are read by viper from aconfig file or environment variables.
type Config struct {
	DBDriver            string `mapstructure:"DB_DRIVER"`
	DBSource            string `mapstructure:"DB_SOURCE"`
	ServerAddress       string `mapstructure:"SERVER_ADDRESS"`
	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAdress   string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`
	RedisAddress        string `mapstructure:"REDIS_ADDRESS"`
	WeatherApiKey       string `mapstructure:"WEATHER_API_KEY"`
	WeatherApiUrl       string `mapstructure:"WEATHER_API_URL"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

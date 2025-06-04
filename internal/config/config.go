package config

import "github.com/spf13/viper"

type Config struct {
	Environment   string `mapstructure:"ENVIRONMENT"`
	MongoUrl      string `mapstructure:"MONGO_URL"`
	MongoUsername string `mapstructure:"MONGO_INITDB_ROOT_USERNAME"`
	MongoPassword string `mapstructure:"MONGO_INITDB_ROOT_PASSWORD"`
	SecretKey     string `mapstructure:"SECRETKEY"`
}

func NewConfig(path, env string) (*Config, error) {
	viper.AutomaticEnv()

	viper.SetConfigName(".env." + env)
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

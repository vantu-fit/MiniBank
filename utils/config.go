package utils

import "github.com/spf13/viper"

type Config struct {
	Evironment           string `mapstructure:"EVIRONMENT"`
	DBDriver             string `mapstructure:"DB_DRIVER"`
	BDSource             string `mapstructure:"DB_SOURCE"`
	MigrateUrl           string `mapstructure:"MIGRATE_URL"`
	ServerAddress        string `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  string `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration string `mapstructure:"REFRESH_TOKEN_DURATION"`
	GrpcServerAddress    string `mapstructure:"GRPC_SERVER_ADDRESS"`
	RedisAddress         string `mapstructure:"REDIS_ADDRESS"`
	EmailSenderName      string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

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

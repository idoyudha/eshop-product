package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	// Config
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"log"`
		AWS
		Redis
		Kafka
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	AWS struct {
		AwsRegion    string `env-required:"true" env:"AWS_REGION"`
		AwsAccessKey string `env-required:"true" env:"AWS_ACCESS_KEY"`
		AwsSecretey  string `env-required:"true" env:"AWS_SECRET_KEY"`
	}

	// Redis
	Redis struct {
		RedisURL      string `env-required:"true" env:"REDIS_URL"`
		RedisPassword string `env-required:"true" env:"REDIS_PASSWORD"`
	}

	// Log
	Log struct {
		Level string `yaml:"log_level"`
	}

	// Kafka
	Kafka struct {
		Broker string `env-required:"true" env:"KAFKA_BROKER"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	// Config
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		DynamoDB
		Redis
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

	// AWS DynamoDB
	DynamoDB struct {
		AwsRegion                  string `env-required:"true" env:"AWS_REGION"`
		AwsDynamoDBServiceEndpoint string `env-required:"true" env:"AWS_DYNAMO_DB_SERVICE_ENDPOINT"`
		AwsDynamoDBAccessKey       string `env-required:"true" env:"AWS_DYNAMO_DB_ACCESS_KEY"`
		AwsDynamoDBSecretey        string `env-required:"true" env:"AWS_DYNAMO_DB_SECRET_KEY"`
	}

	// Redis
	Redis struct {
		RedisURL      string `env-required:"true" env:"REDIS_URL"`
		RedisPassword string `env-required:"true" env:"REDIS_PASSWORD"`
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

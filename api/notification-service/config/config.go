package config

import "github.com/spf13/viper"

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	JwtSecretKey string `json:"jwt_secret_key"`
}

type PsqlDB struct {
	Host             string `json:"host"`
	Port             string `json:"port"`
	User             string `json:"user"`
	Password         string `json:"password"`
	DBName           string `json:"db_name"`
	DBConnectTimeout int    `json:"db_connect_timeout"`
	DBMaxOpen        int    `json:"db_max_open"`
	DBMaxIdle        int    `json:"db_max_idle"`
}

type RabbitMQ struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type EmailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Sending  string `json:"sending"`
	IsTLS    bool   `json:"is_tls"`
}

type Redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Config struct {
	App         App         `json:"app"`
	Psql        PsqlDB      `json:"psql"`
	RabbitMQ    RabbitMQ    `json:"rabbitmq"`
	Redis       Redis       `json:"redis"`
	EmailConfig EmailConfig `json:"email_config"`
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort: viper.GetString("APP_PORT"),
			AppEnv:  viper.GetString("APP_ENV"),

			JwtSecretKey: viper.GetString("JWT_SECRET_KEY"),
		},

		Psql: PsqlDB{
			Host:             viper.GetString("DB_HOST"),
			Port:             viper.GetString("DB_PORT"),
			User:             viper.GetString("DB_USER"),
			Password:         viper.GetString("DB_PASS"),
			DBName:           viper.GetString("DB_NAME"),
			DBConnectTimeout: viper.GetInt("DB_CONNECT_TIMEOUT"),
			DBMaxOpen:        viper.GetInt("DB_MAX_OPEN_CONNECTION"),
			DBMaxIdle:        viper.GetInt("DB_MAX_IDLE_CONNECTION"),
		},

		RabbitMQ: RabbitMQ{
			Host:     viper.GetString("RABBITMQ_HOST"),
			Port:     viper.GetString("RABBITMQ_PORT"),
			User:     viper.GetString("RABBITMQ_USER"),
			Password: viper.GetString("RABBITMQ_PASSWORD"),
		},

		EmailConfig: EmailConfig{
			Host:     viper.GetString("EMAIL_HOST"),
			Port:     viper.GetInt("EMAIL_PORT"),
			Username: viper.GetString("EMAIL_USERNAME"),
			Password: viper.GetString("EMAIL_PASSWORD"),
			Sending:  viper.GetString("EMAIL_SENDING"),
			IsTLS:    viper.GetBool("EMAIL_IS_TLS"),
		},

		Redis: Redis{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
		},
	}
}

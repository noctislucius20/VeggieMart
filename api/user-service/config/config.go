package config

import "github.com/spf13/viper"

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	JwtSecretKey string `json:"jwt_secret_key"`
	JwtIssuer    string `json:"jwt_issuer"`

	UrlUsersService string `json:"URL_USERS_SERVICE"`
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

type Supabase struct {
	Url    string `json:"url"`
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
}

type Redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Config struct {
	App      App      `json:"app"`
	Psql     PsqlDB   `json:"psql"`
	RabbitMQ RabbitMQ `json:"rabbitmq"`
	Storage  Supabase `json:"storage"`
	Redis    Redis    `json:"redis"`
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort: viper.GetString("APP_PORT"),
			AppEnv:  viper.GetString("APP_ENV"),

			JwtSecretKey: viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:    viper.GetString("JWT_ISSUER"),

			UrlUsersService: viper.GetString("URL_USERS_SERVICE"),
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
		Storage: Supabase{
			Url:    viper.GetString("SUPABASE_STORAGE_URL"),
			Key:    viper.GetString("SUPABASE_STORAGE_KEY"),
			Bucket: viper.GetString("SUPABASE_STORAGE_BUCKET"),
		},
		Redis: Redis{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
		},
	}
}

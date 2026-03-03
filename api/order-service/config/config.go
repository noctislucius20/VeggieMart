package config

import "github.com/spf13/viper"

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	JwtSecretKey string `json:"jwt_secret_key"`

	ServerTimeout int `json:"server_timeout"`

	ProductServiceUrl string `json:"product_service_url"`
	UserServiceUrl    string `json:"user_service_url"`

	LatitudeRef  string `json:"latitude_ref"`
	LongitudeRef string `json:"longitude_ref"`
	MaxDistance  int64  `json:"max_distance"`
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

type Redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Elasticsearch struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type PublisherName struct {
	ProductUpdateStock     string `json:"product_update_stock"`
	OrderCreate            string `json:"order_create"`
	EmailUpdateOrderStatus string `json:"email_update_order_status"`
	OrderPaymentSuccess    string `json:"order_payment_success"`
	OrderUpdateStatus      string `json:"order_update_status"`
}

type Config struct {
	App           App           `json:"app"`
	Psql          PsqlDB        `json:"psql"`
	RabbitMQ      RabbitMQ      `json:"rabbitmq"`
	Redis         Redis         `json:"redis"`
	Elasticsearch Elasticsearch `json:"elasticsearch"`
	PublisherName PublisherName `json:"publisher_name"`
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort: viper.GetString("APP_PORT"),
			AppEnv:  viper.GetString("APP_ENV"),

			JwtSecretKey:      viper.GetString("JWT_SECRET_KEY"),
			ServerTimeout:     viper.GetInt("SERVER_TIMEOUT"),
			ProductServiceUrl: viper.GetString("PRODUCT_SERVICE_URL"),
			UserServiceUrl:    viper.GetString("USER_SERVICE_URL"),

			LatitudeRef:  viper.GetString("LATITUDE_REF"),
			LongitudeRef: viper.GetString("LONGITUDE_REF"),
			MaxDistance:  viper.GetInt64("MAX_DISTANCE"),
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

		Redis: Redis{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
		},
		Elasticsearch: Elasticsearch{
			Host: viper.GetString("ELASTICSEARCH_HOST"),
			Port: viper.GetString("ELASTICSEARCH_PORT"),
		},
		PublisherName: PublisherName{
			ProductUpdateStock:     viper.GetString("PRODUCT_UPDATE_STOCK"),
			OrderCreate:            viper.GetString("ORDER_CREATE"),
			EmailUpdateOrderStatus: viper.GetString("EMAIL_UPDATE_ORDER_STATUS"),
			OrderPaymentSuccess:    viper.GetString("ORDER_PAYMENT_SUCCESS"),
			OrderUpdateStatus:      viper.GetString("ORDER_UPDATE_STATUS"),
		},
	}
}

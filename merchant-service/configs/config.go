package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	UrlProductService   string `json:"url_product_service"`
	UrlWarehouseService string `json:"url_warehouse_service"`
	UrlUserService      string `json:"url_user_service"`
}

type SqlDB struct {
	Host           string `json:"host"`
	Port           string `json:"port"`
	User           string `json:"user"`
	Password       string `json:"password"`
	DBName         string `json:"db_name"`
	DBMaxOpenConns int    `json:"db_max_open_conns"`
	DBMaxIdleConns int    `json:"db_max_idle_conns"`
}

type Redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type RabbitMQ struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Supabase struct {
	Url    string `json:"url"`
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
}

type Config struct {
	App      App      `json:"app"`
	SqlDB    SqlDB    `json:"sql_db"`
	Redis    Redis    `json:"redis"`
	RabbitMQ RabbitMQ `json:"rabbitmq"`
	Supabase Supabase `json:"supabase"`
}

func (r *RabbitMQ) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", r.Username, r.Password, r.Host, r.Port)
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort:             viper.GetString("APP_PORT"),
			AppEnv:              viper.GetString("APP_ENV"),
			UrlProductService:   viper.GetString("URL_PRODUCT_SERVICE"),
			UrlWarehouseService: viper.GetString("URL_WAREHOUSE_SERVICE"),
			UrlUserService:      viper.GetString("URL_USER_SERVICE"),
		},
		SqlDB: SqlDB{
			Host:           viper.GetString("DATABASE_HOST"),
			Port:           viper.GetString("DATABASE_PORT"),
			User:           viper.GetString("DATABASE_USER"),
			Password:       viper.GetString("DATABASE_PASSWORD"),
			DBName:         viper.GetString("DATABASE_DB_NAME"),
			DBMaxOpenConns: viper.GetInt("DATABASE_DB_MAX_OPEN_CONNECTION"),
			DBMaxIdleConns: viper.GetInt("DATABASE_DB_MAX_IDLE_CONNECTION"),
		},
		Redis: Redis{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
		},
		RabbitMQ: RabbitMQ{
			Host:     viper.GetString("RABBITMQ_HOST"),
			Port:     viper.GetString("RABBITMQ_PORT"),
			Username: viper.GetString("RABBITMQ_USERNAME"),
			Password: viper.GetString("RABBITMQ_PASSWORD"),
		},
		Supabase: Supabase{
			Url:    viper.GetString("SUPABASE_URL"),
			Key:    viper.GetString("SUPABASE_KEY"),
			Bucket: viper.GetString("SUPABASE_BUCKET"),
		},
	}
}

package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"go-fitness/external/logger/sl"
	ucfg "go.uber.org/config"
	"go.uber.org/fx"
	"log/slog"
	"os"
	"time"
)

type (
	ResultConfig struct {
		fx.Out

		Config   *Config
		Provider ucfg.Provider
	}

	Config struct {
		Env          string `yaml:"env"`
		HTTPServer   `yaml:"http_server"`
		WSServer     `yaml:"ws_server"`
		ENVState     `yaml:"env_state"`
		MongoDB      `yaml:"mongodb"`
		DB           `yaml:"db"`
		VideoService `yaml:"video_service"`
		JWT          string `yaml:"jwt_secret" env:"JWT_SECRET"`
	}

	DB struct {
		MaxOpenConns    int           `yaml:"max_open_conns" env:"MAX_OPEN_CONNS" env-default:"25"`
		MaxIdleConns    int           `yaml:"max_idle_conns"  env:"MAX_IDLE_CONNS" env-default:"25"`
		ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"MAX_lifetime_CONNS" env-default:"3m"`
		MysqlUser       string        `env:"MYSQL_USER" env-default:"root"`
		MysqlPassword   string        `env:"MYSQL_PASSWORD" env-default:"root"`
		MysqlHost       string        `env:"MYSQL_HOST" env-default:"localhost"`
		MysqlPort       string        `env:"MYSQL_PORT" env-default:"3306"`
		MysqlDBName     string        `env:"MYSQL_DBNAME" env-default:"rust"`
	}

	ENVState struct {
		Local string `yaml:"local" env-default:"local"`
		Dev   string `yaml:"dev" env-default:"dev"`
		Prod  string `yaml:"prod" env-default:"prod"`
	}

	WSServer struct {
		AppID   string `yaml:"app_id" env:"PUSHER_APP_ID"`
		Host    string `yaml:"address" env:"PUSHER_HOST" env-default:"localhost"`
		Port    string `yaml:"port" env:"PUSHER_PORT" env-default:"8080"`
		Cluster string `yaml:"cluster" env:"PUSHER_CLUSTER" env-default:"ap1"`
		Secret  string `yaml:"secret" env:"PUSHER_SECRET"`
		Key     string `yaml:"key" env:"PUSHER_KEY"`
		Secure  bool   `yaml:"secure" env:"PUSHER_SECURE" env-default:"false"`
	}

	HTTPServer struct {
		ApiPort     string        `yaml:"api_port" env:"API_PORT" env-default:":8082"`
		Timeout     time.Duration `yaml:"timeout" env-default:"60s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
		StoragePath string        `yaml:"storage_path" env:"STORAGE_PATH" env-default:"./storage"`
	}

	MongoDB struct {
		User          string `yaml:"mongodb_user" env:"MONGODB_USER" env-default:"root"`
		Password      string `yaml:"mongodb_password" env:"MONGODB_PASSWORD" env-default:"root"`
		Host          string `yaml:"mongodb_host" env:"MONGODB_HOST" env-default:"localhost"`
		Port          string `yaml:"mongodb_port" env:"MONGODB_PORT" env-default:"27017"`
		DBName        string `yaml:"mongodb_dbname" env:"MONGODB_DBNAME" env-default:"rust"`
		AuthDatabase  string `yaml:"mongodb_auth_database" env:"MONGODB_AUTH_DATABASE" env-default:"admin"`
		AuthMechanism string `yaml:"mongodb_auth_mechanism" env:"MONGODB_AUTH_MECHANISM" env-default:"SCRAM-SHA-1"`
	}

	VideoService struct {
		VideoPath                 string            `yaml:"video_path" env:"VIDEO_PATH" env-default:"videos"`
		TranscodeVideoWorkerCount int               `yaml:"transcode_worker_count" env:"TRANSCODE_WORKER_COUNT" env-default:"1"`
		Resolutions               map[string]string `yaml:"resolutions" env:"RESOLUTIONS" env-default:"640x360:360,854x480:480,1280x720:720,1920x1080:1080"`
	}
)

func NewConfig() (*Config, error) {
	log := slog.With(
		slog.String("op", "config.NewConfig"),
	)

	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	log.Info("Config path", sl.String("path", configPath))

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %s", err)
	}

	return &cfg, nil
}

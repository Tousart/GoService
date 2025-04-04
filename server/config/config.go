package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPFlags struct {
	HTTPConfigPath string
}

func ParseFlags() HTTPFlags {
	httpCfgPath := flag.String("config", "", "Path to http cfg")
	flag.Parse()
	return HTTPFlags{
		HTTPConfigPath: *httpCfgPath,
	}
}

func MustLoad(cfgPath string, cfg any) {
	if cfgPath == "" {
		log.Fatal("Config path is not set")
	}

	// Существует ли конфиг по указанному пути
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist by this path: %s", cfgPath)
	}

	// Читаем файл конфига и заполняем cfg
	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatalf("error reading config: %s", err)
	}
}

type RabbitMQ struct {
	Host      string `yaml:"host"`
	Port      uint16 `yaml:"port"`
	QueueName string `yaml:"queue_name"`
}

type HTTPConfig struct {
	Address string `yaml:"address"`
}

type Postgres struct {
	Host    string `yaml:"host"`
	Port    uint16 `yaml:"port"`
	DBName  string `yaml:"db_name"`
	SSLMode string `yaml:"sslmode"`
}

type Redis struct {
	Address         string        `yaml:"address"`
	DefaultPassword string        `yaml:"default_password"`
	DBNumber        int           `yaml:"db_number"`
	Duration        time.Duration `yaml:"duration"`
}

type ServerConfig struct {
	HTTPConfig `yaml:"http_server"`
	RabbitMQ   `yaml:"rabbit_mq"`
	Postgres   `yaml:"postgres_db"`
	Redis      `yaml:"redis_db"`
}

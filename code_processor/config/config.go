package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type CodeProcessorFlags struct {
	CodeProcessorConfigPath string
}

func ParseFlags() CodeProcessorFlags {
	processorCfgPath := flag.String("config", "", "Path to processors cfg")
	flag.Parse()
	return CodeProcessorFlags{
		CodeProcessorConfigPath: *processorCfgPath,
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

type CodeProcessorConfig struct {
	RabbitMQ `yaml:"rabbit_mq"`
}

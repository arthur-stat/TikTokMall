package conf

import (
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"gopkg.in/yaml.v2"
)

var (
	conf Config
	once sync.Once
)

type Config struct {
	Kitex  KitexConfig  `yaml:"kitex"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
	Logger LoggerConfig `yaml:"logger"`
}

type KitexConfig struct {
	Service       string `yaml:"service"`
	Address       string `yaml:"address"`
	LogFileName   string `yaml:"log_file_name"`
	LogMaxSize    int    `yaml:"log_max_size"`
	LogMaxBackups int    `yaml:"log_max_backups"`
	LogMaxAge     int    `yaml:"log_max_age"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

func Init() {
	once.Do(func() {
		// TODO: Add configuration loading from file
	})
}

func GetConf() *Config {
	return &conf
}

func LogLevel() klog.Level {
	switch conf.Logger.Level {
	case "debug":
		return klog.LevelDebug
	case "info":
		return klog.LevelInfo
	case "warn":
		return klog.LevelWarn
	case "error":
		return klog.LevelError
	default:
		return klog.LevelInfo
	}
}
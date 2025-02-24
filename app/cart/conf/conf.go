package conf

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kr/pretty"
	"github.com/spf13/viper"
	"gopkg.in/validator.v2"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Env        string           `mapstructure:"env"`
	Service    ServiceConfig    `mapstructure:"service"`
	MySQL      MySQLConfig      `mapstructure:"mysql"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Registry   RegistryConfig   `mapstructure:"registry"`
	Log        LogConfig        `mapstructure:"log"`
	Jaeger     JaegerConfig     `mapstructure:"jaeger"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
}

type ServiceConfig struct {
	Name     string `mapstructure:"name"`
	Port     int    `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
}

type MySQLConfig struct {
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Database     string `mapstructure:"database"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RegistryConfig struct {
	RegistryAddress []string `mapstructure:"registry_address"`
	Username        string   `mapstructure:"username"`
	Password        string   `mapstructure:"password"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	File       string `mapstructure:"file"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type JaegerConfig struct {
	Host         string  `mapstructure:"host"`
	Port         int     `mapstructure:"port"`
	SamplerType  string  `mapstructure:"sampler_type"`
	SamplerParam float64 `mapstructure:"sampler_param"`
	LogSpans     bool    `mapstructure:"log_spans"`
}

type PrometheusConfig struct {
	Port int `mapstructure:"port"`
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	once.Do(initConf)
	return conf
}

// Init 初始化配置
func Init() error {
	once.Do(initConf)
	return nil
}

func initConf() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join("conf", GetEnv()))

	if err := viper.ReadInConfig(); err != nil {
		klog.Fatalf("读取配置文件失败: %v", err)
	}

	conf = new(Config)
	if err := viper.Unmarshal(conf); err != nil {
		klog.Fatalf("解析配置文件失败: %v", err)
	}

	if err := validator.Validate(conf); err != nil {
		klog.Fatalf("验证配置失败: %v", err)
	}

	conf.Env = GetEnv()
	pretty.Printf("%+v\n", conf)
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "test"
	}
	return e
}

func LogLevel() klog.Level {
	level := GetConfig().Service.LogLevel
	switch level {
	case "trace":
		return klog.LevelTrace
	case "debug":
		return klog.LevelDebug
	case "info":
		return klog.LevelInfo
	case "notice":
		return klog.LevelNotice
	case "warn":
		return klog.LevelWarn
	case "error":
		return klog.LevelError
	case "fatal":
		return klog.LevelFatal
	default:
		return klog.LevelInfo
	}
}

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
	HTTP       HTTPConfig       `mapstructure:"http"`
	RPC        RPCConfig        `mapstructure:"rpc"`
	Kitex      KitexConfig      `mapstructure:"kitex"`
	MySQL      MySQLConfig      `mapstructure:"mysql"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Registry   RegistryConfig   `mapstructure:"registry"`
	Log        LogConfig        `mapstructure:"log"`
	Jaeger     JaegerConfig     `mapstructure:"jaeger"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
}

type HTTPConfig struct {
	ServiceName string `mapstructure:"service_name"`
	Address     string `mapstructure:"address"`
}

type RPCConfig struct {
	ServiceName string `mapstructure:"service_name"`
	Port       int    `mapstructure:"port"`
}

type KitexConfig struct {
	LogFileName   string `mapstructure:"log_file_name"`
	LogMaxSize    int    `mapstructure:"log_max_size"`
	LogMaxBackups int    `mapstructure:"log_max_backups"`
	LogMaxAge     int    `mapstructure:"log_max_age"`
	LogLevel      string `mapstructure:"log_level"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
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
	Level string `mapstructure:"level"`
}

type JaegerConfig struct {
	Host         string  `mapstructure:"host"`
	Port         int     `mapstructure:"port"`
	SamplerType  string  `mapstructure:"sampler_type"`
	SamplerParam float64 `mapstructure:"sampler_param"`
	LogSpans     bool    `mapstructure:"log_spans"`
}

type PrometheusConfig struct {
	Port int    `mapstructure:"port"`
	Path string `mapstructure:"path"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join("conf", GetEnv()))

	if err := viper.ReadInConfig(); err != nil {
		klog.Fatalf("failed to read config file: %v", err)
	}

	conf = new(Config)
	if err := viper.Unmarshal(conf); err != nil {
		klog.Fatalf("failed to unmarshal config: %v", err)
	}

	if err := validator.Validate(conf); err != nil {
		klog.Fatalf("config validation failed: %v", err)
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
	level := GetConf().Kitex.LogLevel
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
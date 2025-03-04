package conf

import (
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/spf13/viper"
)

var (
	config Config
	once   sync.Once
)

type Config struct {
	Service    ServiceConfig    `mapstructure:"service"`
	MySQL      MySQLConfig      `mapstructure:"mysql"`
	Registry   RegistryConfig   `mapstructure:"registry"`
	Log        LogConfig        `mapstructure:"log"`
	Jaeger     JaegerConfig     `mapstructure:"jaeger"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	TLS        TLSConfig        `mapstructure:"tls"`
}

type ServiceConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type RegistryConfig struct {
	RegistryAddress []string `mapstructure:"registry_address"`
	Username        string   `mapstructure:"username"`
	Password        string   `mapstructure:"password"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
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

type TLSConfig struct {
	Enable         bool   `mapstructure:"enable"`
	CACertPath     string `mapstructure:"ca_cert_path"`
	ServerCertPath string `mapstructure:"server_cert_path"`
	ServerKeyPath  string `mapstructure:"server_key_path"`
	ClientCertPath string `mapstructure:"client_cert_path"`
	ClientKeyPath  string `mapstructure:"client_key_path"`
}

func Init() {
	once.Do(func() {
		viper.SetConfigName("conf")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./conf/dev")

		if err := viper.ReadInConfig(); err != nil {
			klog.Fatalf("读取配置文件失败: %v", err)
		}

		if err := viper.Unmarshal(&config); err != nil {
			klog.Fatalf("解析配置文件失败: %v", err)
		}
	})
}

func GetConfig() *Config {
	return &config
}

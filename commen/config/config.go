package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

var (
	ConfOnce   sync.Once
	instance   *Config
	configPath = "commen/config/config.yaml"
)

// Config 代表应用程序的配置结构体
type Config struct {
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	RabbitMQ   RabbitMQConfig   `yaml:"rabbitmq"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
	Logstash   LogstashConfig   `yaml:"logstash"`
	JWT        JWTConfig        `yaml:"jwt"`
	GA         GAConfig         `yaml:"GAParams"`
}

// DatabaseConfig 代表数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// RedisConfig 代表Redis配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// RabbitMQConfig 代表RabbitMQ配置
type RabbitMQConfig struct {
	Addr string `yaml:"addr"`
}

// PrometheusConfig 代表Prometheus配置
type PrometheusConfig struct {
	Port int `yaml:"port"`
}

// LogstashConfig 代表Logstash配置
type LogstashConfig struct {
	Addr string `yaml:"addr"`
}

// JWTConfig 代表JWT配置
type JWTConfig struct {
	SecretKey      string `yaml:"secret_key"`
	ExpirationTime string `yaml:"expiration_time"`
}

type GAConfig struct {
	PopulationSize int     `yaml:"PopulationSize"`
	CrossoverRate  float64 `yaml:"CrossoverRate"`
	MutationRate   float64 `yaml:"MutationRate"`
	MaxGenerations int     `yaml:"MaxGenerations"`
	TournamentSize int     `yaml:"TournamentSize"`
	ElitismCount   int     `yaml:"ElitismCount"`
}

// LoadConfig 从YAML文件中加载配置
func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}

// GetConfig 返回配置的单例实例
func GetConfig() *Config {
	ConfOnce.Do(func() {
		conf, err := loadConfig()
		if err != nil {
			log.Fatalf("error loading config: %v", err)
		}
		instance = conf
	})
	return instance
}

// Example usage
//func main() {
//	config, err := LoadConfig()
//	if err != nil {
//		log.Fatalf("error loading config: %v", err)
//	}
//
//	fmt.Printf("Database Host: %s\n", config.Database.Host)
//	fmt.Printf("Redis Address: %s\n", config.Redis.Addr)
//	fmt.Printf("RabbitMQ Address: %s\n", config.RabbitMQ.Addr)
//	fmt.Printf("Prometheus Port: %d\n", config.Prometheus.Port)
//	fmt.Printf("JWT Secret Key: %s\n", config.JWT.SecretKey)
//}

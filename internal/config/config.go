package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 全局配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	OSS      OSSConfig      `mapstructure:"oss"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	RPC      RPCConfig      `mapstructure:"rpc"`
	Security SecurityConfig `mapstructure:"security"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	LogLevel     string `mapstructure:"log_level"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// KafkaConfig Kafka配置
type KafkaConfig struct {
	Brokers     []string `mapstructure:"brokers"`
	TopicPrefix string   `mapstructure:"topic_prefix"`
	GroupID     string   `mapstructure:"group_id"`
	RequireACKs string   `mapstructure:"require_acks"`
	Partitioner string   `mapstructure:"partitioner"`
}

// OSSConfig 对象存储配置
type OSSConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Bucket          string `mapstructure:"bucket"`
	BucketDirectory string `mapstructure:"bucket_directory"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireTime  int    `mapstructure:"expire_time"`  // 小时
	RefreshTime int    `mapstructure:"refresh_time"` // 小时
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

// RPCConfig RPC配置
type RPCConfig struct {
	Services map[string]ServiceConfig `mapstructure:"services"`
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	PasswordEncrypted bool   `mapstructure:"password_encrypted"`
	BCryptCost        int    `mapstructure:"bcrypt_cost"`
	AllowedOrigins    string `mapstructure:"allowed_origins"`
	RateLimit         int    `mapstructure:"rate_limit"`
}

var globalConfig *Config

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	globalConfig = config
	return config, nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}

// setDefaults 设置默认值
func setDefaults() {
	viper.SetDefault("server.port", "8888")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("server.idle_timeout", 60)

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.log_level", "error")

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)

	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.topic_prefix", "douyin")
	viper.SetDefault("kafka.group_id", "douyin-group")
	viper.SetDefault("kafka.require_acks", "WaitForAll")
	viper.SetDefault("kafka.partitioner", "NewRandomPartitioner")

	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expire_time", 24)
	viper.SetDefault("jwt.refresh_time", 168)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 28)

	viper.SetDefault("security.password_encrypted", true)
	viper.SetDefault("security.bcrypt_cost", 12)
	viper.SetDefault("security.rate_limit", 1000)
}

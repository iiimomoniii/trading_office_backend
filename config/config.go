package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Fiber    FiberConfig    `mapstructure:"fiber"`
	Database DBConfig       `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Claude   ClaudeConfig   `mapstructure:"claude"`
	Exchange ExchangeConfig `mapstructure:"exchange"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
}

type FiberConfig struct {
	Address        string `mapstructure:"address"`
	Port           int    `mapstructure:"port"`
	ReadTimeout    int    `mapstructure:"readTimeout"`
	WriteTimeout   int    `mapstructure:"writeTimeout"`
	IdleTimeout    int    `mapstructure:"idleTimeout"`
	ReadBufferSize int    `mapstructure:"readBufferSize"`
	BodyLimitSize  int    `mapstructure:"bodyLimitSize"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn int    `mapstructure:"expiresIn"`
}

type ClaudeConfig struct {
	APIKey string `mapstructure:"apiKey"`
	Model  string `mapstructure:"model"`
}

type ExchangeConfig struct {
	BaseURL       string `mapstructure:"baseUrl"`
	FetchInterval int    `mapstructure:"fetchInterval"`
}

// LoadConfig — โหลด .env.{env} ก่อน แล้ว override ด้วย config.yaml
func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// 1. โหลด .env.{env} ด้วย godotenv (load เข้า OS env vars)
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		fmt.Printf("[config] warning: %s not found: %v\n", envFile, err)
		// fallback .env
		if err2 := godotenv.Load(".env"); err2 != nil {
			fmt.Printf("[config] warning: .env not found\n")
		}
	} else {
		fmt.Printf("[config] loaded %s\n", envFile)
	}

	// 2. โหลด config.yaml ด้วย viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("[config] warning: config.yaml not found, using env vars only\n")
	} else {
		fmt.Printf("[config] loaded config.yaml\n")
	}

	// 3. map OS env vars → viper keys
	overrideFromEnv()

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("[config] unmarshal error: %w", err)
	}

	setDefaults(cfg)

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func overrideFromEnv() {
	envMap := map[string]string{
		"DB_HOST":           "database.host",
		"DB_USER":           "database.user",
		"DB_PASS":           "database.password",
		"DB_NAME":           "database.dbname",
		"DB_SSLMODE":        "database.sslmode",
		"APP_NAME":          "app.name",
		"CLAUDE_API_KEY":    "claude.apiKey",
		"CLAUDE_MODEL":      "claude.model",
		"EXCHANGE_BASE_URL": "exchange.baseUrl",
	}
	for envKey, viperKey := range envMap {
		if v := os.Getenv(envKey); v != "" {
			viper.Set(viperKey, v)
		}
	}

	// int overrides
	if v := os.Getenv("DB_PORT"); v != "" {
		viper.Set("database.port", v)
	}
	if v := os.Getenv("APP_PORT"); v != "" {
		viper.Set("fiber.port", v)
		viper.Set("fiber.address", fmt.Sprintf("0.0.0.0:%s", v))
	}
	if v := os.Getenv("EXCHANGE_FETCH_INTERVAL"); v != "" {
		viper.Set("exchange.fetchInterval", v)
	}
}

func setDefaults(c *Config) {
	if c.Fiber.Address == "" {
		c.Fiber.Address = "0.0.0.0:8080"
	}
	if c.Fiber.Port == 0 {
		c.Fiber.Port = 8080
	}
	if c.Fiber.ReadTimeout == 0 {
		c.Fiber.ReadTimeout = 30000
	}
	if c.Fiber.WriteTimeout == 0 {
		c.Fiber.WriteTimeout = 30000
	}
	if c.Fiber.IdleTimeout == 0 {
		c.Fiber.IdleTimeout = 30000
	}
	if c.Fiber.ReadBufferSize == 0 {
		c.Fiber.ReadBufferSize = 8192
	}
	if c.Fiber.BodyLimitSize == 0 {
		c.Fiber.BodyLimitSize = 10485760
	}
	if c.Database.SSLMode == "" {
		c.Database.SSLMode = "disable"
	}
	if c.Database.Port == 0 {
		c.Database.Port = 5432
	}
	if c.Claude.Model == "" {
		c.Claude.Model = "claude-sonnet-4-20250514"
	}
	if c.Exchange.BaseURL == "" {
		c.Exchange.BaseURL = "https://api.binance.com"
	}
	if c.Exchange.FetchInterval == 0 {
		c.Exchange.FetchInterval = 5
	}
	if c.JWT.ExpiresIn == 0 {
		c.JWT.ExpiresIn = 3600
	}
}

func validate(c *Config) error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	return nil
}

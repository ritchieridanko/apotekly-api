package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string
	}

	Server struct {
		Host            string
		Port            int
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		ShutdownTimeout time.Duration
	}

	JWT struct {
		Secret string
	}

	Database struct {
		Host            string
		Port            int
		User            string
		Pass            string
		Name            string
		SSLMode         string
		DSN             string
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime time.Duration
	}

	Storage struct {
		Provider  string
		Bucket    string
		APIKey    string
		APISecret string
	}

	Image struct {
		MaxSizeBytes int
		AllowedTypes []string
	}

	Tracer struct {
		Endpoint string
	}
}

func Load(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	v.RegisterAlias("server.read_timeout", "server.readtimeout")
	v.RegisterAlias("server.write_timeout", "server.writetimeout")
	v.RegisterAlias("server.shutdown_timeout", "server.shutdowntimeout")
	v.RegisterAlias("database.conn_max_lifetime", "database.connmaxlifetime")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.Database.DSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	return &cfg, nil
}

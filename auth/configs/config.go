package configs

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App
	Auth
	OAuth
	Server
	Client
	Database
	Cache
	Tracer
}

type App struct {
	Name string
	Env  string
}

type Auth struct {
	BCrypt struct {
		Cost int
	}

	JWT struct {
		Issuer    string
		Audiences []string
		Secret    string
		Duration  time.Duration
	}

	TokenDuration struct {
		Session      time.Duration
		Reset        time.Duration
		Verification time.Duration
		EmailChange  time.Duration
	}
}

type OAuth struct {
	Google struct {
		ClientID    string
		Secret      string
		RedirectURL string
	}

	Microsoft struct {
		ClientID    string
		Secret      string
		RedirectURL string
	}

	Duration struct {
		CodeExchange time.Duration
	}
}

type Server struct {
	Host string
	Port int

	Timeout struct {
		Read     time.Duration
		Write    time.Duration
		Shutdown time.Duration
	}
}

type Client struct {
	BaseURL string
}

type Database struct {
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

type Cache struct {
	Host       string
	Port       int
	Pass       string
	MaxRetries int
	BaseDelay  int
}

type Tracer struct {
	Endpoint string
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

package configs

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      `mapstructure:"app"`
	Auth     `mapstructure:"auth"`
	OAuth    `mapstructure:"oauth"`
	Server   `mapstructure:"server"`
	Client   `mapstructure:"client"`
	Database `mapstructure:"database"`
	Cache    `mapstructure:"cache"`
	Tracer   `mapstructure:"tracer"`
	Broker   `mapstructure:"broker"`
}

type App struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type Auth struct {
	BCrypt struct {
		Cost int `mapstructure:"cost"`
	} `mapstructure:"bcrypt"`

	JWT struct {
		Issuer    string        `mapstructure:"issuer"`
		Audiences []string      `mapstructure:"audiences"`
		Secret    string        `mapstructure:"secret"`
		Duration  time.Duration `mapstructure:"duration"`
	} `mapstructure:"jwt"`

	TokenDuration struct {
		Session      time.Duration `mapstructure:"session"`
		Reset        time.Duration `mapstructure:"reset"`
		Verification time.Duration `mapstructure:"verification"`
		EmailChange  time.Duration `mapstructure:"email_change"`
	} `mapstructure:"token_duration"`
}

type OAuth struct {
	Google struct {
		ClientID    string `mapstructure:"client_id"`
		Secret      string `mapstructure:"secret"`
		RedirectURL string `mapstructure:"redirect_url"`
	} `mapstructure:"google"`

	Microsoft struct {
		ClientID    string `mapstructure:"client_id"`
		Secret      string `mapstructure:"secret"`
		RedirectURL string `mapstructure:"redirect_url"`
	} `mapstructure:"microsoft"`

	Duration struct {
		CodeExchange time.Duration `mapstructure:"code_exchange"`
	} `mapstructure:"duration"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	Timeout struct {
		Read     time.Duration `mapstructure:"read"`
		Write    time.Duration `mapstructure:"write"`
		Shutdown time.Duration `mapstructure:"shutdown"`
	} `mapstructure:"timeout"`
}

type Client struct {
	BaseURL string `mapstructure:"base_url"`
}

type Database struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Pass            string `mapstructure:"pass"`
	Name            string `mapstructure:"name"`
	SSLMode         string `mapstructure:"ssl_mode"`
	DSN             string
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type Cache struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Pass       string `mapstructure:"pass"`
	MaxRetries int    `mapstructure:"max_retries"`
	BaseDelay  int    `mapstructure:"base_delay"`
}

type Tracer struct {
	Endpoint string `mapstructure:"endpoint"`
}

type Broker struct {
	Brokers string `mapstructure:"brokers"`

	Timeout struct {
		Batch time.Duration `mapstructure:"batch"`
	} `mapstructure:"timeout"`
}

func Load(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	// read YAML config
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// load env variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
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

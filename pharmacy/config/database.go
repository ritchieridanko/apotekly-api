package config

type dbConfig struct {
	Host            string
	Port            string
	User            string
	Pass            string
	Name            string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

var dbCfg *dbConfig

func loadDBConfig() {
	dbCfg = &dbConfig{
		Host:            getEnv("DB_HOST"),
		Port:            getEnv("DB_PORT"),
		User:            getEnv("DB_USER"),
		Pass:            getEnv("DB_PASS"),
		Name:            getEnv("DB_NAME"),
		SSLMode:         getEnv("DB_SSL_MODE"),
		MaxIdleConns:    getNumberEnv("DB_MAX_IDLE_CONNS"),
		MaxOpenConns:    getNumberEnv("DB_MAX_OPEN_CONNS"),
		ConnMaxLifetime: getNumberEnv("DB_CONN_MAX_LIFETIME"),
	}
}

func DBGetHost() (host string) {
	return dbCfg.Host
}

func DBGetPort() (port string) {
	return dbCfg.Port
}

func DBGetUser() (username string) {
	return dbCfg.User
}

func DBGetPass() (password string) {
	return dbCfg.Pass
}

func DBGetName() (name string) {
	return dbCfg.Name
}

func DBGetSSLMode() (mode string) {
	return dbCfg.SSLMode
}

func DBGetMaxIdleConns() (max int) {
	return dbCfg.MaxIdleConns
}

func DBGetMaxOpenConns() (max int) {
	return dbCfg.MaxOpenConns
}

func DBGetConnMaxLifetime() (max int) {
	return dbCfg.ConnMaxLifetime
}

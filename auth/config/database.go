package config

type dbConfig struct {
	host            string
	port            string
	user            string
	pass            string
	name            string
	sslMode         string
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime int
}

var dbCfg *dbConfig

func LoadDBConfig() {
	dbCfg = &dbConfig{
		host:            GetEnv("DB_HOST"),
		port:            GetEnv("DB_PORT"),
		user:            GetEnv("DB_USER"),
		pass:            GetEnv("DB_PASS"),
		name:            GetEnv("DB_NAME"),
		sslMode:         GetEnv("DB_SSL_MODE"),
		maxIdleConns:    GetNumberEnv("DB_MAX_IDLE_CONNS"),
		maxOpenConns:    GetNumberEnv("DB_MAX_OPEN_CONNS"),
		connMaxLifetime: GetNumberEnv("DB_CONN_MAX_LIFETIME"),
	}
}

func GetDBHost() (host string) {
	return dbCfg.host
}

func GetDBPort() (port string) {
	return dbCfg.port
}

func GetDBUser() (username string) {
	return dbCfg.user
}

func GetDBPass() (password string) {
	return dbCfg.pass
}

func GetDBName() (name string) {
	return dbCfg.name
}

func GetDBSSLMode() (mode string) {
	return dbCfg.sslMode
}

func GetDBMaxIdleConns() (max int) {
	return dbCfg.maxIdleConns
}

func GetDBMaxOpenConns() (max int) {
	return dbCfg.maxOpenConns
}

func GetDBConnMaxLifetime() (max int) {
	return dbCfg.connMaxLifetime
}

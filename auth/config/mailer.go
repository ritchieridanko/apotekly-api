package config

type mailerConfig struct {
	host string
	port int
	user string
	pass string
}

var mailerCfg *mailerConfig

func LoadMailerConfig() {
	mailerCfg = &mailerConfig{
		host: GetEnv("MAILER_HOST"),
		port: GetNumberEnv("MAILER_PORT"),
		user: GetEnv("MAILER_USER"),
		pass: GetEnv("MAILER_PASS"),
	}
}

func GetMailerHost() (host string) {
	return mailerCfg.host
}

func GetMailerPort() (port int) {
	return mailerCfg.port
}

func GetMailerUser() (user string) {
	return mailerCfg.user
}

func GetMailerPass() (pass string) {
	return mailerCfg.pass
}

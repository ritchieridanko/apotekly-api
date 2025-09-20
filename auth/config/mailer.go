package config

type mailerConfig struct {
	Host string
	Port int
	User string
	Pass string
}

var mailerCfg *mailerConfig

func loadMailerConfig() {
	mailerCfg = &mailerConfig{
		Host: getEnv("MAILER_HOST"),
		Port: getNumberEnv("MAILER_PORT"),
		User: getEnv("MAILER_USER"),
		Pass: getEnv("MAILER_PASS"),
	}
}

func MailerGetHost() (host string) {
	return mailerCfg.Host
}

func MailerGetPort() (port int) {
	return mailerCfg.Port
}

func MailerGetUser() (user string) {
	return mailerCfg.User
}

func MailerGetPass() (pass string) {
	return mailerCfg.Pass
}

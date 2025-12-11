package smtp

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// Estructura para el archivo JSON de configuración
type Config struct {
	SMTP SMTPConfigJSON `json:"smtp"`
}

type SMTPConfigJSON struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

// Crear SMTPConfig desde la configuración JSON principal
func NewSMTPConfigFromJSON(host string, port int, username, password, from string) *SMTPConfig {
	return &SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

package config

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Addr string
	Port int
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	Name     string
	User     string
	Password string
}

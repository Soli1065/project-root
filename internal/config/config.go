package config

// Configuration struct holds the configuration for the application
type Configuration struct {
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
}

// NewConfiguration creates a new configuration instance with default values
func NewConfiguration() *Configuration {
	return &Configuration{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "username",
		DBPassword: "password",
		DBName:     "database",
	}
}

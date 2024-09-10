package config

type Config struct {
	AppName          string `mapstructure:"APP_NAME"`
	Env              string `mapstructure:"ENV"`
	Port             string `mapstructure:"PORT"`
	AdminPass        string `mapstructure:"ADMIN_PASS"`
	JwtSecret        string `mapstructure:"JWT_SECRET"`
	JwtTime          int    `mapstructure:"JWT_TIME"`
	CorsMaxAge       int    `mapstructure:"CORS_MAX_AGE"`
	CorsAllowOrigins string `mapstructure:"CORS_ALLOW_ORIGINS"`
	CorsAllowMethods string `mapstructure:"CORS_ALLOW_METHODS"`
	RateLimitMax     int    `mapstructure:"RATE_LIMIT_MAX"`
	RateLimitExp     int    `mapstructure:"RATE_LIMIT_EXPIRATION"`
	CacheExp         int    `mapstructure:"CACHE_EXPIRATION"`
}

var AppConfig *Config

func LoadConfig() error {
	return nil
}

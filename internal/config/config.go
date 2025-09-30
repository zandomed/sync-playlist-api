package config

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Enviroment string

const (
	Development Enviroment = "development"
	Production  Enviroment = "production"
	Staging     Enviroment = "staging"
	Test        Enviroment = "test"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Spotify  SpotifyConfig
	Apple    AppleConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port            string
	Host            string
	Environment     Enviroment
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type SpotifyConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type AppleConfig struct {
	TeamID      string
	KeyID       string
	PrivateKey  string
	RedirectURL string
}

type JWTConfig struct {
	Secret                string
	ExpirationTime        time.Duration
	RefreshExpirationTime time.Duration
}

var (
	instance *Config
	once     sync.Once
)

func (c Config) IsDevelopment() bool {
	return c.Server.Environment != "production"
}

// Get returns the singleton instance of the config
func Get() *Config {
	once.Do(func() {
		config, err := load()
		if err != nil {
			panic("Failed to load config: " + err.Error())
		}
		instance = config
	})
	return instance
}

// Load carga la configuración desde variables de entorno (deprecated, use Get() instead)
func Load() (*Config, error) {
	return load()
}

// load carga la configuración desde variables de entorno
func load() (*Config, error) {
	// Cargar .env si existe (para desarrollo local)
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "9000"),
			Host:            getEnv("HOST", "localhost"),
			Environment:     Enviroment(getEnv("ENVIRONMENT", string(Development))),
			ReadTimeout:     parseDuration(getEnv("READ_TIMEOUT", "10s")),
			WriteTimeout:    parseDuration(getEnv("WRITE_TIMEOUT", "10s")),
			ShutdownTimeout: parseDuration(getEnv("SHUTDOWN_TIMEOUT", "5s")),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "sync-playlist"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", "password"),
			DB:       parseInt(getEnv("REDIS_DB", "0")),
		},
		Spotify: SpotifyConfig{
			ClientID:     getEnv("SPOTIFY_CLIENT_ID", ""),
			ClientSecret: getEnv("SPOTIFY_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("SPOTIFY_REDIRECT_URL", "http://localhost:9000/auth/spotify/callback"),
		},
		Apple: AppleConfig{
			TeamID:      getEnv("APPLE_TEAM_ID", ""),
			KeyID:       getEnv("APPLE_KEY_ID", ""),
			PrivateKey:  getEnv("APPLE_PRIVATE_KEY", ""),
			RedirectURL: getEnv("APPLE_REDIRECT_URL", "http://localhost:9000/auth/apple/callback"),
		},
		JWT: JWTConfig{
			Secret:                getEnv("JWT_SECRET", "your-secret-key"),
			ExpirationTime:        parseDuration(getEnv("JWT_EXPIRATION", "24h")),
			RefreshExpirationTime: parseDuration(getEnv("JWT_REFRESH_EXPIRATION", "100h")),
		},
	}, nil
}

// Funciones helper
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func parseInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func parseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

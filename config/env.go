package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    string
	FirebaseCredentialsPath string
	FirebaseCredentialsJSON string
	RedisEnabled            bool
	UpstashRedisRestURL     string
	UpstashRedisRestToken   string
}

// LoadFromEnv loads config from environment variables only. No interactive prompts.
// Use for Vercel/serverless where stdin is unavailable.
func LoadFromEnv() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:                    os.Getenv("PORT"),
		FirebaseCredentialsPath: os.Getenv("FIREBASE_CREDENTIALS_PATH"),
		FirebaseCredentialsJSON: os.Getenv("FIREBASE_CREDENTIALS_JSON"),
		UpstashRedisRestURL:     os.Getenv("UPSTASH_REDIS_REST_URL"),
		UpstashRedisRestToken:   os.Getenv("UPSTASH_REDIS_REST_TOKEN"),
		RedisEnabled:            os.Getenv("UPSTASH_REDIS_REST_URL") != "" && os.Getenv("UPSTASH_REDIS_REST_TOKEN") != "",
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg
}

func (c *Config) FirebaseEnabled() bool {
	return c.FirebaseCredentialsPath != "" || c.FirebaseCredentialsJSON != ""
}

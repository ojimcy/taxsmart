package config

import (
	"os"
)

type Config struct {
	Port                    string
	SupabaseURL             string
	SupabaseKey             string
	SupabaseJWTSecret       string
	FirebaseCredentialsFile string
	AIProvider              string // "gemini", "openai", "claude"
	AIAPIKey                string
	Environment             string
}

func Load() *Config {
	return &Config{
		Port:                    getEnv("PORT", "8080"),
		SupabaseURL:             getEnv("SUPABASE_URL", ""),
		SupabaseKey:             getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseJWTSecret:       getEnv("SUPABASE_JWT_SECRET", ""),
		FirebaseCredentialsFile: getEnv("FIREBASE_CREDENTIALS_FILE", ""),
		AIProvider:              getEnv("AI_PROVIDER", "gemini"),
		AIAPIKey:                getEnv("AI_API_KEY", ""),
		Environment:             getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

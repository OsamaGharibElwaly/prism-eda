// Package config loads runtime settings from the environment.
package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	MaxUploadBytes int64
	SamplePath     string
	WebDir         string
}

func Load() Config {
	return Config{
		Port:           env("PORT", "8123"),
		MaxUploadBytes: envInt64("PRISM_MAX_UPLOAD_MB", 25) * 1024 * 1024,
		SamplePath:     env("PRISM_SAMPLE", "data/Housing.csv"),
		WebDir:         env("PRISM_WEB_DIR", "web/dist"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt64(key string, fallback int64) int64 {
	v, err := strconv.ParseInt(os.Getenv(key), 10, 64)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Env              string
	AppPort          int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     int
	TLSCertPath      string
	TLSKeyPath       string
	EncryptionKeyB64 string
	StoreSlug        string
	StoreName        string
	MerchantUser     string
	MerchantPass     string
	CustomerUser     string
	CustomerPass     string
}

func Load() (Config, error) {
	cfg := Config{
		Env:              getenv("APP_ENV", "local"),
		AppPort:          getenvInt("APP_PORT", 8443),
		PostgresUser:     getenv("POSTGRES_USER", "nimble"),
		PostgresPassword: getenv("POSTGRES_PASSWORD", "nimble_pw"),
		PostgresDB:       getenv("POSTGRES_DB", "nimble"),
		PostgresHost:     getenv("POSTGRES_HOST", "db"),
		PostgresPort:     getenvInt("POSTGRES_PORT", 5432),
		TLSCertPath:      getenv("APP_TLS_CERT", ""),
		TLSKeyPath:       getenv("APP_TLS_KEY", ""),
		EncryptionKeyB64: getenv("APP_ENCRYPTION_KEY", ""),
		StoreSlug:        getenv("STORE_SLUG", "demo"),
		StoreName:        getenv("STORE_NAME", "Demo Pet Store"),
		MerchantUser:     getenv("MERCHANT_USERNAME", "merchant_demo"),
		MerchantPass:     getenv("MERCHANT_PASSWORD", "merchant_demo_pw"),
		CustomerUser:     getenv("CUSTOMER_USERNAME", "customer_demo"),
		CustomerPass:     getenv("CUSTOMER_PASSWORD", "customer_demo_pw"),
	}

	if cfg.EncryptionKeyB64 == "" {
		return cfg, fmt.Errorf("APP_ENCRYPTION_KEY is required")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getenvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return parsed
}

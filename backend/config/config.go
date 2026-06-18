package config

import (
	"log"
	"os"
)

type Config struct {
	DB         DBConfig
	JWT        JWTConfig
	Blockchain BlockchainConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret string
	Expiry int // hours
}

type BlockchainConfig struct {
	NodeAddr           string
	ContractName       string
	ContractAccount    string
	PrivateKeyPath     string
	PrivateKeyPassword string
}

func Load() *Config {
	privateKeyPath := os.Getenv("XUPER_PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		log.Fatalf("environment variable XUPER_PRIVATE_KEY_PATH is required, please set it before starting the server")
	}
	privateKeyPassword := os.Getenv("XUPER_KEY_PASSWORD")
	if privateKeyPassword == "" {
		log.Fatalf("environment variable XUPER_KEY_PASSWORD is required, please set it before starting the server")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("environment variable JWT_SECRET is required, please set it before starting the server")
	}

	bcCfg := BlockchainConfig{
		NodeAddr:           getEnv("BLOCKCHAIN_NODE_ADDR", "39.156.69.83:37100"),
		ContractName:       getEnv("BLOCKCHAIN_CONTRACT_NAME", "hzy_trace"),
		ContractAccount:    getEnv("XUPER_CONTRACT_ACCOUNT", "XC4103761871843472@xuper"),
		PrivateKeyPath:     privateKeyPath,
		PrivateKeyPassword: privateKeyPassword,
	}

	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "123456"),
			Name:     getEnv("DB_NAME", "hzy_trace"),
		},
		JWT: JWTConfig{
			Secret: jwtSecret,
			Expiry: 24,
		},
		Blockchain: bcCfg,
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr string
	AIProvider string
	OpenAIKey  string
	GroqAPIKey string
}

func Load() *Config {
	envPath := ".env"

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		workDir, _ := os.Getwd()
		envPath = filepath.Join(workDir, ".env")
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			envPath = filepath.Join(workDir, "..", ".env")
		}
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("No se encontró archivo .env en %s: %v, usando variables de entorno del sistema", envPath, err)
	} else {
		log.Printf("Archivo .env cargado correctamente desde: %s", envPath)
	}

	provider := getEnv("AI_PROVIDER", "groq")

	log.Printf("Configuración cargada - Provider: %s, ServerAddr: %s",
		provider, getEnv("SERVER_ADDR", ":8080"))

	return &Config{
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		AIProvider: provider,
		OpenAIKey:  getEnv("OPENAI_API_KEY", ""),
		GroqAPIKey: getEnv("GROQ_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

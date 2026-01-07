package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test-pulpoline-api/internal/client/ai"
	"test-pulpoline-api/internal/config"
	"test-pulpoline-api/internal/handler"
	"test-pulpoline-api/internal/queue"
	"test-pulpoline-api/internal/service"
)

func main() {
	cfg := config.Load()

	requestQueue := queue.NewRequestQueue(10)
	defer requestQueue.Close()

	var aiClient ai.AIClient

	switch cfg.AIProvider {
	case "groq":
		log.Println("Usando Groq API (gratuita)")
		aiClient = ai.NewGroqClient(cfg.GroqAPIKey)
	case "openai":
		log.Println("Usando OpenAI API")
		aiClient = ai.NewClient(cfg.OpenAIKey)
	default:
		log.Printf("Proveedor '%s' no reconocido, usando Groq por defecto", cfg.AIProvider)
		aiClient = ai.NewGroqClient(cfg.GroqAPIKey)
	}

	svc := service.NewService(aiClient)

	httpHandler := handler.NewHandler(requestQueue, svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/process", httpHandler.ProcessText)
	mux.HandleFunc("/health", httpHandler.HealthCheck)

	server := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Servidor iniciado en %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Cerrando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error al cerrar el servidor: %v", err)
	}

	log.Println("Servidor cerrado correctamente")
}

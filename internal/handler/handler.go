package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"test-pulpoline-api/internal/queue"
	"test-pulpoline-api/internal/queue/models"
	"test-pulpoline-api/internal/service"
)

type Handler struct {
	queue   *queue.RequestQueue
	service *service.Service
}

func NewHandler(queue *queue.RequestQueue, svc *service.Service) *Handler {
	return &Handler{
		queue:   queue,
		service: svc,
	}
}

type ProcessTextRequest struct {
	Text string `json:"text"`
}

type ProcessTextResponse struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Response string `json:"response"`
	Status   string `json:"status"`
}

func (h *Handler) ProcessText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al leer el cuerpo: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req ProcessTextRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, fmt.Sprintf("Error al parsear JSON: %v", err), http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		http.Error(w, "El campo 'text' no puede estar vacío", http.StatusBadRequest)
		return
	}

	requestID := uuid.New().String()

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	resultChan := make(chan models.Response, 1)
	errorChan := make(chan error, 1)

	queueRequest := models.Request{
		ID:      requestID,
		Text:    req.Text,
		Context: ctx,
		Result:  resultChan,
		Error:   errorChan,
	}

	if err := h.queue.Enqueue(queueRequest); err != nil {
		log.Printf("No se pudo encolar la solicitud %s: %v. Procesando directamente.", requestID, err)
		go h.service.ProcessRequest(ctx, requestID, req.Text, resultChan, errorChan)
	} else {
		go h.service.ProcessRequest(ctx, requestID, req.Text, resultChan, errorChan)
	}

	select {
	case resp := <-resultChan:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ProcessTextResponse{
			ID:       resp.ID,
			Text:     resp.Text,
			Response: resp.Response,
			Status:   "success",
		})
	case err := <-errorChan:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     requestID,
			"error":  err.Error(),
			"status": "error",
		})
	case <-ctx.Done():
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusRequestTimeout)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     requestID,
			"error":  "Timeout: la solicitud tardó demasiado",
			"status": "timeout",
		})
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "test-pulpoline-api",
	})
}

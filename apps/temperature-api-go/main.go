package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"temperature-api-go/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// The Go smart_home service uses DATABASE_URL; pgx understands the
	// postgres:// URI form as well as the key/value form.
	connString := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smarthome")

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to create database pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to database successfully")

	repo := repository.NewSensorRepository(pool)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", rootHandler)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /temperature", temperatureByLocationHandler(repo))
	mux.HandleFunc("GET /temperature/{sensorID}", temperatureByIDHandler(repo))

	srv := &http.Server{
		Addr:              getEnv("PORT", ":8080"),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited properly")
}

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"service": "temperature-api-go", "status": "ok"})
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// temperatureByLocationHandler serves GET /temperature?location=…
// It is invoked by services.TemperatureService.GetTemperature(location).
func temperatureByLocationHandler(repo *repository.SensorRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		location := strings.TrimSpace(r.URL.Query().Get("location"))
		if location == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "The 'location' query parameter is required.",
			})
			return
		}

		temperature, err := repo.GetByLocation(r.Context(), location)
		if err != nil {
			log.Printf("error querying temperature by location %q: %v", location, err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Failed to query temperature data.",
			})
			return
		}
		if temperature == nil {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": fmt.Sprintf("No temperature sensor found for location '%s'.", location),
			})
			return
		}

		writeJSON(w, http.StatusOK, temperature)
	}
}

// temperatureByIDHandler serves GET /temperature/{sensorID}.
// It is invoked by services.TemperatureService.GetTemperatureByID(sensorID).
func temperatureByIDHandler(repo *repository.SensorRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sensorID := r.PathValue("sensorID")

		temperature, err := repo.GetByID(r.Context(), sensorID)
		if err != nil {
			log.Printf("error querying temperature by id %q: %v", sensorID, err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Failed to query temperature data.",
			})
			return
		}
		if temperature == nil {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": fmt.Sprintf("No temperature sensor found with id '%s'.", sensorID),
			})
			return
		}

		writeJSON(w, http.StatusOK, temperature)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

// getEnv returns the value of an environment variable or a fallback default.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

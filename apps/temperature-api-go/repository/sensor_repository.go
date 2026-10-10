package repository

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"

	"temperature-api-go/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SensorRepository reads temperature sensor data from the shared PostgreSQL
// database using raw SQL.
type SensorRepository struct {
	pool *pgxpool.Pool
}

// NewSensorRepository creates a repository backed by the given connection pool.
func NewSensorRepository(pool *pgxpool.Pool) *SensorRepository {
	return &SensorRepository{pool: pool}
}

// GetByLocation returns the most recently updated temperature sensor for a
// location, or nil if none exists.
func (r *SensorRepository) GetByLocation(ctx context.Context, location string) (*models.TemperatureResponse, error) {
	const query = `
		SELECT id, name, type, location, COALESCE(value, 0), COALESCE(unit, ''), status, last_updated
		FROM sensors
		WHERE type = 'temperature' AND lower(location) = lower($1)
		ORDER BY last_updated DESC
		LIMIT 1
	`

	return r.scanOne(ctx, query, location)
}

// GetByID returns a temperature sensor by its numeric id, or nil if the id is
// not an integer or no matching sensor exists.
func (r *SensorRepository) GetByID(ctx context.Context, sensorID string) (*models.TemperatureResponse, error) {
	id, err := strconv.Atoi(sensorID)
	if err != nil {
		return nil, nil
	}

	const query = `
		SELECT id, name, type, location, COALESCE(value, 0), COALESCE(unit, ''), status, last_updated
		FROM sensors
		WHERE id = $1
		LIMIT 1
	`

	return r.scanOne(ctx, query, id)
}

func (r *SensorRepository) scanOne(ctx context.Context, query string, args ...any) (*models.TemperatureResponse, error) {
	var (
		id          int
		name        string
		sensorType  string
		location    string
		value       float64
		unit        string
		status      string
		lastUpdated time.Time
	)

	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&id, &name, &sensorType, &location, &value, &unit, &status, &lastUpdated); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &models.TemperatureResponse{
		// Random temperature between 15.0 and 30.0, rounded to one decimal place.
		Value:       math.Round((15.0+rand.Float64()*15.0)*10) / 10,
		Unit:        unit,
		Timestamp:   lastUpdated,
		Location:    location,
		Status:      status,
		SensorID:    strconv.Itoa(id),
		SensorType:  sensorType,
		Description: name,
	}, nil
}

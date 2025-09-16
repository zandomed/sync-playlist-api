package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User representa un usuario en el sistema
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}

// ServiceAuth almacena tokens OAuth para servicios de música
type ServiceAuth struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	ServiceName  string    `json:"service_name" db:"service_name"`
	AccessToken  string    `json:"-" db:"access_token"`
	RefreshToken string    `json:"-" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	Scope        string    `json:"scope" db:"scope"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Playlist representa una playlist de cualquier servicio
type Playlist struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	ServiceName  string    `json:"service_name" db:"service_name"`
	ExternalID   string    `json:"external_id" db:"external_id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Public       bool      `json:"public" db:"public"`
	TrackCount   int       `json:"track_count" db:"track_count"`
	OwnerID      string    `json:"owner_id" db:"owner_id"`
	OwnerName    string    `json:"owner_name" db:"owner_name"`
	ImageURL     string    `json:"image_url" db:"image_url"`
	ExternalURL  string    `json:"external_url" db:"external_url"`
	LastSyncedAt time.Time `json:"last_synced_at" db:"last_synced_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Track representa una canción dentro de una playlist
type Track struct {
	ID         uuid.UUID `json:"id" db:"id"`
	PlaylistID uuid.UUID `json:"playlist_id" db:"playlist_id"`
	ExternalID string    `json:"external_id" db:"external_id"`
	Name       string    `json:"name" db:"name"`
	Artist     string    `json:"artist" db:"artist"`
	Album      string    `json:"album" db:"album"`
	Duration   int       `json:"duration" db:"duration"`
	ISRC       string    `json:"isrc" db:"isrc"`
	Position   int       `json:"position" db:"position"`
	ImageURL   string    `json:"image_url" db:"image_url"`
	PreviewURL string    `json:"preview_url" db:"preview_url"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Migration representa un proceso de migración de playlist
type Migration struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	UserID           uuid.UUID       `json:"user_id" db:"user_id"`
	SourcePlaylistID uuid.UUID       `json:"source_playlist_id" db:"source_playlist_id"`
	TargetService    string          `json:"target_service" db:"target_service"`
	TargetPlaylistID *uuid.UUID      `json:"target_playlist_id,omitempty" db:"target_playlist_id"`
	Status           MigrationStatus `json:"status" db:"status"`
	TotalTracks      int             `json:"total_tracks" db:"total_tracks"`
	ProcessedTracks  int             `json:"processed_tracks" db:"processed_tracks"`
	SuccessfulTracks int             `json:"successful_tracks" db:"successful_tracks"`
	FailedTracks     int             `json:"failed_tracks" db:"failed_tracks"`
	ErrorMessage     string          `json:"error_message,omitempty" db:"error_message"`
	StartedAt        *time.Time      `json:"started_at,omitempty" db:"started_at"`
	CompletedAt      *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

func (m *Migration) CalculateProgress() float64 {
	if m.TotalTracks == 0 {
		return 0
	}
	return float64(m.ProcessedTracks) / float64(m.TotalTracks) * 100
}

func (m *Migration) IsCompleted() bool {
	return m.Status == MigrationStatusCompleted ||
		m.Status == MigrationStatusFailed ||
		m.Status == MigrationStatusCancelled
}

// MigrationTrack representa el estado de migración de una canción específica
type MigrationTrack struct {
	ID               uuid.UUID            `json:"id" db:"id"`
	MigrationID      uuid.UUID            `json:"migration_id" db:"migration_id"`
	SourceTrackID    uuid.UUID            `json:"source_track_id" db:"source_track_id"`
	TargetExternalID string               `json:"target_external_id,omitempty" db:"target_external_id"`
	Status           MigrationTrackStatus `json:"status" db:"status"`
	ErrorMessage     string               `json:"error_message,omitempty" db:"error_message"`
	AttemptCount     int                  `json:"attempt_count" db:"attempt_count"`
	ProcessedAt      *time.Time           `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt        time.Time            `json:"created_at" db:"created_at"`
}

// Enums con implementación para SQLX
type MigrationStatus string

const (
	MigrationStatusPending   MigrationStatus = "pending"
	MigrationStatusRunning   MigrationStatus = "running"
	MigrationStatusCompleted MigrationStatus = "completed"
	MigrationStatusFailed    MigrationStatus = "failed"
	MigrationStatusCancelled MigrationStatus = "cancelled"
)

// Scan implementa driver.Valuer
func (ms *MigrationStatus) Scan(value interface{}) error {
	if value == nil {
		*ms = MigrationStatusPending
		return nil
	}
	if str, ok := value.(string); ok {
		*ms = MigrationStatus(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into MigrationStatus", value)
}

// Value implementa driver.Valuer
func (ms MigrationStatus) Value() (driver.Value, error) {
	return string(ms), nil
}

type MigrationTrackStatus string

const (
	MigrationTrackStatusPending    MigrationTrackStatus = "pending"
	MigrationTrackStatusProcessing MigrationTrackStatus = "processing"
	MigrationTrackStatusSuccess    MigrationTrackStatus = "success"
	MigrationTrackStatusFailed     MigrationTrackStatus = "failed"
	MigrationTrackStatusNotFound   MigrationTrackStatus = "not_found"
)

// Scan implementa driver.Valuer
func (mts *MigrationTrackStatus) Scan(value interface{}) error {
	if value == nil {
		*mts = MigrationTrackStatusPending
		return nil
	}
	if str, ok := value.(string); ok {
		*mts = MigrationTrackStatus(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into MigrationTrackStatus", value)
}

// Value implementa driver.Valuer
func (mts MigrationTrackStatus) Value() (driver.Value, error) {
	return string(mts), nil
}

// DTOs para queries complejas
type PlaylistWithUser struct {
	Playlist
	UserEmail string `json:"user_email" db:"user_email"`
	UserName  string `json:"user_name" db:"user_name"`
}

type MigrationProgress struct {
	Migration
	SourcePlaylistName string  `json:"source_playlist_name" db:"source_playlist_name"`
	Progress           float64 `json:"progress"`
}

// Funciones para generar UUIDs
func NewUUID() uuid.UUID {
	return uuid.New()
}

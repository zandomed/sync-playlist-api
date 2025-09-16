package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/models"
	"github.com/zandomed/sync-playlist-api/internal/repository"
)

// Services contiene todos los servicios
type Services struct {
	User      UserService
	Auth      AuthService
	Playlist  PlaylistService
	Migration MigrationService
}

// New crea una nueva instancia de servicios
func New(repos *repository.Repositories, cfg *config.Config) *Services {
	return &Services{
		User:      NewUserService(repos.User),
		Auth:      NewAuthService(repos.Auth, cfg),
		Playlist:  NewPlaylistService(repos.Playlist, repos.Track),
		Migration: NewMigrationService(repos.Migration, repos.Playlist, repos.Track, cfg),
	}
}

// UserService interfaz para lógica de negocio de usuarios
type UserService interface {
	CreateUser(email, name string) (*models.User, error)
	GetUser(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(email, name string) (*models.User, error) {
	// Verificar si el usuario ya existe
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil {
		return existingUser, nil // Usuario ya existe
	}

	user := &models.User{
		Email: email,
		Name:  name,
	}

	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUser(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *userService) UpdateUser(user *models.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("invalid user data: %w", err)
	}
	return s.userRepo.Update(user)
}

// AuthService interfaz para lógica de negocio de autenticación
type AuthService interface {
	SaveAuth(userID uuid.UUID, service, accessToken, refreshToken string, expiresAt time.Time) error
	GetAuth(userID uuid.UUID, service string) (*models.ServiceAuth, error)
	RefreshToken(userID uuid.UUID, service string) (*models.ServiceAuth, error)
	IsAuthenticated(userID uuid.UUID, service string) bool
}

type authService struct {
	authRepo repository.AuthRepository
	config   *config.Config
}

func NewAuthService(authRepo repository.AuthRepository, cfg *config.Config) AuthService {
	return &authService{
		authRepo: authRepo,
		config:   cfg,
	}
}

func (s *authService) SaveAuth(userID uuid.UUID, service, accessToken, refreshToken string, expiresAt time.Time) error {
	auth := &models.ServiceAuth{
		UserID:       userID,
		ServiceName:  service,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	return s.authRepo.CreateOrUpdateAuth(auth)
}

func (s *authService) GetAuth(userID uuid.UUID, service string) (*models.ServiceAuth, error) {
	return s.authRepo.GetByUserAndService(userID, service)
}

func (s *authService) RefreshToken(userID uuid.UUID, service string) (*models.ServiceAuth, error) {
	// TODO: Implementar lógica de refresh según el servicio
	// Por ahora retornamos el auth existente
	return s.authRepo.GetByUserAndService(userID, service)
}

func (s *authService) IsAuthenticated(userID uuid.UUID, service string) bool {
	auth, err := s.authRepo.GetByUserAndService(userID, service)
	if err != nil {
		return false
	}

	// Verificar si el token no ha expirado
	return auth.ExpiresAt.After(time.Now())
}

// PlaylistService interfaz para lógica de negocio de playlists
type PlaylistService interface {
	GetUserPlaylists(userID uuid.UUID, service string, page, limit int) ([]models.Playlist, error)
	GetPlaylist(id uuid.UUID) (*models.Playlist, error)
	SyncPlaylistFromService(ctx context.Context, userID uuid.UUID, service, externalID string) (*models.Playlist, error)
	GetPlaylistTracks(playlistID uuid.UUID) ([]models.Track, error)
}

type playlistService struct {
	playlistRepo repository.PlaylistRepository
	trackRepo    repository.TrackRepository
}

func NewPlaylistService(playlistRepo repository.PlaylistRepository, trackRepo repository.TrackRepository) PlaylistService {
	return &playlistService{
		playlistRepo: playlistRepo,
		trackRepo:    trackRepo,
	}
}

func (s *playlistService) GetUserPlaylists(userID uuid.UUID, service string, page, limit int) ([]models.Playlist, error) {
	offset := (page - 1) * limit
	if service == "" {
		return s.playlistRepo.GetByUserID(userID, limit, offset)
	}
	return s.playlistRepo.GetByUserAndService(userID, service)
}

func (s *playlistService) GetPlaylist(id uuid.UUID) (*models.Playlist, error) {
	return s.playlistRepo.GetByID(id)
}

func (s *playlistService) SyncPlaylistFromService(ctx context.Context, userID uuid.UUID, service, externalID string) (*models.Playlist, error) {
	// TODO: Implementar sincronización real con el servicio externo
	// Por ahora, esto es un placeholder

	playlist := &models.Playlist{
		UserID:      userID,
		ServiceName: service,
		ExternalID:  externalID,
		Name:        "Placeholder Playlist",
		Description: "Synced from " + service,
		TrackCount:  0,
	}

	if err := s.playlistRepo.Create(playlist); err != nil {
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	return playlist, nil
}

func (s *playlistService) GetPlaylistTracks(playlistID uuid.UUID) ([]models.Track, error) {
	return s.trackRepo.GetByPlaylistID(playlistID)
}

// MigrationService interfaz para lógica de negocio de migraciones
type MigrationService interface {
	StartMigration(userID, sourcePlaylistID uuid.UUID, targetService string) (*models.Migration, error)
	GetMigration(id uuid.UUID) (*models.Migration, error)
	GetUserMigrations(userID uuid.UUID) ([]models.MigrationProgress, error)
	CancelMigration(id uuid.UUID) error
	ProcessMigration(ctx context.Context, migrationID uuid.UUID) error
}

type migrationService struct {
	migrationRepo repository.MigrationRepository
	playlistRepo  repository.PlaylistRepository
	trackRepo     repository.TrackRepository
	config        *config.Config
}

func NewMigrationService(
	migrationRepo repository.MigrationRepository,
	playlistRepo repository.PlaylistRepository,
	trackRepo repository.TrackRepository,
	cfg *config.Config,
) MigrationService {
	return &migrationService{
		migrationRepo: migrationRepo,
		playlistRepo:  playlistRepo,
		trackRepo:     trackRepo,
		config:        cfg,
	}
}

func (s *migrationService) StartMigration(userID, sourcePlaylistID uuid.UUID, targetService string) (*models.Migration, error) {
	// Verificar que la playlist existe y pertenece al usuario
	playlist, err := s.playlistRepo.GetByID(sourcePlaylistID)
	if err != nil {
		return nil, fmt.Errorf("playlist not found: %w", err)
	}

	if playlist.UserID != userID {
		return nil, fmt.Errorf("playlist does not belong to user")
	}

	// Obtener tracks de la playlist
	tracks, err := s.trackRepo.GetByPlaylistID(sourcePlaylistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist tracks: %w", err)
	}

	// Crear migración
	migration := &models.Migration{
		UserID:           userID,
		SourcePlaylistID: sourcePlaylistID,
		TargetService:    targetService,
		TotalTracks:      len(tracks),
		ProcessedTracks:  0,
		SuccessfulTracks: 0,
		FailedTracks:     0,
	}

	if err := s.migrationRepo.Create(migration); err != nil {
		return nil, fmt.Errorf("failed to create migration: %w", err)
	}

	// Crear registros de track migration
	trackIDs := make([]uuid.UUID, len(tracks))
	for i, track := range tracks {
		trackIDs[i] = track.ID
	}

	if err := s.migrationRepo.CreateTrackMigrations(migration.ID, trackIDs); err != nil {
		return nil, fmt.Errorf("failed to create track migrations: %w", err)
	}

	// Iniciar procesamiento en background
	go func() {
		ctx := context.Background()
		if err := s.ProcessMigration(ctx, migration.ID); err != nil {
			// Log error - en producción usar logger
			fmt.Printf("Migration processing failed: %v\n", err)
		}
	}()

	return migration, nil
}

func (s *migrationService) GetMigration(id uuid.UUID) (*models.Migration, error) {
	return s.migrationRepo.GetByID(id)
}

func (s *migrationService) GetUserMigrations(userID uuid.UUID) ([]models.MigrationProgress, error) {
	return s.migrationRepo.GetByUserID(userID)
}

func (s *migrationService) CancelMigration(id uuid.UUID) error {
	return s.migrationRepo.UpdateStatus(id, models.MigrationStatusCancelled, "Cancelled by user")
}

func (s *migrationService) ProcessMigration(ctx context.Context, migrationID uuid.UUID) error {
	// Obtener migración
	_, err := s.migrationRepo.GetByID(migrationID)
	if err != nil {
		return fmt.Errorf("failed to get migration: %w", err)
	}

	// Actualizar status a running
	if err := s.migrationRepo.UpdateStatus(migrationID, models.MigrationStatusRunning, ""); err != nil {
		return fmt.Errorf("failed to update migration status: %w", err)
	}

	// Obtener tracks a migrar
	trackMigrations, err := s.migrationRepo.GetTrackMigrations(migrationID)
	if err != nil {
		return fmt.Errorf("failed to get track migrations: %w", err)
	}

	successful := 0
	failed := 0

	// Procesar cada track
	for i, trackMigration := range trackMigrations {
		select {
		case <-ctx.Done():
			// Contexto cancelado
			s.migrationRepo.UpdateStatus(migrationID, models.MigrationStatusCancelled, "Context cancelled")
			return ctx.Err()
		default:
			// Simular procesamiento de track
			// TODO: Implementar lógica real de migración aquí
			success := s.processTrackMigration(ctx, trackMigration)

			if success {
				successful++
				s.migrationRepo.UpdateTrackStatus(trackMigration.ID, models.MigrationTrackStatusSuccess, "target_id_placeholder", "")
			} else {
				failed++
				s.migrationRepo.UpdateTrackStatus(trackMigration.ID, models.MigrationTrackStatusFailed, "", "Failed to migrate track")
			}

			// Actualizar progreso
			processed := i + 1
			s.migrationRepo.UpdateProgress(migrationID, processed, successful, failed)
		}
	}

	// Finalizar migración
	finalStatus := models.MigrationStatusCompleted
	errorMsg := ""
	if failed == len(trackMigrations) {
		finalStatus = models.MigrationStatusFailed
		errorMsg = "All tracks failed to migrate"
	}

	return s.migrationRepo.UpdateStatus(migrationID, finalStatus, errorMsg)
}

// processTrackMigration simula el procesamiento de un track individual
// TODO: Reemplazar con lógica real de migración
func (s *migrationService) processTrackMigration(ctx context.Context, trackMigration models.MigrationTrack) bool {
	// Simular trabajo con sleep
	time.Sleep(100 * time.Millisecond)

	// Simular 80% de éxito
	return time.Now().UnixNano()%10 < 8
}

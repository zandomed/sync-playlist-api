package repository

import (
	"github.com/zandomed/sync-playlist-api/pkg/database"
)

// Repositories contiene todos los repositorios
type Repositories struct {
	// User      UserRepository
	Auth AuthRepository
	// Playlist  PlaylistRepository
	// Track     TrackRepository
	// Migration MigrationRepository
}

// New crea una nueva instancia de repositorios
func New(db *database.DB) *Repositories {
	return &Repositories{
		// User:      NewUserRepository(db),
		Auth: NewAuthRepository(db),
		// Playlist:  NewPlaylistRepository(db),
		// Track:     NewTrackRepository(db),
		// Migration: NewMigrationRepository(db),
	}
}

// // UserRepository interfaz para operaciones de usuario
// type UserRepository interface {
// 	Create(user *models.User) error
// 	GetByID(id uuid.UUID) (*models.User, error)
// 	GetByEmail(email string) (*models.User, error)
// 	Update(user *models.User) error
// 	Delete(id uuid.UUID) error
// }

// type userRepository struct {
// 	db *database.DB
// }

// func NewUserRepository(db *database.DB) UserRepository {
// 	return &userRepository{db: db}
// }

// func (r *userRepository) Create(user *models.User) error {
// 	user.ID = models.NewUUID()
// 	user.CreatedAt = time.Now()
// 	user.UpdatedAt = time.Now()

// 	query := `
// 		INSERT INTO users (id, email, name, created_at, updated_at)
// 		VALUES (:id, :email, :name, :created_at, :updated_at)`

// 	_, err := r.db.NamedExec(query, user)
// 	return err
// }

// func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
// 	var user models.User
// 	query := `SELECT * FROM users WHERE id = $1`
// 	err := r.db.Get(&user, query, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func (r *userRepository) GetByEmail(email string) (*models.User, error) {
// 	var user models.User
// 	query := `SELECT * FROM users WHERE email = $1`
// 	err := r.db.Get(&user, query, email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func (r *userRepository) Update(user *models.User) error {
// 	user.UpdatedAt = time.Now()
// 	query := `
// 		UPDATE users
// 		SET name = :name, updated_at = :updated_at
// 		WHERE id = :id`

// 	_, err := r.db.NamedExec(query, user)
// 	return err
// }

// func (r *userRepository) Delete(id uuid.UUID) error {
// 	query := `DELETE FROM users WHERE id = $1`
// 	_, err := r.db.Exec(query, id)
// 	return err
// }

// // AuthRepository interfaz para operaciones de autenticación
// // type AuthRepository interface {
// // 	CreateOrUpdateAuth(auth *models.ServiceAuth) error
// // 	GetByUserAndService(userID uuid.UUID, service string) (*models.ServiceAuth, error)
// // 	DeleteExpiredTokens() error
// // 	GetUserServices(userID uuid.UUID) ([]models.ServiceAuth, error)
// // }

// // type authRepository struct {
// // 	db *database.DB
// // }

// // func NewAuthRepository(db *database.DB) AuthRepository {
// // 	return &authRepository{db: db}
// // }

// // func (r *authRepository) CreateOrUpdateAuth(auth *models.ServiceAuth) error {
// // 	auth.UpdatedAt = time.Now()

// // 	query := `
// // 		INSERT INTO service_auths (id, user_id, service_name, access_token, refresh_token, expires_at, scope, created_at, updated_at)
// // 		VALUES (:id, :user_id, :service_name, :access_token, :refresh_token, :expires_at, :scope, :created_at, :updated_at)
// // 		ON CONFLICT (user_id, service_name)
// // 		DO UPDATE SET
// // 			access_token = EXCLUDED.access_token,
// // 			refresh_token = EXCLUDED.refresh_token,
// // 			expires_at = EXCLUDED.expires_at,
// // 			scope = EXCLUDED.scope,
// // 			updated_at = EXCLUDED.updated_at`

// // 	if auth.ID == uuid.Nil {
// // 		auth.ID = models.NewUUID()
// // 		auth.CreatedAt = time.Now()
// // 	}

// // 	_, err := r.db.NamedExec(query, auth)
// // 	return err
// // }

// // func (r *authRepository) GetByUserAndService(userID uuid.UUID, service string) (*models.ServiceAuth, error) {
// // 	var auth models.ServiceAuth
// // 	query := `SELECT * FROM service_auths WHERE user_id = $1 AND service_name = $2`
// // 	err := r.db.Get(&auth, query, userID, service)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return &auth, nil
// // }

// // func (r *authRepository) DeleteExpiredTokens() error {
// // 	query := `DELETE FROM service_auths WHERE expires_at < CURRENT_TIMESTAMP`
// // 	_, err := r.db.Exec(query)
// // 	return err
// // }

// // func (r *authRepository) GetUserServices(userID uuid.UUID) ([]models.ServiceAuth, error) {
// // 	var auths []models.ServiceAuth
// // 	query := `SELECT * FROM service_auths WHERE user_id = $1 ORDER BY service_name`
// // 	err := r.db.Select(&auths, query, userID)
// // 	return auths, err
// // }

// // PlaylistRepository interfaz para operaciones de playlist
// type PlaylistRepository interface {
// 	Create(playlist *models.Playlist) error
// 	GetByID(id uuid.UUID) (*models.Playlist, error)
// 	GetByUserID(userID uuid.UUID, limit, offset int) ([]models.Playlist, error)
// 	GetByUserAndService(userID uuid.UUID, service string) ([]models.Playlist, error)
// 	Update(playlist *models.Playlist) error
// 	Delete(id uuid.UUID) error
// 	SyncPlaylist(playlist *models.Playlist) error
// }

// type playlistRepository struct {
// 	db *database.DB
// }

// func NewPlaylistRepository(db *database.DB) PlaylistRepository {
// 	return &playlistRepository{db: db}
// }

// func (r *playlistRepository) Create(playlist *models.Playlist) error {
// 	playlist.ID = models.NewUUID()
// 	playlist.CreatedAt = time.Now()
// 	playlist.UpdatedAt = time.Now()
// 	playlist.LastSyncedAt = time.Now()

// 	query := `
// 		INSERT INTO playlists (id, user_id, service_name, external_id, name, description, public,
// 			track_count, owner_id, owner_name, image_url, external_url, last_synced_at, created_at, updated_at)
// 		VALUES (:id, :user_id, :service_name, :external_id, :name, :description, :public,
// 			:track_count, :owner_id, :owner_name, :image_url, :external_url, :last_synced_at, :created_at, :updated_at)`

// 	_, err := r.db.NamedExec(query, playlist)
// 	return err
// }

// func (r *playlistRepository) GetByID(id uuid.UUID) (*models.Playlist, error) {
// 	var playlist models.Playlist
// 	query := `SELECT * FROM playlists WHERE id = $1`
// 	err := r.db.Get(&playlist, query, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &playlist, nil
// }

// func (r *playlistRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]models.Playlist, error) {
// 	var playlists []models.Playlist
// 	query := `
// 		SELECT * FROM playlists
// 		WHERE user_id = $1
// 		ORDER BY created_at DESC
// 		LIMIT $2 OFFSET $3`

// 	err := r.db.Select(&playlists, query, userID, limit, offset)
// 	return playlists, err
// }

// func (r *playlistRepository) GetByUserAndService(userID uuid.UUID, service string) ([]models.Playlist, error) {
// 	var playlists []models.Playlist
// 	query := `
// 		SELECT * FROM playlists
// 		WHERE user_id = $1 AND service_name = $2
// 		ORDER BY name`

// 	err := r.db.Select(&playlists, query, userID, service)
// 	return playlists, err
// }

// func (r *playlistRepository) Update(playlist *models.Playlist) error {
// 	playlist.UpdatedAt = time.Now()
// 	query := `
// 		UPDATE playlists
// 		SET name = :name, description = :description, public = :public,
// 			track_count = :track_count, image_url = :image_url, updated_at = :updated_at
// 		WHERE id = :id`

// 	_, err := r.db.NamedExec(query, playlist)
// 	return err
// }

// func (r *playlistRepository) Delete(id uuid.UUID) error {
// 	query := `DELETE FROM playlists WHERE id = $1`
// 	_, err := r.db.Exec(query, id)
// 	return err
// }

// func (r *playlistRepository) SyncPlaylist(playlist *models.Playlist) error {
// 	playlist.UpdatedAt = time.Now()
// 	playlist.LastSyncedAt = time.Now()

// 	query := `
// 		UPDATE playlists
// 		SET track_count = :track_count, last_synced_at = :last_synced_at, updated_at = :updated_at
// 		WHERE id = :id`

// 	_, err := r.db.NamedExec(query, playlist)
// 	return err
// }

// // TrackRepository interfaz para operaciones de tracks
// type TrackRepository interface {
// 	CreateBatch(tracks []models.Track) error
// 	GetByPlaylistID(playlistID uuid.UUID) ([]models.Track, error)
// 	SearchTrack(name, artist, isrc string) (*models.Track, error)
// 	DeleteByPlaylistID(playlistID uuid.UUID) error
// }

// type trackRepository struct {
// 	db *database.DB
// }

// func NewTrackRepository(db *database.DB) TrackRepository {
// 	return &trackRepository{db: db}
// }

// func (r *trackRepository) CreateBatch(tracks []models.Track) error {
// 	if len(tracks) == 0 {
// 		return nil
// 	}

// 	// Construcción de query batch optimizada
// 	valueStrings := make([]string, 0, len(tracks))
// 	valueArgs := make([]interface{}, 0, len(tracks)*10)

// 	for i, track := range tracks {
// 		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
// 			i*10+1, i*10+2, i*10+3, i*10+4, i*10+5, i*10+6, i*10+7, i*10+8, i*10+9, i*10+10))

// 		if track.ID == uuid.Nil {
// 			track.ID = models.NewUUID()
// 		}
// 		if track.CreatedAt.IsZero() {
// 			track.CreatedAt = time.Now()
// 		}

// 		valueArgs = append(valueArgs, track.ID, track.PlaylistID, track.ExternalID, track.Name,
// 			track.Artist, track.Album, track.Duration, track.ISRC, track.Position, track.CreatedAt)
// 	}

// 	query := fmt.Sprintf(`
// 		INSERT INTO tracks (id, playlist_id, external_id, name, artist, album, duration, isrc, position, created_at)
// 		VALUES %s`, strings.Join(valueStrings, ","))

// 	_, err := r.db.Exec(query, valueArgs...)
// 	return err
// }

// func (r *trackRepository) GetByPlaylistID(playlistID uuid.UUID) ([]models.Track, error) {
// 	var tracks []models.Track
// 	query := `SELECT * FROM tracks WHERE playlist_id = $1 ORDER BY position`
// 	err := r.db.Select(&tracks, query, playlistID)
// 	return tracks, err
// }

// func (r *trackRepository) SearchTrack(name, artist, isrc string) (*models.Track, error) {
// 	var track models.Track

// 	// Búsqueda optimizada: primero por ISRC, luego por full-text search
// 	if isrc != "" {
// 		query := `SELECT * FROM tracks WHERE isrc = $1 LIMIT 1`
// 		err := r.db.Get(&track, query, isrc)
// 		if err == nil {
// 			return &track, nil
// 		} else if err != sql.ErrNoRows {
// 			return nil, err
// 		}
// 	}

// 	// Full-text search como fallback
// 	query := `
// 		SELECT * FROM tracks
// 		WHERE to_tsvector('english', name || ' ' || artist) @@ plainto_tsquery($1)
// 		ORDER BY ts_rank(to_tsvector('english', name || ' ' || artist), plainto_tsquery($1)) DESC
// 		LIMIT 1`

// 	searchTerm := fmt.Sprintf("%s %s", name, artist)
// 	err := r.db.Get(&track, query, searchTerm)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &track, nil
// }

// func (r *trackRepository) DeleteByPlaylistID(playlistID uuid.UUID) error {
// 	query := `DELETE FROM tracks WHERE playlist_id = $1`
// 	_, err := r.db.Exec(query, playlistID)
// 	return err
// }

// // MigrationRepository interfaz para operaciones de migración
// type MigrationRepository interface {
// 	Create(migration *models.Migration) error
// 	GetByID(id uuid.UUID) (*models.Migration, error)
// 	GetByUserID(userID uuid.UUID) ([]models.MigrationProgress, error)
// 	UpdateStatus(id uuid.UUID, status models.MigrationStatus, errorMsg string) error
// 	UpdateProgress(id uuid.UUID, processed, successful, failed int) error
// 	CreateTrackMigrations(migrationID uuid.UUID, trackIDs []uuid.UUID) error
// 	GetTrackMigrations(migrationID uuid.UUID) ([]models.MigrationTrack, error)
// 	UpdateTrackStatus(trackMigrationID uuid.UUID, status models.MigrationTrackStatus, targetID, errorMsg string) error
// }

// type migrationRepository struct {
// 	db *database.DB
// }

// func NewMigrationRepository(db *database.DB) MigrationRepository {
// 	return &migrationRepository{db: db}
// }

// func (r *migrationRepository) Create(migration *models.Migration) error {
// 	migration.ID = models.NewUUID()
// 	migration.CreatedAt = time.Now()
// 	migration.UpdatedAt = time.Now()
// 	migration.Status = models.MigrationStatusPending

// 	query := `
// 		INSERT INTO migrations (id, user_id, source_playlist_id, target_service, status,
// 			total_tracks, processed_tracks, successful_tracks, failed_tracks, created_at, updated_at)
// 		VALUES (:id, :user_id, :source_playlist_id, :target_service, :status,
// 			:total_tracks, :processed_tracks, :successful_tracks, :failed_tracks, :created_at, :updated_at)`

// 	_, err := r.db.NamedExec(query, migration)
// 	return err
// }

// func (r *migrationRepository) GetByID(id uuid.UUID) (*models.Migration, error) {
// 	var migration models.Migration
// 	query := `SELECT * FROM migrations WHERE id = $1`
// 	err := r.db.Get(&migration, query, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &migration, nil
// }

// func (r *migrationRepository) GetByUserID(userID uuid.UUID) ([]models.MigrationProgress, error) {
// 	var migrations []models.MigrationProgress
// 	query := `
// 		SELECT m.*, p.name as source_playlist_name
// 		FROM migrations m
// 		JOIN playlists p ON m.source_playlist_id = p.id
// 		WHERE m.user_id = $1
// 		ORDER BY m.created_at DESC`

// 	err := r.db.Select(&migrations, query, userID)

// 	// Calcular progreso para cada migración
// 	for i := range migrations {
// 		migrations[i].Progress = migrations[i].CalculateProgress()
// 	}

// 	return migrations, err
// }

// func (r *migrationRepository) UpdateStatus(id uuid.UUID, status models.MigrationStatus, errorMsg string) error {
// 	now := time.Now()
// 	query := `
// 		UPDATE migrations
// 		SET status = $2, error_message = $3, updated_at = $4`

// 	args := []interface{}{id, status, errorMsg, now}

// 	switch status {
// 	case models.MigrationStatusRunning:
// 		query += `, started_at = $5`
// 		args = append(args, now)
// 	case models.MigrationStatusCompleted, models.MigrationStatusFailed:
// 		query += `, completed_at = $5`
// 		args = append(args, now)
// 	}

// 	query += ` WHERE id = $1`

// 	_, err := r.db.Exec(query, args...)
// 	return err
// }

// func (r *migrationRepository) UpdateProgress(id uuid.UUID, processed, successful, failed int) error {
// 	query := `
// 		UPDATE migrations
// 		SET processed_tracks = $2, successful_tracks = $3, failed_tracks = $4, updated_at = $5
// 		WHERE id = $1`

// 	_, err := r.db.Exec(query, id, processed, successful, failed, time.Now())
// 	return err
// }

// func (r *migrationRepository) CreateTrackMigrations(migrationID uuid.UUID, trackIDs []uuid.UUID) error {
// 	if len(trackIDs) == 0 {
// 		return nil
// 	}

// 	valueStrings := make([]string, 0, len(trackIDs))
// 	valueArgs := make([]interface{}, 0, len(trackIDs)*3)

// 	for i, trackID := range trackIDs {
// 		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
// 		valueArgs = append(valueArgs, models.NewUUID(), migrationID, trackID)
// 	}

// 	query := fmt.Sprintf(`
// 		INSERT INTO migration_tracks (id, migration_id, source_track_id)
// 		VALUES %s`, strings.Join(valueStrings, ","))

// 	_, err := r.db.Exec(query, valueArgs...)
// 	return err
// }

// func (r *migrationRepository) GetTrackMigrations(migrationID uuid.UUID) ([]models.MigrationTrack, error) {
// 	var tracks []models.MigrationTrack
// 	query := `SELECT * FROM migration_tracks WHERE migration_id = $1 ORDER BY created_at`
// 	err := r.db.Select(&tracks, query, migrationID)
// 	return tracks, err
// }

// func (r *migrationRepository) UpdateTrackStatus(trackMigrationID uuid.UUID, status models.MigrationTrackStatus, targetID, errorMsg string) error {
// 	query := `
// 		UPDATE migration_tracks
// 		SET status = $2, target_external_id = $3, error_message = $4,
// 			processed_at = $5, attempt_count = attempt_count + 1
// 		WHERE id = $1`

// 	_, err := r.db.Exec(query, trackMigrationID, status, targetID, errorMsg, time.Now())
// 	return err
// }

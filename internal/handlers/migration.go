package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/zandomed/sync-playlist-api/internal/models"

	"github.com/zandomed/sync-playlist-api/internal/services"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

// MigrationHandler maneja operaciones de migración
type MigrationHandler struct {
	migrationService services.MigrationService
	logger           *logger.Logger
	upgrader         websocket.Upgrader
}

func NewMigrationHandler(migrationService services.MigrationService, logger *logger.Logger) *MigrationHandler {
	return &MigrationHandler{
		migrationService: migrationService,
		logger:           logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// En producción, verificar orígenes permitidos
				return true
			},
		},
	}
}

// StartMigrationRequest estructura para iniciar migración
type StartMigrationRequest struct {
	SourcePlaylistID string `json:"source_playlist_id" validate:"required,uuid"`
	TargetService    string `json:"target_service" validate:"required,service"`
}

// StartMigration inicia una nueva migración
func (h *MigrationHandler) StartMigration(c echo.Context) error {
	claims, err := getUserFromContext(c)
	if err != nil {
		return sendError(c, http.StatusUnauthorized, err, "Invalid token")
	}

	req := new(StartMigrationRequest)
	if err := c.Bind(req); err != nil {
		return sendError(c, http.StatusBadRequest, err, "Invalid request body")
	}

	if err := c.Validate(req); err != nil {
		return sendError(c, http.StatusBadRequest, err, "Validation failed")
	}

	sourcePlaylistID, err := uuid.Parse(req.SourcePlaylistID)
	if err != nil {
		return sendError(c, http.StatusBadRequest, err, "Invalid source playlist ID")
	}

	migration, err := h.migrationService.StartMigration(claims.UserID, sourcePlaylistID, req.TargetService)
	if err != nil {
		return sendError(c, http.StatusInternalServerError, err, "Failed to start migration")
	}

	return sendSuccess(c, http.StatusCreated, migration, "Migration started successfully")
}

// GetMigrationStatus obtiene el estado de una migración
func (h *MigrationHandler) GetMigrationStatus(c echo.Context) error {
	claims, err := getUserFromContext(c)
	if err != nil {
		return sendError(c, http.StatusUnauthorized, err, "Invalid token")
	}

	migrationIDStr := c.Param("id")
	migrationID, err := uuid.Parse(migrationIDStr)
	if err != nil {
		return sendError(c, http.StatusBadRequest, err, "Invalid migration ID")
	}

	migration, err := h.migrationService.GetMigration(migrationID)
	if err != nil {
		return sendError(c, http.StatusNotFound, err, "Migration not found")
	}

	// Verificar que la migración pertenece al usuario
	if migration.UserID != claims.UserID {
		return sendError(c, http.StatusForbidden, nil, "Access denied")
	}

	return sendSuccess(c, http.StatusOK, migration, "")
}

// GetMigrationProgress obtiene el progreso detallado de una migración
func (h *MigrationHandler) GetMigrationProgress(c echo.Context) error {
	claims, err := getUserFromContext(c)
	if err != nil {
		return sendError(c, http.StatusUnauthorized, err, "Invalid token")
	}

	migrationIDStr := c.Param("id")
	migrationID, err := uuid.Parse(migrationIDStr)
	if err != nil {
		return sendError(c, http.StatusBadRequest, err, "Invalid migration ID")
	}

	migration, err := h.migrationService.GetMigration(migrationID)
	if err != nil {
		return sendError(c, http.StatusNotFound, err, "Migration not found")
	}

	// Verificar que la migración pertenece al usuario
	if migration.UserID != claims.UserID {
		return sendError(c, http.StatusForbidden, nil, "Access denied")
	}

	progress := map[string]interface{}{
		"migration":           migration,
		"progress_percent":    migration.CalculateProgress(),
		"is_completed":        migration.IsCompleted(),
		"estimated_remaining": h.calculateEstimatedTime(migration),
	}

	return sendSuccess(c, http.StatusOK, progress, "")
}

// CancelMigration cancela una migración en curso
func (h *MigrationHandler) CancelMigration(c echo.Context) error {
	claims, err := getUserFromContext(c)
	if err != nil {
		return sendError(c, http.StatusUnauthorized, err, "Invalid token")
	}

	migrationIDStr := c.Param("id")
	migrationID, err := uuid.Parse(migrationIDStr)
	if err != nil {
		return sendError(c, http.StatusBadRequest, err, "Invalid migration ID")
	}

	migration, err := h.migrationService.GetMigration(migrationID)
	if err != nil {
		return sendError(c, http.StatusNotFound, err, "Migration not found")
	}

	// Verificar que la migración pertenece al usuario
	if migration.UserID != claims.UserID {
		return sendError(c, http.StatusForbidden, nil, "Access denied")
	}

	// Solo se puede cancelar si está pending o running
	if migration.IsCompleted() {
		return sendError(c, http.StatusBadRequest, nil, "Cannot cancel completed migration")
	}

	if err := h.migrationService.CancelMigration(migrationID); err != nil {
		return sendError(c, http.StatusInternalServerError, err, "Failed to cancel migration")
	}

	return sendSuccess(c, http.StatusOK, nil, "Migration cancelled successfully")
}

// WebSocketHandler maneja conexiones WebSocket para progreso en tiempo real
func (h *MigrationHandler) WebSocketHandler(c echo.Context) error {
	migrationIDStr := c.Param("id")
	migrationID, err := uuid.Parse(migrationIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid migration ID")
	}

	// Upgrade la conexión HTTP a WebSocket
	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := ws.Close(); err != nil {
			h.logger.Sugar().Errorf("Error closing websocket: %v", err)
		}
	}()

	// TODO: En una implementación real, aquí necesitarías:
	// 1. Verificar autenticación via token en query params o headers
	// 2. Verificar que el usuario tiene acceso a esta migración
	// 3. Subscribirse a eventos de progreso (usando Redis pub/sub o similar)

	// Por ahora, simulamos envío de progreso cada segundo
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Obtener estado actual de la migración
			migration, err := h.migrationService.GetMigration(migrationID)
			if err != nil {
				h.logger.Sugar().Errorf("Failed to get migration for websocket: %v", err)
				return nil
			}

			// Enviar progreso al cliente
			progress := map[string]interface{}{
				"migration_id":      migration.ID,
				"status":            migration.Status,
				"total_tracks":      migration.TotalTracks,
				"processed_tracks":  migration.ProcessedTracks,
				"successful_tracks": migration.SuccessfulTracks,
				"failed_tracks":     migration.FailedTracks,
				"progress_percent":  migration.CalculateProgress(),
				"is_completed":      migration.IsCompleted(),
				"timestamp":         time.Now().Unix(),
			}

			if err := ws.WriteJSON(progress); err != nil {
				h.logger.Sugar().Errorf("Failed to write websocket message: %v", err)
				return nil
			}

			// Si la migración está completa, cerrar conexión
			if migration.IsCompleted() {
				// Enviar mensaje final
				finalMessage := map[string]interface{}{
					"type":    "completion",
					"status":  migration.Status,
					"message": "Migration completed",
				}
				if err := ws.WriteJSON(finalMessage); err != nil {
					h.logger.Sugar().Errorf("Failed to write websocket message: %v", err)
					return nil
				}
				return nil
			}

		case <-c.Request().Context().Done():
			return nil
		}
	}
}

// GetUserMigrations obtiene todas las migraciones del usuario
func (h *MigrationHandler) GetUserMigrations(c echo.Context) error {
	claims, err := getUserFromContext(c)
	if err != nil {
		return sendError(c, http.StatusUnauthorized, err, "Invalid token")
	}

	migrations, err := h.migrationService.GetUserMigrations(claims.UserID)
	if err != nil {
		return sendError(c, http.StatusInternalServerError, err, "Failed to get migrations")
	}

	return sendSuccess(c, http.StatusOK, migrations, "")
}

// Helper functions
func (h *MigrationHandler) calculateEstimatedTime(migration *models.Migration) string {
	if migration.ProcessedTracks == 0 || migration.StartedAt == nil {
		return "Unknown"
	}

	elapsed := time.Since(*migration.StartedAt)
	avgTimePerTrack := elapsed / time.Duration(migration.ProcessedTracks)
	remainingTracks := migration.TotalTracks - migration.ProcessedTracks

	if remainingTracks <= 0 {
		return "0s"
	}

	estimated := avgTimePerTrack * time.Duration(remainingTracks)
	return estimated.Truncate(time.Second).String()
}

package health

import (
	"time"

	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/domain"
)

type GetStatusUseCase struct {
	config *config.Config
}

func NewGetStatusUseCase(config *config.Config) *GetStatusUseCase {
	return &GetStatusUseCase{
		config: config,
	}
}

func (uc *GetStatusUseCase) Execute() *domain.HealthStatus {
	return &domain.HealthStatus{
		Status:      "healthy",
		Version:     "1.0.0",
		Environment: string(uc.config.Server.Environment),
		Timestamp:   time.Now(),
	}
}

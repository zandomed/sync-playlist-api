package auth

import (
	"context"

	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
)

type VerifyTokenRequest struct {
	Token string
}

type VerifyTokenResponse struct {
	Valid  bool
	UserID string
}

type VerifyTokenUseCase struct {
	verificationRepo repositories.VerificationRepository
}

func NewVerifyTokenUseCase(verificationRepo repositories.VerificationRepository) *VerifyTokenUseCase {
	return &VerifyTokenUseCase{
		verificationRepo: verificationRepo,
	}
}

func (uc *VerifyTokenUseCase) Execute(ctx context.Context, req VerifyTokenRequest) (*VerifyTokenResponse, error) {
	// Find the verification token
	verificationToken, err := uc.verificationRepo.FindByToken(ctx, req.Token)
	if err != nil {
		return &VerifyTokenResponse{
			Valid:  false,
			UserID: "",
		}, nil
	}

	// Validate the token for frontend verification
	if err := verificationToken.ValidateForFrontend(); err != nil {
		return &VerifyTokenResponse{
			Valid:  false,
			UserID: "",
		}, nil
	}

	// Mark the token as used
	if err := verificationToken.MarkAsUsed(); err != nil {
		return &VerifyTokenResponse{
			Valid:  false,
			UserID: "",
		}, nil
	}

	// Update the token in the database
	if err := uc.verificationRepo.Update(ctx, verificationToken); err != nil {
		return &VerifyTokenResponse{
			Valid:  false,
			UserID: "",
		}, nil
	}

	return &VerifyTokenResponse{
		Valid:  true,
		UserID: verificationToken.UserID().String(),
	}, nil
}

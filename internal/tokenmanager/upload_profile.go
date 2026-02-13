package tokenmanager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/netbill/restkit/tokens"
)

func (m *Manager) GenerateUploadProfileMediaToken(
	OwnerAccountID uuid.UUID,
	UploadSessionID uuid.UUID,
) (string, error) {
	tkn, err := tokens.UploadContentClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   OwnerAccountID.String(),
			Issuer:    m.Issuer,
			Audience:  []string{m.Issuer},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(m.ProfileMediaUploadTTL)),
		},
		UploadSessionID: UploadSessionID,
		ResourceID:      OwnerAccountID.String(),
		ResourceType:    ProfileResource,
	}.GenerateJWT(m.UploadSK)
	if err != nil {
		return "", fmt.Errorf("failed to generate upload profile media token, cause: %w", err)
	}

	return tkn, nil
}

func (m *Manager) ParseUploadProfileContentToken(token string) (tokens.UploadContentClaims, error) {
	res, err := tokens.ParseUploadFilesClaims(token, m.UploadSK)
	if err != nil {
		return tokens.UploadContentClaims{}, fmt.Errorf(
			"failed to validate upload profile media token, cause: %w", err,
		)
	}

	if res.ResourceType != ProfileResource {
		return tokens.UploadContentClaims{}, fmt.Errorf("invalid upload token resource type")
	}

	audSuccess := false
	for _, aud := range res.Audience {
		if aud == m.Issuer {
			audSuccess = true
			break
		}
	}
	if !audSuccess {
		return tokens.UploadContentClaims{}, fmt.Errorf("invalid upload token audience")
	}

	_, err = uuid.Parse(res.Subject)
	if err != nil {
		return tokens.UploadContentClaims{}, fmt.Errorf("invalid upload token subject: %w", err)
	}

	_, err = uuid.Parse(res.ResourceID)
	if err != nil {
		return tokens.UploadContentClaims{}, fmt.Errorf("invalid upload token resource ID: %w", err)
	}

	return res, nil
}

package tokenmanager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/netbill/restkit/tokens"
)

type Manager struct {
	uploadSK string

	profileMediaUploadTTL time.Duration
}

const (
	ProfilesActor   = "profiles-svc"
	ProfileResource = "profile"
)

func New(uploadSK string, profileMediaUploadTTL time.Duration) Manager {
	return Manager{
		uploadSK:              uploadSK,
		profileMediaUploadTTL: profileMediaUploadTTL,
	}
}

func (m Manager) NewUploadProfileMediaToken(
	OwnerAccountID uuid.UUID,
	UploadSessionID uuid.UUID,
) (string, error) {
	tkn, err := tokens.UploadContentClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   OwnerAccountID.String(),
			Issuer:    ProfilesActor,
			Audience:  []string{ProfilesActor},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(m.profileMediaUploadTTL)),
		},
		UploadSessionID: UploadSessionID,
		ResourceID:      OwnerAccountID.String(),
		Resource:        ProfileResource,
	}.GenerateJWT(m.uploadSK)

	if err != nil {
		return "", fmt.Errorf("failed to generate upload profile media token, cause: %w", err)
	}

	return tkn, nil
}

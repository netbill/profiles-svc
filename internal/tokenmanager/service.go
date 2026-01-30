package tokenmanager

import (
	"time"

	"github.com/google/uuid"
	"github.com/netbill/restkit/tokens"
)

type Manager struct {
	uploadSK string

	profileMediaUploadTTL time.Duration
}

const (
	ProfilesService = "profiles-svc"
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
	return tokens.NewUploadFileToken(
		tokens.GenerateUploadFilesJwtRequest{
			OwnerAccountID:  OwnerAccountID,
			UploadSessionID: UploadSessionID,
			ResourceID:      OwnerAccountID.String(),
			Resource:        ProfileResource,
			Issuer:          ProfilesService,
			Audience:        []string{ProfilesService},
			Ttl:             m.profileMediaUploadTTL,
		}, m.uploadSK)
}

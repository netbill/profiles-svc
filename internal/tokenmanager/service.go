package tokenmanager

import (
	"time"

	"github.com/google/uuid"
	"github.com/netbill/restkit/tokens"
)

type Manager struct {
	uploadSK string
}

const (
	ProfilesService       = "profiles-svc"
	ProfileResource       = "profile"
	ProfileMediaUploadTTL = 1 * time.Hour
)

func New(issuer, uploadSK string) Manager {
	return Manager{
		uploadSK: uploadSK,
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
			Ttl:             ProfileMediaUploadTTL,
		}, m.uploadSK)
}

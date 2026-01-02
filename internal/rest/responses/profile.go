package responses

import (
	"github.com/netbill/pagi"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/resources"
)

func Profile(m models.Profile) resources.Profile {
	resp := resources.Profile{
		Data: resources.ProfileData{
			Id:   m.AccountID,
			Type: resources.ProfileType,
			Attributes: resources.ProfileAttributes{
				Username:    m.Username,
				Pseudonym:   m.Pseudonym,
				Description: m.Description,
				Avatar:      m.Avatar,
				Official:    m.Official,
				UpdatedAt:   m.UpdatedAt,
				CreatedAt:   m.CreatedAt,
			},
		},
	}

	return resp
}

func ProfileCollection(m pagi.Page[[]models.Profile]) resources.ProfilesCollection {
	resp := resources.ProfilesCollection{
		Data: make([]resources.ProfileData, 0, len(m.Data)),
	}

	for _, el := range m.Data {
		p := Profile(el).Data

		resp.Data = append(resp.Data, p)
	}

	return resp
}

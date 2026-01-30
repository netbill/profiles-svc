package responses

import (
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/resources"
	"github.com/netbill/restkit/pagi"
)

func Profile(m models.Profile) resources.Profile {
	resp := resources.Profile{
		Data: resources.ProfileData{
			Id:   m.AccountID,
			Type: "profile",
			Attributes: resources.ProfileAttributes{
				Username:    m.Username,
				Pseudonym:   m.Pseudonym,
				Description: m.Description,
				Official:    m.Official,
				Avatar:      m.Avatar,
				UpdatedAt:   m.UpdatedAt,
				CreatedAt:   m.CreatedAt,
			},
		},
	}

	return resp
}

func ProfileCollection(r *http.Request, m pagi.Page[[]models.Profile]) resources.ProfilesCollection {
	data := make([]resources.ProfileData, len(m.Data))

	for i, profile := range m.Data {
		data[i] = Profile(profile).Data
	}

	links := pagi.BuildPageLinks(r, m.Page, m.Size, m.Total)

	return resources.ProfilesCollection{
		Data: data,
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}
}

func UpdateProfileSession(uploadLinks models.UpdateProfileMedia, profile models.Profile) resources.UpdateProfileSession {
	return resources.UpdateProfileSession{
		Data: resources.UpdateProfileSessionData{
			Id:   uploadLinks.UploadSessionID,
			Type: "update_profile_session",
			Attributes: resources.UpdateProfileSessionDataAttributes{
				UploadToken: uploadLinks.UploadToken,
				UploadUrl:   uploadLinks.Links.UploadURL,
				GetUrl:      uploadLinks.Links.GetURL,
			},
			Relationships: resources.UpdateProfileSessionDataRelationships{
				Profile: &resources.UpdateProfileSessionDataRelationshipsProfile{
					Data: resources.UpdateProfileSessionDataRelationshipsProfileData{
						Id:   profile.AccountID,
						Type: "profile",
					},
				},
			},
		},
		Included: []resources.ProfileData{
			Profile(profile).Data,
		},
	}
}

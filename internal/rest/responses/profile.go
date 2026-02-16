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
				Version:     m.Version,
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

func UpdateProfileSession(profile models.Profile, links models.UploadProfileMediaLinks) resources.UploadProfileMediaLinks {
	return resources.UploadProfileMediaLinks{
		Data: resources.UploadProfileMediaLinksData{
			Id:   profile.AccountID,
			Type: "profile_avatar_upload_links",
			Attributes: resources.UploadProfileMediaLinksDataAttributes{
				Avatar: resources.UploadProfileMediaLinksDataAttributesAvatar{
					Key:        links.Avatar.Key,
					UploadUrl:  links.Avatar.UploadURL,
					PreloadUrl: links.Avatar.PreloadUrl,
				},
			},
			Relationships: resources.UploadProfileMediaLinksDataRelationships{
				Profile: &resources.UploadProfileMediaLinksDataRelationshipsProfile{
					Data: resources.UploadProfileMediaLinksDataRelationshipsProfileData{
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

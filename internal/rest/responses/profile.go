package responses

import (
	"net/http"

	"github.com/netbill/profiles-svc/internal/core/models"
	resources2 "github.com/netbill/profiles-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

func Profile(m models.Profile) resources2.Profile {
	resp := resources2.Profile{
		Data: resources2.ProfileData{
			Id:   m.AccountID,
			Type: "profile",
			Attributes: resources2.ProfileAttributes{
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

func ProfileCollection(r *http.Request, m pagi.Page[[]models.Profile]) resources2.ProfilesCollection {
	data := make([]resources2.ProfileData, len(m.Data))

	for i, profile := range m.Data {
		data[i] = Profile(profile).Data
	}

	links := pagi.BuildPageLinks(r, m.Page, m.Size, m.Total)

	return resources2.ProfilesCollection{
		Data: data,
		Links: resources2.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}
}

func UploadProfileMediaLinks(profile models.Profile, links models.UploadProfileMediaLinks) resources2.UploadProfileMediaLinks {
	return resources2.UploadProfileMediaLinks{
		Data: resources2.UploadProfileMediaLinksData{
			Id:   profile.AccountID,
			Type: "profile_avatar_upload_links",
			Attributes: resources2.UploadProfileMediaLinksDataAttributes{
				Avatar: resources2.UploadProfileMediaLinksDataAttributesAvatar{
					Key:        links.Avatar.Key,
					UploadUrl:  links.Avatar.UploadURL,
					PreloadUrl: links.Avatar.PreloadUrl,
				},
			},
			Relationships: resources2.UploadProfileMediaLinksDataRelationships{
				Profile: &resources2.UploadProfileMediaLinksDataRelationshipsProfile{
					Data: resources2.UploadProfileMediaLinksDataRelationshipsProfileData{
						Id:   profile.AccountID,
						Type: "profile",
					},
				},
			},
		},
		Included: []resources2.ProfileData{
			Profile(profile).Data,
		},
	}
}

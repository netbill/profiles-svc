package responses

import (
	"net/http"

	"github.com/netbill/profiles-svc/internal/api/rest/scope"
	"github.com/netbill/profiles-svc/internal/models"
	"github.com/netbill/profiles-svc/pkg/oapi"
	"github.com/netbill/restkit/pagi"
)

func Profile(r *http.Request, m models.Profile) oapi.Profile {
	return oapi.Profile{
		Data: profileData(r, m),
	}
}

func profileData(r *http.Request, m models.Profile) oapi.ProfileData {
	res := oapi.ProfileData{
		Id:   m.AccountID,
		Type: "profile",
		Attributes: oapi.ProfileAttributes{
			Username:    m.Username,
			Pseudonym:   m.Pseudonym,
			Description: m.Description,
			Version:     m.Version,
			UpdatedAt:   m.UpdatedAt,
			CreatedAt:   m.CreatedAt,
		},
	}
	if m.AvatarKey != nil {
		url := scope.ResolverURL(r, *m.AvatarKey)
		res.Attributes.AvatarUrl = &url
	}

	return res
}

func ProfileCollection(r *http.Request, page pagi.Page[[]models.Profile]) oapi.ProfilesCollection {
	data := make([]oapi.ProfileData, len(page.Data))
	for i, profile := range page.Data {
		data[i] = profileData(r, profile)
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	return oapi.ProfilesCollection{
		Data: data,
		Links: oapi.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}
}

func UploadProfileMediaLinks(r *http.Request, profile models.Profile, links models.UploadProfileMediaLinks) oapi.UploadProfileMediaLinks {
	return oapi.UploadProfileMediaLinks{
		Data: oapi.UploadProfileMediaLinksData{
			Id:   profile.AccountID,
			Type: "profile_upload_links",
			Attributes: oapi.UploadProfileMediaLinksDataAttributes{
				Avatar: oapi.UploadProfileMediaLinksDataAttributesAvatar{
					Key:        links.Avatar.Key,
					UploadUrl:  links.Avatar.UploadURL,
					PreloadUrl: links.Avatar.PreloadUrl,
				},
			},
			Relationships: oapi.UploadProfileMediaLinksDataRelationships{
				Profile: &oapi.UploadProfileMediaLinksDataRelationshipsProfile{
					Data: oapi.UploadProfileMediaLinksDataRelationshipsProfileData{
						Id:   profile.AccountID,
						Type: "profile",
					},
				},
			},
		},
		Included: []oapi.ProfileData{
			profileData(r, profile),
		},
	}
}

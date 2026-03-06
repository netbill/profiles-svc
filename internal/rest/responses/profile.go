package responses

import (
	"net/http"

	"github.com/netbill/profiles-svc/internal/models"
	"github.com/netbill/profiles-svc/internal/rest/scope"
	"github.com/netbill/profiles-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

func Profile(r *http.Request, m models.Profile) resources.Profile {
	return resources.Profile{
		Data: profileData(r, m),
	}
}

func profileData(r *http.Request, m models.Profile) resources.ProfileData {
	res := resources.ProfileData{
		Id:   m.AccountID,
		Type: "profile",
		Attributes: resources.ProfileAttributes{
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

func ProfileCollection(r *http.Request, page pagi.Page[[]models.Profile]) resources.ProfilesCollection {
	data := make([]resources.ProfileData, len(page.Data))
	for i, profile := range page.Data {
		data[i] = profileData(r, profile)
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

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

func UploadProfileMediaLinks(r *http.Request, profile models.Profile, links models.UploadProfileMediaLinks) resources.UploadProfileMediaLinks {
	return resources.UploadProfileMediaLinks{
		Data: resources.UploadProfileMediaLinksData{
			Id:   profile.AccountID,
			Type: "profile_upload_links",
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
			profileData(r, profile),
		},
	}
}

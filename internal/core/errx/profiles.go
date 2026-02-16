package errx

import (
	"github.com/netbill/ape"
)

var (
	ErrorProfileNotExists = ape.DeclareError("PROFILE_NOT_EXISTS")

	ErrorNoContentUploaded = ape.DeclareError("NO_CONTENT_UPLOADED")

	ErrorProfileAvatarKeyIsInvalid        = ape.DeclareError("PROFILE_AVATAR_KEY_IS_INVALID")
	ErrorProfileAvatarContentIsExceedsMax = ape.DeclareError("PROFILE_AVATAR_CONTENT_EXCEEDS_MAX")
	ErrorProfileAvatarResolutionIsInvalid = ape.DeclareError("PROFILE_AVATAR_RESOLUTION_IS_INVALID")
	ErrorProfileAvatarFormatIsNotAllowed  = ape.DeclareError("PROFILE_AVATAR_FORMAT_IS_NOT_ALLOWED")
)

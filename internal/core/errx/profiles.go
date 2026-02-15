package errx

import (
	"github.com/netbill/ape"
)

var (
	ErrorProfileNotExists = ape.DeclareError("PROFILE_NOT_EXISTS")

	ErrorNoContentUploaded = ape.DeclareError("NO_CONTENT_UPLOADED")

	ErrorProfileAvatarKeyIsInvalid     = ape.DeclareError("PROFILE_AVATAR_KEY_IS_INVALID")
	ErrorProfileAvatarContentIsInvalid = ape.DeclareError("PROFILE_AVATAR_CONTENT_IS_INVALID")
)

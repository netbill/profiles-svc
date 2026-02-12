package errx

import (
	"github.com/netbill/ape"
)

var (
	ErrorProfileNotExists = ape.DeclareError("PROFILE_NOT_EXISTS")

	ErrorNoContentUploaded             = ape.DeclareError("NO_CONTENT_UPLOADED")
	ErrorProfileAvatarContentIsInvalid = ape.DeclareError("PROFILE_AVATAR_CONTENT_IS_INVALID")
)

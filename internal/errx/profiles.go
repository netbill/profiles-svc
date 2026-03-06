package errx

import (
	"github.com/netbill/ape"
)

var (
	ErrorProfileNotExists             = ape.DeclareError("PROFILE_NOT_EXISTS")
	ErrorProfileUploadedAvatarInvalid = ape.DeclareError("PROFILE_UPLOADED_AVATAR_INVALID")
)

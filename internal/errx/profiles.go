package errx

import (
	"github.com/netbill/ape"
)

var (
	ErrorProfileNotExists             = ape.DeclareError("PROFILE_NOT_EXISTS")
	ErrorProfileUploadedAvatarInvalid = ape.DeclareError("PROFILE_UPLOADED_AVATAR_INVALID")
	ErrorProfileUsernameNotValid      = ape.DeclareError("PROFILE_USERNAME_NOT_VALID")
	ErrorProfileUsernameTaken         = ape.DeclareError("PROFILE_USERNAME_TAKEN")
)

package errx

import "github.com/netbill/ape"

var (
	ErrorContentTypeIsNotAllowed = ape.DeclareError("CONTENT_TYPE_IS_NOT_ALLOWED")

	ErrorNoAvatarUpload = ape.DeclareError("NO_AVATAR_UPLOAD")
)

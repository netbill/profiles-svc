package errx

import "github.com/netbill/ape"

var (
	ErrorContentTypeIsNotAllowed = ape.DeclareError("CONTENT_TYPE_IS_NOT_ALLOWED")
	ErrorContentLengthExceed     = ape.DeclareError("CONTENT_LENGTH_EXCEED")

	ErrorNoAvatarUpload = ape.DeclareError("NO_AVATAR_UPLOAD")
)

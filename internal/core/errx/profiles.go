package errx

import (
	"github.com/netbill/restkit/ape"
)

var ErrorProfileNotFound = ape.DeclareError("PROFILE_NOT_FOUND")
var ErrorUsernameAlreadyTaken = ape.DeclareError("USERNAME_ALREADY_TAKEN")
var ErrorUsernameIsNotAllowed = ape.DeclareError("USERNAME_IS_NOT_ALLOWED")

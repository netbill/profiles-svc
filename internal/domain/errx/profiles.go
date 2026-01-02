package errx

import (
	"github.com/netbill/ape"
)

var ErrorProfileNotFound = ape.DeclareError("PROFILE_NOT_FOUND")

var ErrorProfileAlreadyExists = ape.DeclareError("PROFILE_FOR_USER_ALREADY_EXISTS")

var ErrorUsernameAlreadyTaken = ape.DeclareError("USERNAME_ALREADY_TAKEN")

var ErrorUsernameIsNotValid = ape.DeclareError("USERNAME_IS_NOT_VALID")

var ErrorSexIsNotValid = ape.DeclareError("SEX_IS_NOT_VALID")

var ErrorBirthdateIsNotValid = ape.DeclareError("BIRTHDATE_IS_NOT_VALID")

var ErrorUserTooYoung = ape.DeclareError("USER_TOO_YOUNG")

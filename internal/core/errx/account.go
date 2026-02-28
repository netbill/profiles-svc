package errx

import "github.com/netbill/ape"

var (
	ErrorAccountNotExists     = ape.DeclareError("ACCOUNT_NOT_EXISTS")
	ErrorAccountAlreadyExists = ape.DeclareError("ACCOUNT_ALREADY_EXISTS")
	ErrorAccountDeleted       = ape.DeclareError("ACCOUNT_DELETED")
)

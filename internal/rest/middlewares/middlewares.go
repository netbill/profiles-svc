package middlewares

import (
	"net/http"

	"github.com/netbill/logium"
	"github.com/netbill/restkit/mdlv"
)

type Service struct {
	accountAccessSK string
	uploadFilesSK   string

	log *logium.Logger
}

type Config struct {
	AccountAccessSK string
	UploadFilesSK   string
}

func New(
	log *logium.Logger,
	cfg Config,
) Service {
	return Service{
		accountAccessSK: cfg.AccountAccessSK,
		uploadFilesSK:   cfg.UploadFilesSK,
		log:             log,
	}
}

func (s Service) AccountAuth() func(next http.Handler) http.Handler {
	return mdlv.AccountAuth(s.log, accountDataCtxKey, s.accountAccessSK)
}

func (s Service) AccountRolesGrant(
	allowedRoles map[string]bool,
) func(http.Handler) http.Handler {
	return mdlv.AccountRoleGrant(s.log, accountDataCtxKey, allowedRoles)
}

func (s Service) ConfirmUploadFiles(scope string) func(next http.Handler) http.Handler {
	return mdlv.ConfirmUploadFiles(s.log, uploadFilesCtxKey, scope, s.uploadFilesSK)
}

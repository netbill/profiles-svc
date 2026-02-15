package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/profiles-svc/resources"
	"github.com/netbill/restkit"
)

func DeleteUploadProfileAvatar(r *http.Request) (req resources.DeleteUploadProfileAvatar, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = restkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("delete_profile_upload_avatar")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	return req, errs.Filter()
}

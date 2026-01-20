package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/profiles-svc/resources"
)

func UpdateProfileUsername(r *http.Request) (req resources.UpdateProfileUsername, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(&req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("update_profile_username")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	if chi.URLParam(r, "user_id") == req.Data.Id.String() {
		errs["data/id"] = fmt.Errorf("query user_id and body data/id do not match")
	}

	return req, errs.Filter()
}

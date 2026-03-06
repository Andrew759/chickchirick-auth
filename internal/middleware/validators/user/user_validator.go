package user

import (
	mainService "chickChirick/cmd/service"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	"chickChirick/internal/middleware/service"
	"chickChirick/internal/model/user"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type UserValidatorFactory struct{}

func (uvf UserValidatorFactory) NewValidator(dbDecorator mainService.DBDecorator, opts ...middleware.ValidatorOption) middleware.Validator {
	var vOptions middleware.ValidatorOptions
	for _, opt := range opts {
		opt(&vOptions)
	}

	return UserValidator{
		DBDecorator:      dbDecorator,
		ValidatorOptions: vOptions,
	}
}

type UserValidator struct {
	DBDecorator mainService.DBDecorator
	middleware.ValidatorOptions
}

func (uv UserValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u user.User

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			c_http.NewResponse().SendError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		if errorList := validateRequestRules(u); len(errorList) > 0 {
			errResponse := c_http.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), config.UserUserKey, &u)
		next(w, r.WithContext(ctx))
	}
}

func validateRequestRules(u user.User) []error {
	var errList []error

	if strings.TrimSpace(u.Name) == "" || !service.IsHasCorrectLength(u.Name, 256) {
		errList = append(errList, errors.New("invalid name"))
	}
	if strings.TrimSpace(u.Surname) == "" || !service.IsHasCorrectLength(u.Surname, 256) {
		errList = append(errList, errors.New("invalid surname"))
	}
	if strings.TrimSpace(u.Login) == "" || !service.IsLogin(u.Login) || !service.IsHasCorrectLength(u.Login, 256) {
		errList = append(errList, errors.New("invalid login"))
	}
	if strings.TrimSpace(u.Phone) == "" || !service.IsPhoneNumber(u.Phone) {
		errList = append(errList, errors.New("invalid phone"))
	}

	return errList
}

package user

import (
	mainService "chickChirick/cmd/service"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	"chickChirick/internal/model/user"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type PhotoValidatorFactory struct{}

func (pvf PhotoValidatorFactory) NewValidator(dbDecorator mainService.DBDecorator, opts ...middleware.ValidatorOption) middleware.Validator {
	var vOptions middleware.ValidatorOptions
	for _, opt := range opts {
		opt(&vOptions)
	}

	return &PhotoValidator{
		DBDecorator:      dbDecorator,
		ValidatorOptions: vOptions,
	}
}

type PhotoValidator struct {
	DBDecorator mainService.DBDecorator
	middleware.ValidatorOptions
}

func (pv PhotoValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p user.Photo

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			c_http.NewResponse().SendError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if errorList := pv.validateRequestRules(p); len(errorList) > 0 {
			errResponse := c_http.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		if pv.IsDBValidationActivated() {
			errContext := pv.validateAndSendResponseByDBRules(p)
			if errContext != nil {
				c_http.NewResponse().SendError(w, errContext.Message, errContext.Code)
				return
			}
		}

		ctx := context.WithValue(r.Context(), config.UserPhotoKey, &p)
		next(w, r.WithContext(ctx))
	}
}

func (pv PhotoValidator) validateRequestRules(p user.Photo) []error {
	var errList []error

	if strings.TrimSpace(p.FileUuid.String()) == "" {
		errList = append(errList, errors.New("invalid uuid"))
	}

	return errList
}

func (pv PhotoValidator) validateAndSendResponseByDBRules(p user.Photo) *middleware.ValidatorErrorContext {
	_, err := user.GetUserById(pv.DBDecorator.GormInterface, p.UserId)
	if err != nil && errors.Is(err, user.UserNotFoundErr) {
		return &middleware.ValidatorErrorContext{
			Message: err.Error(),
			Code:    http.StatusNotFound,
		}
	}
	if err != nil {
		return &middleware.ValidatorErrorContext{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	return nil
}

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

type MetaValidatorFactory struct{}

func (mvf MetaValidatorFactory) NewValidator(dbDecorator mainService.DBDecorator, opts ...middleware.ValidatorOption) middleware.Validator {
	var vOptions middleware.ValidatorOptions
	for _, opt := range opts {
		opt(&vOptions)
	}

	return &MetaValidator{
		DBDecorator:      dbDecorator,
		ValidatorOptions: vOptions,
	}
}

type MetaValidator struct {
	DBDecorator mainService.DBDecorator
	middleware.ValidatorOptions
}

func (mv MetaValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m user.Meta

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			c_http.NewResponse().SendError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if errorList := mv.validateRequestRules(m); len(errorList) > 0 {
			errResponse := c_http.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		if mv.IsDBValidationActivated() {
			errContext := mv.validateAndSendResponseByDBRules(m)
			if errContext != nil {
				c_http.NewResponse().SendError(w, errContext.Message, errContext.Code)
				return
			}
		}

		ctx := context.WithValue(r.Context(), config.UserMetaKey, &m)
		next(w, r.WithContext(ctx))
	}
}

func (mv MetaValidator) validateRequestRules(m user.Meta) []error {
	var errList []error

	if strings.TrimSpace(m.UserUuid.String()) == "" {
		errList = append(errList, errors.New("invalid meta user uuid"))
	}

	return errList
}

func (mv MetaValidator) validateAndSendResponseByDBRules(m user.Meta) *middleware.ValidatorErrorContext {
	_, err := user.GetUserById(mv.DBDecorator.GormInterface, m.UserId)
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

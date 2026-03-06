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

type CreatePropertyValidatorFactory struct{}

func (pvf CreatePropertyValidatorFactory) NewValidator(dbDecorator mainService.DBDecorator, opts ...middleware.ValidatorOption) middleware.Validator {
	var vOptions middleware.ValidatorOptions
	for _, opt := range opts {
		opt(&vOptions)
	}

	return &CreatePropertyValidator{
		DBDecorator:      dbDecorator,
		ValidatorOptions: vOptions,
	}
}

type CreatePropertyValidator struct {
	DBDecorator mainService.DBDecorator
	middleware.ValidatorOptions
}

func (pv CreatePropertyValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p user.Property

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

		ctx := context.WithValue(r.Context(), config.UserPropertyKey, &p)
		next(w, r.WithContext(ctx))
	}
}

func (pv CreatePropertyValidator) validateRequestRules(p user.Property) []error {
	var errList []error

	if strings.TrimSpace(p.Email) != "" && !service.IsEmail(p.Email) {
		errList = append(errList, errors.New("invalid email"))
	}
	if p.Password != nil && !service.IsHasCorrectLength(*p.Password, 1024) {
		errList = append(errList, errors.New("invalid password"))
	}

	return errList
}

func (pv CreatePropertyValidator) validateAndSendResponseByDBRules(p user.Property) *middleware.ValidatorErrorContext {
	_, err := user.GetUserById(pv.DBDecorator.GormInterface, p.UserId)
	if err != nil && errors.Is(err, user.UserNotFoundErr) {
		return &middleware.ValidatorErrorContext{
			Message: err.Error(),
			Code:    http.StatusNotFound,
		}
	}
	_, err = user.HasProperty(pv.DBDecorator.GormInterface, p)
	if err != nil && errors.Is(err, user.PropertyForUserAlreadyExistsErr) {
		return &middleware.ValidatorErrorContext{
			Message: err.Error(),
			Code:    http.StatusConflict,
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

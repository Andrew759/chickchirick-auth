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
)

type BanValidatorFactory struct{}

func (bvf BanValidatorFactory) NewValidator(dbDecorator mainService.DBDecorator, opts ...middleware.ValidatorOption) middleware.Validator {
	var vOptions middleware.ValidatorOptions
	for _, opt := range opts {
		opt(&vOptions)
	}

	return &BanValidator{
		DBDecorator:      dbDecorator,
		ValidatorOptions: vOptions,
	}
}

type BanValidator struct {
	DBDecorator mainService.DBDecorator
	middleware.ValidatorOptions
}

func (bv BanValidator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b user.Ban

		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			c_http.NewResponse().SendError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if errorList := bv.validateRequestRules(b); len(errorList) > 0 {
			errResponse := c_http.NewResponse()
			errResponse.AddErrorsToErrorContainer(errorList)

			errResponse.Send(w, http.StatusBadRequest)
			return
		}

		if bv.IsDBValidationActivated() {
			errContext := bv.validateAndSendResponseByDBRules(b)
			if errContext != nil {
				c_http.NewResponse().SendError(w, errContext.Message, errContext.Code)
				return
			}
		}

		ctx := context.WithValue(r.Context(), config.UserBanKey, &b)
		next(w, r.WithContext(ctx))
	}
}

func (bv BanValidator) validateRequestRules(b user.Ban) []error {
	var errList []error

	if b.BannedUserId == 0 {
		errList = append(errList, errors.New("invalid meta user uuid"))
	}

	return errList
}

func (bv BanValidator) validateAndSendResponseByDBRules(b user.Ban) *middleware.ValidatorErrorContext {
	_, err := user.GetUserById(bv.DBDecorator.GormInterface, b.UserId)
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

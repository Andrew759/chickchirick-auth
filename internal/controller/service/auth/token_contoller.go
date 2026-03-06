package auth

import (
	"chickchirick-auth/internal/controller/c_controller"
	token "chickchirick-auth/internal/model/auth"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenController struct {
	Controller c_controller.Controller
}

func (tc *TokenController) HandleRequest() {
	tc.Controller.E.GET("/tokens", func(c *gin.Context) {
		tc.GetTokens(c)
	})

	tc.Controller.ServeMux.HandleFunc("POST /token", func(w http.ResponseWriter, r *http.Request) {
		tc.CreateToken(w, c_http.NewRequest(r))
	})

	tc.Controller.ServeMux.HandleFunc("GET /token/{id}", func(w http.ResponseWriter, r *http.Request) {
		tc.GetToken(w, c_http.NewRequest(r))
	})

	tc.Controller.ServeMux.HandleFunc("PUT /token/{id}", func(w http.ResponseWriter, r *http.Request) {
		tc.UpdateToken(w, c_http.NewRequest(r))
	})

	tc.Controller.ServeMux.HandleFunc("DELETE /token/{id}", func(w http.ResponseWriter, r *http.Request) {
		tc.DeleteToken(w, c_http.NewRequest(r))
	})
}

func (tc *TokenController) GetTokens(c *gin.Context) {
	tokens, err := token.GetTokens(tc.Controller.DI.DBDecorator.GDB())
	if err != nil {
		c.String()
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, tokens, http.StatusOK)
}

func (tc *TokenController) GetToken(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := token.GetTokenById(tc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Token not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, t, http.StatusOK)
}

func (tc *TokenController) CreateToken(w http.ResponseWriter, r *c_http.Request) {
	var t token.Token
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := token.CreateToken(tc.Controller.Dependencies.DBDecorator.GDB(), &t); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, t, http.StatusCreated)
}

func (tc *TokenController) UpdateToken(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var t token.Token
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = token.UpdateTokenById(tc.Controller.Dependencies.DBDecorator.GDB(), &t, id)
	if err != nil && errors.Is(err, token.TokenNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, t, http.StatusOK)
}

func (tc *TokenController) DeleteToken(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = token.DeleteTokenById(tc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, token.TokenNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

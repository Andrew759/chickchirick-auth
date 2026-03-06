package auth

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	session "chickChirick/internal/model/auth"
	"encoding/json"
	"errors"
	"net/http"
)

type SessionController struct {
	Controller c_controller.Controller
}

func (sc *SessionController) HandleRequest() {
	sc.Controller.ServeMux.HandleFunc("GET /sessions", func(w http.ResponseWriter, r *http.Request) {
		sc.GetSessions(w)
	})

	sc.Controller.ServeMux.HandleFunc("POST /session", func(w http.ResponseWriter, r *http.Request) {
		sc.CreateSession(w, c_http.NewRequest(r))
	})

	sc.Controller.ServeMux.HandleFunc("GET /session/{id}", func(w http.ResponseWriter, r *http.Request) {
		sc.GetSession(w, c_http.NewRequest(r))
	})

	sc.Controller.ServeMux.HandleFunc("PUT /session/{id}", func(w http.ResponseWriter, r *http.Request) {
		sc.UpdateSession(w, c_http.NewRequest(r))
	})

	sc.Controller.ServeMux.HandleFunc("DELETE /session/{id}", func(w http.ResponseWriter, r *http.Request) {
		sc.DeleteSession(w, c_http.NewRequest(r))
	})
}

func (sc *SessionController) GetSessions(w http.ResponseWriter) {
	codes, err := session.GetSessions(sc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, codes, http.StatusOK)
}

func (sc *SessionController) GetSession(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := session.GetSessionById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Session not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, c, http.StatusOK)
}

func (sc *SessionController) CreateSession(w http.ResponseWriter, r *c_http.Request) {
	var s session.Session
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := session.CreateSession(sc.Controller.Dependencies.DBDecorator.GDB(), &s); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, s, http.StatusCreated)
}

func (sc *SessionController) UpdateSession(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var s session.Session
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = session.UpdateSessionById(sc.Controller.Dependencies.DBDecorator.GDB(), &s, id)
	if err != nil && errors.Is(err, session.SessionNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, s, http.StatusOK)
}

func (sc *SessionController) DeleteSession(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = session.DeleteSessionById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, session.SessionNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

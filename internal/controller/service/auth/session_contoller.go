package auth

import (
	"chickchirick-auth/internal/controller/c_controller"
	session "chickchirick-auth/internal/model/auth"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SessionController struct {
	Controller c_controller.Controller
}

func (sc *SessionController) RegisterRoutes() {
	e := sc.Controller.E

	e.GET("/sessions", sc.GetSessions)
	e.POST("/session", sc.CreateSession)
	e.GET("/session/:id", sc.GetSession)
	e.PUT("/session/:id", sc.UpdateSession)
	e.DELETE("/session/:id", sc.DeleteSession)
}

func (sc *SessionController) GetSessions(c *gin.Context) {
	sessions, err := session.GetSessions(sc.Controller.DI.DBDecorator.GDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (sc *SessionController) GetSession(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	iId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	s, err := session.GetSessionById(sc.Controller.DI.DBDecorator.GDB(), iId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (sc *SessionController) CreateSession(c *gin.Context) {
	var s session.Session

	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := session.CreateSession(sc.Controller.DI.DBDecorator.GDB(), &s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, s)
}

func (sc *SessionController) UpdateSession(c *gin.Context) {
	id := c.Param("id")

	var s session.Session
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	iId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = session.UpdateSessionById(sc.Controller.DI.DBDecorator.GDB(), &s, iId)
	if err != nil {
		if errors.Is(err, session.SessionNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, s)
}

func (sc *SessionController) DeleteSession(c *gin.Context) {
	id := c.Param("id")

	iId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = session.DeleteSessionById(sc.Controller.DI.DBDecorator.GDB(), iId)
	if err != nil {
		if errors.Is(err, session.SessionNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session: " + err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

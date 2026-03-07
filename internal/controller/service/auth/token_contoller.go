package auth

import (
	"chickchirick-auth/internal/controller/c_controller"
	token "chickchirick-auth/internal/model/auth"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TokenController struct {
	Controller c_controller.Controller
}

func (tc *TokenController) RegisterRoutes() {
	e := tc.Controller.E

	e.GET("/tokens", tc.GetTokens)
	e.POST("/token", tc.CreateToken)
	e.GET("/token/:id", tc.GetToken)
	e.PUT("/token/:id", tc.UpdateToken)
	e.DELETE("/token/:id", tc.DeleteToken)
}

func (tc *TokenController) GetTokens(c *gin.Context) {
	tokens, err := token.GetTokens(tc.Controller.DI.DBDecorator.GDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (tc *TokenController) GetToken(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	iId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	t, err := token.GetTokenById(tc.Controller.DI.DBDecorator.GDB(), iId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (tc *TokenController) CreateToken(c *gin.Context) {
	var t token.Token
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := token.CreateToken(tc.Controller.DI.DBDecorator.GDB(), &t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func (tc *TokenController) UpdateToken(c *gin.Context) {
	id := c.Param("id")
	iId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var t token.Token
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	err = token.UpdateTokenById(tc.Controller.DI.DBDecorator.GDB(), &t, iId)
	if err != nil {
		if errors.Is(err, token.TokenNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, t)
}

func (tc *TokenController) DeleteToken(c *gin.Context) {
	id := c.Param("id")
	iId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = token.DeleteTokenById(tc.Controller.DI.DBDecorator.GDB(), iId)
	if err != nil {
		if errors.Is(err, token.TokenNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token: " + err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

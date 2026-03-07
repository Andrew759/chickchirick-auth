package auth

import (
	"chickchirick-auth/internal/controller/c_controller"
	code "chickchirick-auth/internal/model/auth"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CodeController struct {
	Controller c_controller.Controller
}

func (cc *CodeController) RegisterRoutes() {
	e := cc.Controller.E

	e.GET("/codes", cc.GetCodes)
	e.POST("/code", cc.CreateCode)
	e.GET("/code/:id", cc.GetCode)
	e.PUT("/code/:id", cc.UpdateCode)
	e.DELETE("/code/:id", cc.DeleteCode)
}

func (cc *CodeController) GetCodes(c *gin.Context) {
	cds, err := code.GetCodes(cc.Controller.DI.DBDecorator.GDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cds)
}

func (cc *CodeController) GetCode(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	iId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	cd, err := code.GetCodeById(cc.Controller.DI.DBDecorator.GDB(), iId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Code not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, cd)
}

func (cc *CodeController) CreateCode(c *gin.Context) {
	var cd code.Code
	if err := c.ShouldBindJSON(&c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := code.CreateCode(cc.Controller.DI.DBDecorator.GDB(), &cd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create code: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cd)
}

func (cc *CodeController) UpdateCode(c *gin.Context) {
	id := c.Param("id")
	iId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var cd code.Code
	if err := c.ShouldBindJSON(&c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	err = code.UpdateCodeById(cc.Controller.DI.DBDecorator.GDB(), &cd, iId)
	if err != nil {
		if errors.Is(err, code.CodeNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update code: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, cd)
}

func (cc *CodeController) DeleteCode(c *gin.Context) {
	id := c.Param("id")
	iId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = code.DeleteCodeById(cc.Controller.DI.DBDecorator.GDB(), iId)
	if err != nil {
		if errors.Is(err, code.CodeNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete code: " + err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

package auth

import (
	"chickchirick-auth/internal/controller/c_controller"
	user "chickchirick-auth/internal/model/auth"
	"chickchirick-auth/internal/request"
	"chickchirick-auth/internal/service"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Controller c_controller.Controller
}

func (uc *UserController) RegisterRoutes() {
	e := uc.Controller.E

	group := e.Group("/user")
	//group.Use(middleware.AuthMiddleware())
	group.POST("/login", uc.GetUser)
	group.POST("/create", uc.CreateUser)
}

func (uc *UserController) GetUser(c *gin.Context) {
	var req request.GetUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := user.GetUserByUuidAndPass(
		uc.Controller.DI.DBDecorator.GDB(),
		req.UserUuid,
		req.Password,
	)

	if err != nil && errors.Is(err, user.UserNotFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, u)
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var u user.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := user.CreateUser(uc.Controller.DI.DBDecorator.GDB(), &u); err != nil &&
		errors.Is(err, user.UserAlreadyExistsErr) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	at, rt, err := service.CreateTokens(ctx, *uc.Controller.DI.RedisDecorator, u.UserUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user token: " + err.Error()})
		return
	}

	c.SetCookie("access_token", at.Token, int(at.Lt.Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", rt.Token, int(rt.Lt.Seconds()), "/auth/refresh", "", false, true)

	c.JSON(http.StatusCreated, u)
}

package auth

import (
	"chickchirick-auth/internal/controller/c_controller"
	token "chickchirick-auth/internal/model/auth"
	"chickchirick-auth/internal/request"
	"chickchirick-auth/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Controller c_controller.Controller
}

func (ac *AuthController) RegisterRoutes() {
	e := ac.Controller.E

	e.POST("/auth/login", ac.Login)
	e.POST("/auth/refresh", ac.RefreshToken)
	e.POST("/auth/logout", ac.Logout)
}

func (ac *AuthController) Login(c *gin.Context) {
	var ctr request.CreateTokenRequest
	if err := c.ShouldBindJSON(&ctr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := token.GetUserByUuidAndPass(ac.Controller.DI.DBDecorator.GDB(), ctr.UserUuid, ctr.Password)
	if err != nil && errors.Is(err, token.UserNotFoundErr) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	at, rt, err := service.CreateTokens(c.Request.Context(), *ac.Controller.DI.RedisDecorator, user.UserUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tokens"})
		return
	}

	ac.setTokenCookies(c, at.Token, rt.Token, int(at.Lt.Seconds()), int(rt.Lt.Seconds()))

	c.JSON(http.StatusOK, gin.H{"payload": "Successfully logged in"})
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	oldRtStr, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	at, rt, err := service.RefreshToken(c.Request.Context(), *ac.Controller.DI.RedisDecorator, oldRtStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired or invalid"})
		return
	}

	ac.setTokenCookies(c, at.Token, rt.Token, int(at.Lt.Seconds()), int(rt.Lt.Seconds()))

	c.Status(http.StatusNoContent)
}

func (ac *AuthController) Logout(c *gin.Context) {
	rtStr, err := c.Cookie("refresh_token")
	if err == nil {
		_ = service.Logout(c.Request.Context(), *ac.Controller.DI.RedisDecorator, rtStr)
	}

	//Отчистка cookie
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/auth/refresh", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func (ac *AuthController) setTokenCookies(c *gin.Context, at, rt string, atMaxAge, rtMaxAge int) {
	c.SetCookie("access_token", at, atMaxAge, "/", "", false, true)

	//Refresh Token доступен только по пути /auth/refresh
	c.SetCookie("refresh_token", rt, rtMaxAge, "/auth/refresh", "", false, true)
}

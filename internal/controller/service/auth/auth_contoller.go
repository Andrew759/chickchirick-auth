package auth

import (
	"chickchirick-auth/internal/controller/c_controller"
	token "chickchirick-auth/internal/model/auth"
	"chickchirick-auth/internal/request"
	"chickchirick-auth/internal/service"
	"chickchirick-auth/pkg/chirick_config"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type AuthController struct {
	Controller c_controller.Controller
}

func (ac *AuthController) RegisterRoutes() {
	e := ac.Controller.E

	group := e.Group("/auth")

	group.POST("/login", ac.Login)
	group.POST("/refresh", ac.RefreshToken)
	group.POST("/logout", ac.Logout)
	group.GET("/validate", ac.ValidateToken)
}

func (ac *AuthController) Login(c *gin.Context) {
	var ctr request.CreateTokenRequest
	if err := c.ShouldBindJSON(&ctr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user, err := token.GetUserByUuidAndPass(ac.Controller.DI.DBDecorator.GDB(), ctr.UserUuid, ctr.Password)
	if err != nil && errors.Is(err, token.UserNotFoundErr) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	at, rt, err := service.CreateTokens(c.Request.Context(), *ac.Controller.DI.RedisDecorator, user.UserUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tokens"})
		return
	}

	ac.setTokenCookies(c, at.Token, rt.Token, int(at.Lt.Seconds()), int(rt.Lt.Seconds()))

	c.JSON(http.StatusOK, gin.H{"payload": "successfully logged in"})
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	oldRtStr, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	at, rt, err := service.RefreshToken(c.Request.Context(), *ac.Controller.DI.RedisDecorator, oldRtStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired or invalid"})
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

	c.JSON(http.StatusOK, gin.H{"payload": "Logged out"})
}

func (ac *AuthController) setTokenCookies(c *gin.Context, at, rt string, atMaxAge, rtMaxAge int) {
	c.SetCookie("access_token", at, atMaxAge, "/", "", false, true)

	//Refresh Token доступен только по пути /auth/refresh
	c.SetCookie("refresh_token", rt, rtMaxAge, "/auth/refresh", "", false, true)
}

// ValidateToken проверяет токен из других микросервисов
func (ac *AuthController) ValidateToken(c *gin.Context) {
	tokenStr, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		return
	}

	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	claims := &service.Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString(chirick_config.SecretKey)), nil
	})

	if err != nil || !t.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payload": gin.H{
			"valid":     true,
			"user_uuid": claims.UserUuid,
		},
	})
}

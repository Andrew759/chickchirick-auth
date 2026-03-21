package factory

import (
	"chickchirick-auth/internal/controller/c_controller"
	internalService "chickchirick-auth/internal/controller/service/auth"

	"github.com/gin-gonic/gin"
)

type AuthServer struct{}

func InitAuthServer(e *gin.Engine, aDIC *c_controller.DIContainer) {
	authServer := AuthServer{}

	authServer.initAuthService(e, aDIC)
	authServer.initUserService(e, aDIC)
}

func (as AuthServer) initAuthService(e *gin.Engine, aDIC *c_controller.DIContainer) internalService.AuthController {
	authService := internalService.AuthController{
		Controller: c_controller.Controller{
			E:  e,
			DI: aDIC,
		},
	}
	authService.RegisterRoutes()

	return authService
}

func (as AuthServer) initUserService(e *gin.Engine, aDIC *c_controller.DIContainer) internalService.UserController {
	userService := internalService.UserController{
		Controller: c_controller.Controller{
			E:  e,
			DI: aDIC,
		},
	}
	userService.RegisterRoutes()

	return userService
}

package factory

import (
	"chickchirick-auth/internal/controller/c_controller"
	internalService "chickchirick-auth/internal/controller/service/auth"

	"github.com/gin-gonic/gin"
)

type AuthServer struct{}

func InitAuthServer(e *gin.Engine, aDIC *c_controller.DIContainer) {
	authServer := AuthServer{}

	authServer.initCodeService(e, aDIC)
	authServer.initSessionService(e, aDIC)
	authServer.initTokenService(e, aDIC)
}

func (as AuthServer) initCodeService(e *gin.Engine, aDIC *c_controller.DIContainer) internalService.CodeController {
	codeService := internalService.CodeController{
		Controller: c_controller.Controller{
			E:  e,
			DI: aDIC,
		},
	}
	codeService.RegisterRoutes()

	return codeService
}

func (as AuthServer) initSessionService(e *gin.Engine, aDIC *c_controller.DIContainer) internalService.SessionController {
	sessionService := internalService.SessionController{
		Controller: c_controller.Controller{
			E:  e,
			DI: aDIC,
		},
	}
	sessionService.RegisterRoutes()

	return sessionService
}

func (as AuthServer) initTokenService(e *gin.Engine, aDIC *c_controller.DIContainer) internalService.TokenController {
	tokenService := internalService.TokenController{
		Controller: c_controller.Controller{
			E:  e,
			DI: aDIC,
		},
	}
	tokenService.RegisterRoutes()

	return tokenService
}

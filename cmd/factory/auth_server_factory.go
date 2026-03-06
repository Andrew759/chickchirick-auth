package factory

import (
	"chickChirick-auth/internal/controller/c_controller"
	internalService "chickChirick/internal/controller/service/auth"
	"net/http"
)

type AuthServer struct {
	*http.ServeMux
	c_controller.DIContainer
}

func InitAuthServer(mux *http.ServeMux, abstractDiContainer c_controller.DIContainer) AuthServer {
	authServer := AuthServer{
		ServeMux:    mux,
		DIContainer: abstractDiContainer,
	}

	authServer.initCodeService()
	authServer.initSessionService()
	authServer.initTokenService()

	return authServer
}

func (as AuthServer) initCodeService() internalService.CodeController {
	codeService := internalService.CodeController{
		Controller: c_controller.Controller{
			ServeMux:     as.ServeMux,
			Dependencies: as.DIContainer,
		},
	}
	codeService.HandleRequest()

	return codeService
}

func (as AuthServer) initSessionService() internalService.SessionController {
	sessionService := internalService.SessionController{
		Controller: c_controller.Controller{
			ServeMux:     as.ServeMux,
			Dependencies: as.DIContainer,
		},
	}
	sessionService.HandleRequest()

	return sessionService
}

func (as AuthServer) initTokenService() internalService.TokenController {
	tokenService := internalService.TokenController{
		Controller: c_controller.Controller{
			ServeMux:     as.ServeMux,
			Dependencies: as.DIContainer,
		},
	}
	tokenService.HandleRequest()

	return tokenService
}

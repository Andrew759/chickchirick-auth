package factory

import (
	"chickchirick-auth/cmd/service"
	"chickchirick-auth/internal/controller/c_controller"
	//TODO: подумать - оставить или удалить профилировщик
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
)

func BuildAndServe(dbDecorator *service.DBDecorator, redisDecorator *service.RedisDecorator) {
	err := BuildServer(dbDecorator, redisDecorator)
	if err != nil {
		panic(err)
	}
}

func BuildServer(dbDecorator *service.DBDecorator, redisDecorator *service.RedisDecorator) error {
	e := gin.Default()
	InitAuthServer(e, &c_controller.DIContainer{DBDecorator: dbDecorator, RedisDecorator: redisDecorator})

	err := e.Run()
	if err != nil {
		return err
	}

	return nil
}

package factory

import (
	"chickchirick-auth/cmd/service"
	"chickchirick-auth/internal/controller/c_controller"
	"chickchirick-auth/pkg/chirick_config"

	//TODO: подумать - оставить или удалить профилировщик
	_ "net/http/pprof"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func BuildAndServe(dbDecorator *service.DBDecorator, redisDecorator *service.RedisDecorator) error {
	return BuildServer(dbDecorator, redisDecorator)
}

func BuildServer(dbDecorator *service.DBDecorator, redisDecorator *service.RedisDecorator) error {
	e := gin.Default()

	//TODO: доработать CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{viper.GetString(chirick_config.UserApp)}

	e.Use(cors.New(config))

	InitAuthServer(e, &c_controller.DIContainer{DBDecorator: dbDecorator, RedisDecorator: redisDecorator})

	err := e.Run()
	if err != nil {
		return err
	}

	return nil
}

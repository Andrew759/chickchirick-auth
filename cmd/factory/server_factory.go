package factory

import (
	"chickchirick-auth/cmd/service"
	"net/http"
	//TODO: подумать - оставить или удалить профилировщик
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
)

func BuildAndServe(dbDecorator service.DBDecorator, redisDecorator service.RedisDecorator) {
	err := BuildServer(dbDecorator, redisDecorator)
	if err != nil {
		panic(err)
	}
}

func BuildServer(dbDecorator service.DBDecorator, redisDecorator service.RedisDecorator) error {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err := r.Run()
	if err != nil {
		return err
	}

	return nil
}

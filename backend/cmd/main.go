package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/notblinkyet/docker-pinger/backend/internal/config"
	"github.com/notblinkyet/docker-pinger/backend/src/app"
	"github.com/notblinkyet/docker-pinger/backend/src/services"
	"github.com/notblinkyet/docker-pinger/backend/src/storage"
	handlers "github.com/notblinkyet/docker-pinger/backend/src/transport"
)

func main() {
	cfg := config.MustLoad()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	layerStorage, err := storage.Open(cfg)
	if err != nil {
		panic(err)
	}
	logger.Println("success connect to DB")
	defer layerStorage.Close()
	layerService := services.NewServices(layerStorage)
	layerHandler := handlers.NewHandlers(layerService)

	router := gin.Default()
	api := router.Group("/api")
	{
		containers := api.Group("/containers")
		{
			containers.GET("", layerHandler.Containers.GetAll)
			containers.POST("", layerHandler.Containers.Create)
			containers.DELETE(":ip", layerHandler.Containers.Delete)
		}
		pings := api.Group("/pings")
		{
			pings.GET("", layerHandler.Pings.GetAll)
			pings.GET("/last", layerHandler.Pings.GetLast)
			pings.POST("", layerHandler.Pings.Create)
		}
	}
	app := app.New(router, cfg.Server.Port, cfg.Server.Host, cfg.Server.TimeOut)

	go func() {
		app.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	logger.Println("app is working")

	<-stop

	if err := app.Stop(cfg.Server.TimeOut); err != nil {
		panic(err)
	}
	logger.Println("app is closed")
}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/notblinkyet/docker-pinger/backend/internal/api/pinger"
	"github.com/notblinkyet/docker-pinger/backend/internal/app"
	"github.com/notblinkyet/docker-pinger/backend/internal/config"
	"github.com/notblinkyet/docker-pinger/backend/internal/services"
	"github.com/notblinkyet/docker-pinger/backend/internal/storage"
	handlers "github.com/notblinkyet/docker-pinger/backend/internal/transport"
)

func main() {
	cfg := config.MustLoad()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger.Println(cfg)
	layerStorage, err := storage.Open(cfg)
	if err != nil {
		panic(err)
	}
	logger.Println("success connect to DB")
	defer layerStorage.Close()
	pingerApi := pinger.NewPingerApi(&cfg.PingerApi)
	layerService := services.NewServices(layerStorage, pingerApi, logger)
	layerHandler := handlers.NewHandlers(layerService)

	router := gin.Default()
	api := router.Group("/backend")
	{
		containers := api.Group("/containers")
		{
			containers.GET("", layerHandler.Container.GetAll)
			containers.POST("", layerHandler.Container.Create)
			containers.DELETE(":ip", layerHandler.Container.Delete)
		}
		pings := api.Group("/pings")
		{
			pings.GET("", layerHandler.Ping.GetAll)
			pings.GET("/last", layerHandler.Ping.GetLast)
			pings.POST("", layerHandler.Ping.Create)
		}
	}
	app := app.New(router, cfg.Server.Port, cfg.Server.Host, cfg.Server.TimeOut)

	go func() {
		app.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	logger.Println("backend app is working")

	<-stop

	if err := app.Stop(cfg.Server.TimeOut); err != nil {
		panic(err)
	}
	logger.Println("backend app is closed")
}

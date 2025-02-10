package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/notblinkyet/docker-pinger/pinger/internal/app"
	"github.com/notblinkyet/docker-pinger/pinger/internal/client"
	"github.com/notblinkyet/docker-pinger/pinger/internal/config"
	"github.com/notblinkyet/docker-pinger/pinger/internal/handlers"
	"github.com/notblinkyet/docker-pinger/pinger/internal/pinger"
	"github.com/notblinkyet/docker-pinger/pinger/internal/services"
)

func main() {
	cfg := config.MustLoad()
	client := client.NewApi(&cfg.Clinet)
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	_, _ = client, logger
	ips, err := client.GetContainers()
	if err != nil {
		logger.Fatal(err)
	}
	pinger := pinger.NewPinger(ips, client, logger)
	layerService := services.NewServices(pinger.Ips)
	layerHandler := handlers.NewHandlers(layerService)
	router := gin.Default()
	main := router.Group("/pinger")
	{
		pings := main.Group("/ping")
		{
			pings.POST("", layerHandler.Ping.Create)
			pings.DELETE(":ip", layerHandler.Ping.Delete)
		}
	}

	app := app.NewApp(router, cfg.Server.Port, cfg.Server.Host, cfg.Server.Timeout)

	go app.Run()

	go pinger.PingOnceAfterDelay(cfg.Delay)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	logger.Println("pinger app is working")
	<-stop
	if err := app.Stop(cfg.Server.Timeout); err != nil {
		logger.Fatal(err)
	}
	logger.Println("pinger app is closed")
}

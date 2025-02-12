package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/notblinkyet/docker-pinger/pinger/internal/app"
	"github.com/notblinkyet/docker-pinger/pinger/internal/client"
	"github.com/notblinkyet/docker-pinger/pinger/internal/config"
	"github.com/notblinkyet/docker-pinger/pinger/internal/handlers"
	"github.com/notblinkyet/docker-pinger/pinger/internal/pinger"
	"github.com/notblinkyet/docker-pinger/pinger/internal/redis"
	"github.com/notblinkyet/docker-pinger/pinger/internal/services"
)

func main() {
	cfg := config.MustLoad()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger.Println(cfg)
	client := client.NewApi(&cfg.Clinet)
	ips, err := client.GetContainers()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(ips)
	redis, err := redis.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Fatal(err)
	}
	pinger := pinger.NewPinger(ips, client, logger, redis)
	layerService := services.NewServices(pinger.Ips)
	layerHandler := handlers.NewHandlers(layerService)
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(
		cors.New(cors.Config{
			AllowOrigins: []string{fmt.Sprintf("http://%s:%d", cfg.Clinet.Host, cfg.Clinet.Port)},
			AllowMethods: []string{"POST", "DELETE"},
			MaxAge:       12 * time.Hour,
		}),
		gin.Logger(),
	)
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

	go pinger.StartPinging(cfg.Delta)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	logger.Println("pinger app is working")
	<-stop
	if err := app.Stop(cfg.Server.Timeout); err != nil {
		logger.Fatal(err)
	}
	logger.Println("pinger app is closed")
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/notblinkyet/docker-pinger/backend/internal/services"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
)

type IPingHandler interface {
	GetAll(ctx *gin.Context)
	GetLast(ctx *gin.Context)
	Create(ctx *gin.Context)
}

type PingHandler struct {
	PingService services.IPingService
}

func NewPingHandler(service services.IPingService) *PingHandler {
	return &PingHandler{
		PingService: service,
	}
}

func (handler *PingHandler) GetAll(ctx *gin.Context) {
	pings, err := handler.PingService.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"data": pings})
}

func (handler *PingHandler) GetLast(ctx *gin.Context) {
	pings, err := handler.PingService.GetLast()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"data": pings})
}

func (handler *PingHandler) Create(ctx *gin.Context) {
	var data []models.Ping
	if err := ctx.ShouldBindBodyWithJSON(data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	handler.PingService.Create(data)
	ctx.AbortWithStatus(http.StatusOK)
}

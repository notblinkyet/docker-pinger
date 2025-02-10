package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
	"github.com/notblinkyet/docker-pinger/pinger/internal/services"
)

var (
	ErrBadIp error = errors.New("ip string is empty")
)

type IPingHandler interface {
	Delete(ctx *gin.Context)
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
func (handler *PingHandler) Create(ctx *gin.Context) {
	var container models.Container
	err := json.NewDecoder(ctx.Request.Body).Decode(&container)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := handler.PingService.Create(container.Ip); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

func (handler *PingHandler) Delete(ctx *gin.Context) {
	ip := ctx.Param("ip")
	if ip == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": ErrBadIp.Error()})
		return
	}
	if err := handler.PingService.Delete(ip); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

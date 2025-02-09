package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
	"github.com/notblinkyet/docker-pinger/backend/src/services"
)

type IContainerHandler interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type ContainerHandler struct {
	ContainerService services.IContainerService
}

func NewContainerHandler(service services.IContainerService) *ContainerHandler {
	return &ContainerHandler{
		ContainerService: service,
	}
}

func (handler *ContainerHandler) Create(ctx *gin.Context) {
	var container models.Container
	if err := ctx.ShouldBindBodyWithJSON(&container); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := handler.ContainerService.Create(container.Ip); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

func (handler *ContainerHandler) GetAll(ctx *gin.Context) {
	containers, err := handler.ContainerService.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"data": containers})
}

func (handler *ContainerHandler) Delete(ctx *gin.Context) {
	ip := ctx.Param("ip")
	if ip == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := handler.ContainerService.Delete(ip); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

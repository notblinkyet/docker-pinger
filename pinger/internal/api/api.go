package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
	"github.com/notblinkyet/docker-pinger/pinger/internal/config"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrInternalServer = errors.New("internal server error")
)

type Api struct {
	cfg *config.Api
}

func (a *Api) GetContainers() ([]models.Container, error) {
	url := fmt.Sprint("http://%s:%d%s", a.cfg.Host, a.cfg.Port, a.cfg.GetEndpoint)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 400 {
			return nil, ErrBadRequest
		}
		return nil, ErrInternalServer
	}
	defer resp.Body.Close()
	var containers []models.Container

	err = json.NewDecoder(resp.Body).Decode(&containers)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

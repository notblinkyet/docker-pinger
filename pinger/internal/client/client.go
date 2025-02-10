package client

import (
	"bytes"
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

type Client struct {
	Cfg *config.Client
}

func NewApi(cfg *config.Client) *Client {
	return &Client{
		Cfg: cfg,
	}
}

type ContainerResponse struct {
	Data []models.Container `json:"data"`
}

func (a *Client) GetContainers() ([]string, error) {
	url := fmt.Sprintf("http://%s:%d%s", a.Cfg.Host, a.Cfg.Port, a.Cfg.GetEndpoint)
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
	var containers ContainerResponse

	err = json.NewDecoder(resp.Body).Decode(&containers)
	if err != nil {
		return nil, err
	}
	ips := make([]string, len(containers.Data))
	for i, container := range containers.Data {
		ips[i] = container.Ip
	}
	return ips, nil
}

func (a *Client) Post(pings []models.Ping) error {
	url := fmt.Sprintf("http://%s:%d%s", a.Cfg.Host, a.Cfg.Port, a.Cfg.PostEndpoint)
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(pings)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 400 {
			return ErrBadRequest
		}
		return ErrInternalServer
	}
	return nil
}

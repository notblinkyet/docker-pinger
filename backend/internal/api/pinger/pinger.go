package pinger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/notblinkyet/docker-pinger/backend/internal/config"
	"github.com/notblinkyet/docker-pinger/backend/pkg/models"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrInternalServer = errors.New("internal server error")
)

type PingerApi struct {
	cfg *config.PingerApi
}

func NewPingerApi(cfg *config.PingerApi) *PingerApi {
	return &PingerApi{
		cfg: cfg,
	}
}

func (a *PingerApi) Post(ip string) error {
	url := fmt.Sprintf("http://%s:%d%s", a.cfg.Host, a.cfg.Port, a.cfg.PostEndpoint)
	var (
		buf       bytes.Buffer
		container models.Container
	)
	container.Ip = ip
	err := json.NewEncoder(&buf).Encode(&container)
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

func (a *PingerApi) Delete(ip string) error {
	url := fmt.Sprintf("http://%s:%d%s/:%s", a.cfg.Host, a.cfg.Port, a.cfg.PostEndpoint, ip)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
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

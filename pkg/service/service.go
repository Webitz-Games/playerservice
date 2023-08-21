package service

import (
	"playerapi/pkg/config"
	"playerapi/pkg/player/api"
)

type PlayerService struct {
	Name   string
	config *config.Config
	api.RequestHandlers
}

func New(config *config.Config, name string, requestHandlers api.RequestHandlers) *PlayerService {
	return &PlayerService{
		Name:            name,
		config:          config,
		RequestHandlers: requestHandlers,
	}
}

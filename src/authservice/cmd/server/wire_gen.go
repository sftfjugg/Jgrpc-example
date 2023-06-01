// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"authservice/config"
	"authservice/service/v1/server"
)

// Injectors from wire.go:

func InitServer(cfg string) (*Server, error) {
	configConfig := config.NewConfig(cfg)
	client := NewRedis(configConfig)
	tracerProvider, err := NewTrace(configConfig)
	if err != nil {
		return nil, err
	}
	repository := serverV1.NewRepository(client, configConfig, tracerProvider)
	group := NewRunGroup()
	logger := NewLogger()
	authorizationServer := serverV1.NewServer(configConfig, client, repository, logger)
	server := NewGrpcServer(authorizationServer)
	serverServer := NewServer(repository, configConfig, group, logger, server, tracerProvider)
	return serverServer, nil
}

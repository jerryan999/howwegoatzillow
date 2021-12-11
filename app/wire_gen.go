// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/golang/mock/gomock"
	"github.com/google/wire"
	"github.com/zillow/howwegoatzillow/libs/config"
	"github.com/zillow/howwegoatzillow/libs/db"
	"github.com/zillow/howwegoatzillow/libs/http"
	"github.com/zillow/howwegoatzillow/libs/kafka"
	"github.com/zillow/howwegoatzillow/libs/logger"
	"github.com/zillow/howwegoatzillow/mocks/db"
	"github.com/zillow/howwegoatzillow/mocks/kafka"
)

// Injectors from wire.go:

func InitializeServer() (*MyServer, func()) {
	appConfig := config.NewAppConfig()
	serverConfig := NewServerConfig(appConfig)
	tracer := NewTracer()
	loggerLogger, cleanup := logger.NewLogger(tracer)
	server := NewServer(serverConfig, loggerLogger, tracer)
	httpConfig := NewHttpServiceConfig(appConfig)
	leveledLogger := http.NewLeveledLogger(loggerLogger)
	provider := http.NewClientProvider(tracer, leveledLogger)
	dbConfig := NewDbConfig(appConfig)
	dbProvider := db.NewProvider()
	kafkaConfig := NewKafkaConfig(appConfig)
	client := kafka.NewClient(kafkaConfig, tracer, loggerLogger)
	myService := &MyService{
		HTTPConfig:         httpConfig,
		HTTPClientProvider: provider,
		DBConfig:           dbConfig,
		DBProvider:         dbProvider,
		KafkaConfig:        kafkaConfig,
		KafkaClient:        client,
	}
	myServer := NewMyServer(server, myService)
	return myServer, func() {
		cleanup()
	}
}

func InitializeServerTestable(ctrl *gomock.Controller) (*ServerTestable, func()) {
	appConfig := config.NewAppConfig()
	serverConfig := NewServerConfig(appConfig)
	tracer := NewTracer()
	loggerLogger, cleanup := logger.NewLogger(tracer)
	server := NewServer(serverConfig, loggerLogger, tracer)
	httpConfig := NewHttpServiceConfig(appConfig)
	leveledLogger := http.NewLeveledLogger(loggerLogger)
	provider := http.NewClientProvider(tracer, leveledLogger)
	dbConfig := NewDbConfig(appConfig)
	mockProvider := mock_db.NewMockProvider(ctrl)
	kafkaConfig := NewKafkaConfig(appConfig)
	mockClient := mock_kafka.NewMockClient(ctrl)
	myService := &MyService{
		HTTPConfig:         httpConfig,
		HTTPClientProvider: provider,
		DBConfig:           dbConfig,
		DBProvider:         mockProvider,
		KafkaConfig:        kafkaConfig,
		KafkaClient:        mockClient,
	}
	myServer := NewMyServer(server, myService)
	mockWriter := mock_kafka.NewMockWriter(ctrl)
	serverTestable := &ServerTestable{
		Server:     myServer,
		DBProvider: mockProvider,
		KProvider:  mockClient,
		KWriter:    mockWriter,
	}
	return serverTestable, func() {
		cleanup()
	}
}

// wire.go:

// This is in a separate common package
var ZCommonSet = wire.NewSet(
	NewServerConfig,
	NewServer, config.NewAppConfig, NewKafkaConfig, kafka.NewClient, wire.Bind(new(kafka.Logger), new(logger.Logger)), logger.NewLogger, NewTracer,
	NewDbConfig, db.NewProvider, NewHttpServiceConfig, http.NewClientProvider, wire.Bind(new(http.Logger), new(logger.Logger)), http.NewLeveledLogger,
)

var ZCommonMockSet = wire.NewSet(
	NewServerConfig,
	NewServer, config.NewAppConfig, NewKafkaConfig, logger.NewLogger, NewTracer,
	NewDbConfig,
	NewHttpServiceConfig, http.NewClientProvider, wire.Bind(new(http.Logger), new(logger.Logger)), http.NewLeveledLogger, mock_kafka.NewMockClient, mock_kafka.NewMockWriter, wire.Bind(new(kafka.Client), new(*mock_kafka.MockClient)), mock_db.NewMockProvider, wire.Bind(new(db.Provider), new(*mock_db.MockProvider)),
)

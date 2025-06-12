package config

import (
	"scalper/services"

	"go.uber.org/zap"
)

type Services struct {
	UserPhoneService *services.UserPhoneService
}

func InitServices(repos *Repositories, logger *zap.Logger) *Services {
	return &Services{
		UserPhoneService: services.NewUserPhoneService(repos.UserPhone, logger),
	}
}

package config

import (
	"context"
	"scalper/repositories"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDBConfig struct {
	URI      string
	Database string
}

type Repositories struct {
	UserPhone *repositories.UserPhoneRepository
}

func InitMongoDB(cfg MongoDBConfig, logger *zap.Logger) (*mongo.Client, *Repositories, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	repos := &Repositories{
		UserPhone: repositories.NewUserPhoneRepository(client, cfg.Database, logger),
	}

	if err := repos.UserPhone.EnsureIndexes(ctx); err != nil {
		return nil, nil, err
	}
	return client, repos, nil
}

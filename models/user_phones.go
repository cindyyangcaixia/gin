package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserPhone struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`
	SerialNumber string             `bson:"serial_number" json:"serial_number" validate:"required"`
	PhoneNumber  string             `bson:"phone_number" json:"phone_number" validate:"required,phone"`
	ParentId     primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty" validate:"omitempty"`
	// ParentLevel  int32              `bson:"parent_level,omitempty"`
	Level       int32              `bson:"level,omitempty" validate:"gte=0"`
	AccessToken string             `bson:"access_token,omitempty"`
	AmID        primitive.ObjectID `bson:"am_id" json:"am_id" validate:"required"`
	CreateTime  time.Time          `bson:"create_time"`
	UpdateTime  time.Time          `bson:"update_time"`
}

func (UserPhone) CollectionName() string {
	return "userphones"
}

func EnsureIndexes(ctx context.Context, coll *mongo.Collection) error {
	serialNumberIndex := mongo.IndexModel{
		Keys:    bson.M{"serial_number": 1},
		Options: options.Index().SetUnique(true),
	}

	phoneNumberIndex := mongo.IndexModel{
		Keys:    bson.M{"phone_number": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := coll.Indexes().CreateMany(ctx, []mongo.IndexModel{serialNumberIndex, phoneNumberIndex})
	if err != nil {
		return fmt.Errorf("failed to create indexes %w", err)
	}
	return nil
}

package repositories

import (
	"context"
	"time"

	"scalper/models"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type UserPhoneRepository struct {
	coll   *mongo.Collection
	logger *zap.Logger
}

func NewUserPhoneRepository(client *mongo.Client, dbName string, logger *zap.Logger) *UserPhoneRepository {
	return &UserPhoneRepository{
		coll:   client.Database(dbName).Collection("userphones"),
		logger: logger,
	}
}

func (r *UserPhoneRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{"serial_number": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"phone_number": 1}, Options: options.Index().SetUnique(true)},
	})
	return err
}

func (r *UserPhoneRepository) InsertOne(ctx context.Context, userPhone *models.UserPhone) (*mongo.InsertOneResult,
	error) {
	userPhone.CreateTime = time.Now()
	userPhone.UpdateTime = time.Now()
	data, err := r.coll.InsertOne(ctx, userPhone)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return data, nil
}

func (r *UserPhoneRepository) FindOne(ctx context.Context, phoneNumber string) (*models.UserPhone, error) {
	var result models.UserPhone
	err := r.coll.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode((&result))
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &result, nil
}

func (r *UserPhoneRepository) ListUserPhones(ctx context.Context, phoneNumber string, serialNumber *string, page,
	limit int64) ([]*models.UserPhone, int64, error) {
	filter := bson.M{}
	if phoneNumber != "" {
		filter["phone_number"] = phoneNumber
	}
	if serialNumber != nil && *serialNumber != "" {
		filter["serial_number"] = *serialNumber
	}

	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)

	cursor, err := r.coll.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, errors.New(err.Error())
	}
	defer cursor.Close(ctx)

	var results []*models.UserPhone
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, errors.New(err.Error())
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, errors.New(err.Error())
	}

	return results, total, nil
}

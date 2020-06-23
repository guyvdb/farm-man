package mongo

import (
	"context"
	"fmt"
	"github.com/guyvdb/farm-man/platform/model/infrastructure"
	"github.com/guyvdb/farm-man/platform/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	//	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInfrastructureRepository struct {
	database *mongo.Database
}

func NewMongoInfrastructureRepository(database *mongo.Database) repository.InfrastructureRepository {
	return &MongoInfrastructureRepository{
		database: database,
	}
}

func (r *MongoInfrastructureRepository) FindBoundsById(id string) (*infrastructure.Bounds, error) {
	var result infrastructure.Bounds

	filter := bson.M{"_id": id}
	col := r.database.Collection("bounds")

	retval := col.FindOne(context.Background(), filter)
	err := retval.Decode(&result)

	//err := r.database.Collection("bounds").FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *MongoInfrastructureRepository) ListAllBounds() ([]*infrastructure.Bounds, error) {
	return []*infrastructure.Bounds{}, nil
}

func (r *MongoInfrastructureRepository) InsertBounds(bounds *infrastructure.Bounds) error {
	_, err := r.database.Collection("bounds").InsertOne(context.Background(), bounds)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoInfrastructureRepository) UpdateBounds(bounds *infrastructure.Bounds) error {
	return nil
}

func (r *MongoInfrastructureRepository) DeleteBounds(bounds *infrastructure.Bounds) error {
	return nil
}

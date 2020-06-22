package mongo

import (
	"github.com/guyvdb/farm-man/platform/adapter"
	"github.com/guyvdb/farm-man/platform/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
)

type MongoAdapter struct {
	client   *mongo.Client
	database string
}

func NewMongoAdapter(uri string, database string) (adapter.Adapter, error) {
	client, err := mongoConnect(uri)
	if err != nil {
		return nil, err
	}
	return &MongoAdapter{client: client, database: database}, nil
}

func mongoConnect(uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (a *MongoAdapter) GetInfrastructureRepository() repository.InfrastructureRepository {
	return NewMongoInfrastructureRepository()
}

func (a *MongoAdapter) GetSequenceRepository() repository.SequenceRepository {
	db := a.client.Database(a.database)
	return NewMongoSequenceRepository(db)
}

package mongo

import (
	"github.com/guyvdb/farm-man/platform/model/infrastructure"
	"github.com/guyvdb/farm-man/platform/repository"
)

type MongoInfrastructureRepository struct {
}

func NewMongoInfrastructureRepository() repository.InfrastructureRepository {
	return &MongoInfrastructureRepository{}
}

func (r *MongoInfrastructureRepository) GetRootBounds() *infrastructure.Bounds {
	return nil
}

package adapter

import (
	"github.com/guyvdb/farm-man/platform/repository"
)

type Adapter interface {
	GetSequenceRepository() repository.SequenceRepository
	GetInfrastructureRepository() repository.InfrastructureRepository
}

package domain

// This is the domain service. It encorporates application specific rules that are not implemented
// in the model layer. Models have enterprise or business specific rules in them.

// It interfaces with the adapter.Adapter interface to retreive various repository interfaces. Code
// in this package may not interface with any database or framework specific components. Everything
// must be abstracted to adapter.Adapter interface, repository interfaces or the specific models

import (
	"github.com/guyvdb/farm-man/platform/model/infrastructure"
	"github.com/guyvdb/farm-man/platform/model/sequence"
)

type DomainService interface {
	SequenceDomainService
	InfrastructureDomainService
}

type SequenceDomainService interface {
	CreateSequence(prefix string, padding int) error
	DeleteSequence(prefix string) error
	ResetSequence(prefix string) error
	DeleteAllSequences() error
	NextSequence(prefix string, seperator string) sequence.Sequence
}

type InfrastructureDomainService interface {
	CreateBounds(btype infrastructure.BoundsType, name string, parent *infrastructure.Bounds) (*infrastructure.Bounds, error)
	UpdateBounds(bounds *infrastructure.Bounds) error
	DeleteBounds(bounds *infrastructure.Bounds) error

	CreateTank() (*infrastructure.Tank, error)
	UpdateTank(tank *infrastructure.Tank) error
	DeleteTank(tank *infrastructure.Tank) error
}

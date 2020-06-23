package repository

import (
	"github.com/guyvdb/farm-man/platform/model/infrastructure"
)

type InfrastructureRepository interface {

	// Bounds is an abstraction that segements the farm by area such as farm, tunnel, row, etc.
	// The root bounds is that root of a tree of all bounds.
	//	GetRootBounds() (*infrastructure.Bounds, error) // this should move to the domain service

	FindBoundsById(id string) (*infrastructure.Bounds, error)
	ListAllBounds() ([]*infrastructure.Bounds, error)
	InsertBounds(bounds *infrastructure.Bounds) error
	UpdateBounds(bounds *infrastructure.Bounds) error
	DeleteBounds(bounds *infrastructure.Bounds) error
}

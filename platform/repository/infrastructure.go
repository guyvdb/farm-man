package repository

import (
	"github.com/guyvdb/farm-man/platform/model/infrastructure"
)

type InfrastructureRepository interface {

	// Bounds is an abstraction that segements the farm by area such as farm, tunnel, row, etc.
	// The root bounds is that root of a tree of all bounds.
	GetRootBounds() *infrastructure.Bounds
}

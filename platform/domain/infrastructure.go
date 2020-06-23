package domain

import (
	"github.com/guyvdb/farm-man/platform/model/infrastructure"
)

/* ------------------------------------------------------------------------
 * Create a new bounds. Assign it a new ID and save it to the underlying
 * datastore
 * --------------------------------------------------------------------- */
func (s *Service) CreateBounds(btype infrastructure.BoundsType, name string, parent *infrastructure.Bounds) (*infrastructure.Bounds, error) {
	seq := s.adapter.GetSequenceRepository().Next(infrastructure.BoundsTypePrefix(btype), "-")
	bounds := infrastructure.NewBounds(seq, btype, name, parent)

	err := s.adapter.GetInfrastructureRepository().InsertBounds(bounds)
	if err != nil {
		return nil, err
	}
	return bounds, nil
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func (s *Service) UpdateBounds(bounds *infrastructure.Bounds) error {
	return nil
}

func (s *Service) DeleteBounds(bounds *infrastructure.Bounds) error {
	return nil
}

func (s *Service) CreateTank() (*infrastructure.Tank, error) {
	return nil, nil
}

func (s *Service) UpdateTank(tank *infrastructure.Tank) error {
	return nil
}

func (s *Service) DeleteTank(tank *infrastructure.Tank) error {
	return nil
}

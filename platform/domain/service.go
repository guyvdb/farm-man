package domain

import (
	"github.com/guyvdb/farm-man/platform/adapter"
)

type Service struct {
	adapter adapter.Adapter
}

func NewDomainService(adapter adapter.Adapter) DomainService {
	return &Service{
		adapter: adapter,
	}
}

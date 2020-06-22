package repository

import (
	"github.com/guyvdb/farm-man/platform/model/sequence"
)

type SequenceRepository interface {
	Next(prefix string, seperator string) sequence.Sequence
}

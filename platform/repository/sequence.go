package repository

import (
	"github.com/guyvdb/farm-man/platform/model/sequence"
)

type SequenceRepository interface {
	CreateSequence(prefix string, padding int) error
	DeleteSequence(prefix string) error
	ResetSequence(prefix string) error
	DeleteAllSequences() error
	Next(prefix string, seperator string) sequence.Sequence
}

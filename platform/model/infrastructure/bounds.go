package infrastructure

// This is an abstract model that represents a tree like containment of farm, growing areas, houses, bins, etc.
import (
	//	"fmt"
	//	"github.com/google/uuid"
	"github.com/guyvdb/farm-man/platform/model/sequence"
)

type BoundsType int

const (
	BOUNDSTYPE_FARM BoundsType = iota
	BOUNDSTYPE_AREA
	BOUNDSTYPE_GREENHOUSE
	BOUNDSTYPE_ROW
	BOUNDSTYPE_BIN
	//BOUNDSTYPE_RAFT
	//BOUNDSTYPE_BUCKET
)

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
type Bounds struct {
	Id       sequence.Sequence `bson:"_id" json:"id,omitempty"`
	Type     BoundsType        `bson:"typeid",json:"typeid"`
	ParentId sequence.Sequence `bson:"parentid", json:"parentid"`
	Name     string            `bson:"name",json:"name"`
	Children []*Bounds
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func NewBounds(seq sequence.Sequence, btype BoundsType, name string, parent *Bounds) *Bounds {
	if parent == nil {
		return &Bounds{
			Id:       seq,
			Type:     btype,
			Name:     name,
			Children: make([]*Bounds, 0, 8),
		}
	} else {
		return &Bounds{
			Id:       seq,
			Type:     btype,
			ParentId: parent.Id,
			Name:     name,
			Children: make([]*Bounds, 0, 8),
		}
	}
}

/* ------------------------------------------------------------------------
 * Set a list of children onto the container
 * --------------------------------------------------------------------- */
func (b *Bounds) SetChildren(children []*Bounds) {
	for _, child := range children {
		b.Children = append(b.Children, child)
	}
}

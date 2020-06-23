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

var boundsPrefixes map[BoundsType]string = map[BoundsType]string{
	BOUNDSTYPE_FARM:       "FARM",
	BOUNDSTYPE_AREA:       "AREA",
	BOUNDSTYPE_GREENHOUSE: "GREENHOUSE",
	BOUNDSTYPE_ROW:        "ROW",
	BOUNDSTYPE_BIN:        "BIN",
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
func BoundsTypePrefix(btype BoundsType) string {
	p, ok := boundsPrefixes[btype]
	if ok {
		return p
	}
	return "BOUNDS_TYPE_ERROR"
}

/* ------------------------------------------------------------------------
 *
 * --------------------------------------------------------------------- */
type Bounds struct {
	Id       sequence.Sequence `bson:"_id" json:"id,omitempty"`
	Type     BoundsType        `bson:"tid",json:"tid"`
	ParentId sequence.Sequence `bson:"pid,omitempty", json:"pid"`
	Name     string            `bson:"name",json:"name"`
	Children []*Bounds         `bson:"-"`
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

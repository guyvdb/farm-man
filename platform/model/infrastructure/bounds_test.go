package infrastructure

import (
	"github.com/google/uuid"

	//"fmt"
	"github.com/guyvdb/farm-man/platform/model/infrastructure"
	"testing"
)

func CreateRootBounds() *infrastructure.Bounds {
	return infrastructure.NewBounds(infrastructure.BOUNDSTYPE_FARM, "farm", nil)
}

// func NewBounds(ctype BoundsType, name string, parent *Bounds) *Bounds {

func TestBoundsNewWithoutParent(t *testing.T) {
	c := CreateRootBounds()

	if c.Name != "farm" {
		t.Errorf("Name error")
	}

	if c.Type != infrastructure.BOUNDSTYPE_FARM {
		t.Errorf("Type error")
	}
}

func TestBoundsNewWithParent(t *testing.T) {
	p := CreateRootBounds()
	c := infrastructure.NewBounds(infrastructure.BOUNDSTYPE_AREA, "field 1", p)
	if c.ParentId.String() != p.Id.String() {
		t.Errorf("Parent error")
	}
}

func TestBoundsNewWithIdWithoutParent(t *testing.T) {
	id := uuid.New()
	c := infrastructure.NewBoundsWithId(id, infrastructure.BOUNDSTYPE_AREA, "field 1", nil)
	if c.Id.String() != id.String() {
		t.Errorf("WithId error")
	}
}

func TestBoundsNewWithIdWithParent(t *testing.T) {
	id := uuid.New()
	p := CreateRootBounds()
	c := infrastructure.NewBoundsWithId(id, infrastructure.BOUNDSTYPE_AREA, "field 1", p)

	if c.Id.String() != id.String() {
		t.Errorf("WithId id error")
	}

	if c.ParentId.String() != p.Id.String() {
		t.Errorf("WithId parent error")
	}

}

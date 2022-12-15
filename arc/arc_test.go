/******************************************************************************
Filename: arc.go
Names: Arsh Banerjee, Kenny Lam, and Kenar Vyas
NetId:
Description:
*****************************************************************************/

package arc

import (
	"testing"
)

// Helper Functions:

func TestArcInit(t *testing.T) {
	arc := NewArc(8)
	if arc == nil {
		t.Errorf("Error initializing ARC")
	}
	if arc.RemainingStorage() != 8 {
		t.Errorf("Size of ARC cache is %d bytes, should be 8 bytes", arc.RemainingStorage())
	}
	if arc.Len() != 0 {
		t.Errorf("Bindings in ARC cache is %d, should be 0", arc.Len())
	}
}

func TestTooSmallArcInit(t *testing.T) {
	if NewArc(0) != nil || NewArc(1) != nil {
		t.Errorf("initializing an ARC object even when size is incorrect")
	}
}

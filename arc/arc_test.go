/******************************************************************************
Filename: arc.go
Names: Arsh Banerjee, Kenny Lam, and Kenar Vyas
NetId:
Description:
*****************************************************************************/

package arc

import (
	"log"
	"strconv"
	"testing"
)

// Helper Functions:
// Checks to see if Get exists and returns correct value
func confirmGet(arc *ARC, val string) {
	Val, err := arc.Get(val)
	if err != true {
		log.Fatalf("Key: %s not found", val)
	}
	Value := string(Val)
	if Value != val {
		log.Fatalf("Get of key: %s is %s, should be %s", val, Value, val)
	}
}

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

func TestAddingUntilCapacity(t *testing.T) {
	arc := NewArc(16)
	if arc == nil {
		t.Errorf("Error initializing ARC")
	}

	values := []int{1, 2, 3, 4, 5, 6, 7, 8}

	for i := 0; i < len(values); i++ {
		val := strconv.Itoa(i)
		arc.Set(val, []byte(val))
	}
}

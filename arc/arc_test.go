/******************************************************************************
Filename: arc.go
Names: Arsh Banerjee, Kenny Lam, and Kenar Vyas
NetId:
Description:
*****************************************************************************/

package arc

import (
	"testing"
	"fmt"
)

// Helper Functions:

func TestArcInit(t *testing.T) {
	fmt.Println("Testing init..")
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
	fmt.Println("Testing too small..")

	if NewArc(0) != nil || NewArc(1) != nil {
		t.Errorf("initializing an ARC object even when size is incorrect")
	}
}

func TestBasicSetGet(t *testing.T) {
	fmt.Println("Testing basic set get..")

	arc := NewArc(8)

	if arc == nil {
		t.Errorf("Error initializing ARC")
	}

	arc.Set("a", []byte{'a'})

	if CacheToStr(arc.cacheOrder[:arc.splitIndex+1]) != "a x x x " {
		t.Errorf("Error set")
	}

	if CacheToStr(arc.cacheOrder[arc.splitIndex+1:]) != "x x x x " {
		t.Errorf("Error set")
	}

	arc.Set("b", []byte{'b'})
	arc.Set("c", []byte{'c'})
	arc.Set("d", []byte{'d'})

	if CacheToStr(arc.cacheOrder[:arc.splitIndex+1]) != "d c b a " {
		fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
		t.Errorf("Error set")
	}

	if CacheToStr(arc.cacheOrder[arc.splitIndex+1:]) != "x x x x " {
		t.Errorf("Error set")
	}

	arc.Set("b", []byte{'b'})

	if CacheToStr(arc.cacheOrder[:arc.splitIndex+1]) != "d c a " {
		fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
		t.Errorf("Error set")
	}

	if CacheToStr(arc.cacheOrder[arc.splitIndex+1:]) != "x x x x b " {
		t.Errorf("Error set")
	}

	arc.Get("a") 

	if CacheToStr(arc.cacheOrder[:arc.splitIndex+1]) != "d c " {
		fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
		t.Errorf("Error get")
	}

	if CacheToStr(arc.cacheOrder[arc.splitIndex+1:]) != "x x x x b a " {
		fmt.Println(CacheToStr(arc.cacheOrder[arc.splitIndex+1:]))		
		t.Errorf("Error get")
	}	

	arc.Get("b")

	if CacheToStr(arc.cacheOrder[:arc.splitIndex+1]) != "d c " {
		fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
		t.Errorf("Error get")
	}

	if CacheToStr(arc.cacheOrder[arc.splitIndex+1:]) != "x x x x a b " {
		fmt.Println(CacheToStr(arc.cacheOrder[arc.splitIndex+1:]))		
		t.Errorf("Error get")
	}	

	arc.Get("c")

	if CacheToStr(arc.cacheOrder[:arc.splitIndex+1]) != "d " {
		fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
		t.Errorf("Error get")
	}

	if CacheToStr(arc.cacheOrder[arc.splitIndex+1:]) != "x x x x a b c " {
		fmt.Println(CacheToStr(arc.cacheOrder[arc.splitIndex+1:]))		
		t.Errorf("Error get")
	}	

}

func TestGhostBuffer(t *testing.T) {

	fmt.Println("Testing ghost buffer..")

	arc := NewArc(8)
	arc.Set("a", []byte{'a'})
	arc.Set("b", []byte{'a'})
	arc.Set("c", []byte{'a'})
	arc.Set("d", []byte{'a'})

	if CacheToStr(arc.cacheOrder[:arc.splitIndex+1]) != "d c b a " {
		fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
		t.Errorf("Error set")
	}

	if CacheToStr(arc.cacheOrder[arc.splitIndex+1:]) != "x x x x " {
		t.Errorf("Error set")
	}

	arc.Set("e", []byte{'e'})

	if CacheToStr(arc.cacheOrder[:arc.splitIndex+1]) != "e d c " {
		fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
		t.Errorf("Error set")
	}

	if CacheToStr(arc.b1CacheOrder) != "b a x x x x x x " {
		fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
		t.Errorf("Error set")
	}

	_, ok := arc.Get("b")

	if ok {
		t.Errorf("Error Get")
	}

	if CacheToStr(arc.cacheOrder[:arc.splitIndex+1]) != "e d c x " {
		fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
		t.Errorf("Error set")
	}

	arc.Set("f", []byte{'f'})
	arc.Get("a")
	arc.Get("c")
	if CacheToStr(arc.cacheOrder[:arc.splitIndex+1]) != "f e d x " {
		fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
		t.Errorf("Error set")
	}
	if CacheToStr(arc.cacheOrder[arc.splitIndex+1:]) != "x x x c " {
		fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
		t.Errorf("Error set")
	}
}

func CacheToStr(order []string) string {
	st := ""
	for i := 0; i < len(order); i++ {
		if order[i] == "" {
			st += "x "
		} else {
			st += order[i] + " "
		}
	}
	return st
}

func PrintAll(arc * ARC) {
	fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
	fmt.Println(CacheToStr(arc.cacheOrder[arc.splitIndex+1:]))		
	fmt.Println(CacheToStr(arc.b1CacheOrder))		
	fmt.Println(CacheToStr(arc.b2CacheOrder))		
}
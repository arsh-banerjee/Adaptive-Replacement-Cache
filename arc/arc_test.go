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
	"fmt"
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

// Turns cache into readable string with x demarcating empty slots
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

// Prints out all caches for easy debugging
func PrintAll(arc * ARC) {
	fmt.Println(CacheToStr(arc.cacheOrder[:arc.splitIndex+1]))		
	fmt.Println(CacheToStr(arc.cacheOrder[arc.splitIndex+1:]))		
	fmt.Println(CacheToStr(arc.b1CacheOrder))		
	fmt.Println(CacheToStr(arc.b2CacheOrder))		
}

// Tests:

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
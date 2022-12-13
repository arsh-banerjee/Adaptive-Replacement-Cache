/******************************************************************************
Filename: arc.go
Names: Arsh Banerjee, Kenny Lam, and Kenar Vyas
NetId: arshb
Description:
*****************************************************************************/

package arc

import (
	"sync"
)

type ARC struct {
	limit        int
	currentUsage int
	T            map[string]*entry
	splitIndex   int //Index that divides T into t1 and t2
	b1           map[string]string
	b2           map[string]string
	cacheOrder   []string // represents the order of cache entries
	lock         *sync.Mutex
}

type entry struct {
	Key   string
	Value []byte
}

// NewARC returns a pointer to a new ARC with a capacity to store limit bytes
func NewArc(limit int) *ARC {
	return &ARC{
		limit: limit, lock: new(sync.Mutex),
		currentUsage: 0,
		T:            make(map[string]*entry, limit),
		splitIndex:   int(limit/2) - 1,
		b1:           make(map[string]string),
		b2:           make(map[string]string),
		cacheOrder:   make([]string, limit),
		lock:         new(sync.Mutex),
	}
}

// MaxStorage returns the maximum number of bytes this ARC can store
func (arc *ARC) MaxStorage() int {
	return arc.limit
}

// RemainingStorage returns the number of unused bytes available in this ARC
func (arc *ARC) RemainingStorage() int {
	arc.lock.Lock()
	defer arc.lock.Unlock()

	return arc.limit - arc.currentUsage
}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (arc *ARC) Get(key string) (value []byte, ok bool) {

	val, prs := arc.T[key]

	if prs {
		var index int
		for i := 0; i < arc.limit; i++ {
			if key == arc.cacheOrder[i] {
				index = i
				break
			}
		}
		// if in LRU portion of cache
		if index < arc.splitIndex+1 {

		} else {
			// if in LFU portion of cache, move the cache entry to the front of the LFU list
			newOrder := make([]string, arc.splitIndex)
			firstHalf := arc.cacheOrder[arc.splitIndex+1 : index]
			secondHalf := arc.cacheOrder[index:]
			val := arc.cacheOrder[index]
			copy(newOrder, arc.cacheOrder[0:arc.splitIndex+1])
			newOrder = append(newOrder, val)
			newOrder = append(newOrder, firstHalf...)
			newOrder = append(newOrder, secondHalf...)
			arc.cacheOrder = newOrder

		}
		return val.Value, true
	}

	_, prsB1 := arc.b1[key]
	if prsB1 {

		return nil, false

	}

	_, prsB2 := arc.b2[key]
	if prsB2 {

		return nil, false
	}

	return nil, false
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (arc *ARC) Remove(key string) (value []byte, ok bool) {

	return nil, false
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (arc *ARC) Set(key string, value []byte) bool {
	arc.lock.Lock()
	defer arc.lock.Unlock()

	return false
}

// Len returns the number of bindings in the ARC.
func (arc *ARC) Len() int {
	return arc.currentUsage
}

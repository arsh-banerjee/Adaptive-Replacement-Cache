/******************************************************************************
Filename: arc.go
Names: Arsh Banerjee, Kenny Lam, and Kenar Vyas
NetId: arshb
Description:
*****************************************************************************/

package arc

import (
	"log"
	"sync"
)

type ARC struct {
	limit        int
	currentUsage int
	len          int
	T            map[string]*entry
	splitIndex   int //Index that divides T into t1 and t2. cacheOrder[splitIndex] is the last element in T1

	// TODO: ensure that b1/b2 do not exceed a size of limit, and update
	// all functions to check for the size of b1/b2 before adding new entries
	b1         map[string]string
	b2         map[string]string
	cacheOrder []string // represents the order of cache entries.
	// cacheOrder[0] represents the most recent cache entry,
	// while cacheOrder[arc.limit - 1] represents the most frequently used
	lock *sync.Mutex
}

type entry struct {
	Key   string
	Value []byte
}

// NewARC returns a pointer to a new ARC with a capacity to store limit bytes
func NewArc(limit int) *ARC {
	return &ARC{
		limit:        limit,
		lock:         new(sync.Mutex),
		currentUsage: 0,
		T:            make(map[string]*entry, limit),
		splitIndex:   int(limit/2) - 1,
		b1:           make(map[string]string),
		b2:           make(map[string]string),
		cacheOrder:   make([]string, limit),
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
			val := arc.cacheOrder[index]
			temp := val
			for i := index; i < arc.limit-1; i++ {
				innerTemp := arc.cacheOrder[i]
				arc.cacheOrder[i] = temp
				temp = innerTemp
			}
			arc.cacheOrder[arc.limit-1] = val

		}
		return val.Value, true
	}

	_, prsB1 := arc.b1[key]
	if prsB1 {
		// evict a key from T2, expand T1, and put the key into B2
		evictedKey := arc.cacheOrder[arc.splitIndex+1]
		if evictedKey != "" {
			arc.b2[evictedKey] = evictedKey
			arc.cacheOrder[arc.splitIndex+1] = ""
			delete(arc.T, evictedKey)
		}
		arc.splitIndex++

		return nil, false

	}

	_, prsB2 := arc.b2[key]
	if prsB2 {
		evictedKey := arc.cacheOrder[arc.splitIndex]
		if evictedKey != "" {
			arc.b1[evictedKey] = evictedKey
			arc.cacheOrder[arc.splitIndex] = ""
			delete(arc.T, evictedKey)
		}
		arc.splitIndex--

		return nil, false
	}

	return nil, false
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (arc *ARC) Remove(key string) (value []byte, ok bool) {

	val, ok := arc.T[key]

	if !ok {
		return nil, false
	}

	index := arc.GetIndex(key)

	if index == -1 {
		log.Fatalf("key not found")
	} else {
		arc.cacheOrder = append(arc.cacheOrder[:index], arc.cacheOrder[index+1:]...)

		// TODO: make sure this is consistent with Get
		if index < arc.splitIndex+1 {
			arc.splitIndex -= 1
		} else {
			arc.splitIndex += 1
		}
	}

	delete(arc.T, key)

	// Update b1 / b2
	_, ok = arc.b1[key]
	if ok {
		delete(arc.b1, key)
	}
	_, ok = arc.b1[key]

	if ok {
		delete(arc.b2, key)
	}

	arc.currentUsage -= (len(val.Value) + len(val.Key))
	arc.len -= 1

	return val.Value, true
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (arc *ARC) Set(key string, value []byte) bool {
	arc.lock.Lock()
	defer arc.lock.Unlock()

	if len(value)+len(key) > arc.limit {
		return false
	}

	val, ok := arc.T[key]

	if !ok {
		arc.currentUsage += len(key) + len(value)
		arc.len += 1

		for arc.RemainingStorage() < 0 {
			evict_key := arc.cacheOrder[0] // Evicting from L1, least recently used
			_, ok := arc.Remove(evict_key)
			if !ok {
				log.Fatalf("Remove failed in Set")
			}
		}
		arc.T[key] = &entry{Key: key, Value: value}
		// TODO: update cache order and split index - confirm this is correct after get is finished
		arc.cacheOrder = insert(arc.cacheOrder, arc.splitIndex-1, key)
		arc.splitIndex += 1
		arc.Get(key) // mark as used

	} else {
		if arc.RemainingStorage()+len(val.Value)-len(value) < 0 {
			return false
		}

		arc.T[key].Value = value
		arc.currentUsage = (arc.currentUsage - len(val.Value)) + len(value)
		arc.Get(key) // mark as used
	}

	return true
}

// Len returns the number of bindings in the ARC.
func (arc *ARC) Len() int {
	return arc.len
}

func (arc *ARC) GetIndex(key string) int {
	for i := 0; i < arc.limit; i++ {
		if key == arc.cacheOrder[i] {
			return i
		}
	}
	return -1
}

func insert(a []string, index int, value string) []string {
	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

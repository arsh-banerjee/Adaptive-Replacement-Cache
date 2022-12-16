/******************************************************************************
Filename: arc.go
Names: Arsh Banerjee, Kenny Lam, and Kenar Vyas
NetId: arshb
Description:
*****************************************************************************/

package arc

import (
	"fmt"
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

	// represents the order of the entries b1/b2 cache. b#CacheOrder[0] represents the most recently added
	b1CacheOrder []string
	b2CacheOrder []string
	b1Size       int
	b2Size       int
	lock         *sync.Mutex
}

type entry struct {
	Key   string
	Value []byte
}

// NewARC returns a pointer to a new ARC with a capacity to store limit bytes
func NewArc(limit int) *ARC {

	if limit < 2 {
		return nil
	}
	return &ARC{
		limit:        limit,
		lock:         new(sync.Mutex),
		len:          0,
		currentUsage: 0,
		T:            make(map[string]*entry, limit),
		splitIndex:   int(limit/2) - 1,
		b1:           make(map[string]string),
		b2:           make(map[string]string),
		cacheOrder:   make([]string, limit),
		b1CacheOrder: make([]string, limit),
		b1Size:       0,
		b2CacheOrder: make([]string, limit),
		b2Size:       0,
	}
}

// MaxStorage returns the maximum number of bytes this ARC can store
func (arc *ARC) MaxStorage() int {
	return arc.limit
}

// RemainingStorage returns the number of unused bytes available in this ARC
func (arc *ARC) RemainingStorage() int {
	// arc.lock.Lock()
	// defer arc.lock.Unlock()

	return arc.limit - arc.currentUsage
}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (arc *ARC) Get(key string) (value []byte, ok bool) {
	// arc.lock.Lock()
	// defer arc.lock.Unlock()
	fmt.Printf("")
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
		// move key to LFU portion of cache, shift everything to the left
		if index < arc.splitIndex+1 {
			keyVal := arc.cacheOrder[index]
			for i := index; i < len(arc.cacheOrder)-1; i++ {
				arc.cacheOrder[i] = arc.cacheOrder[i+1]
			}
			arc.cacheOrder[len(arc.cacheOrder)-1] = keyVal
			arc.splitIndex--

			// if in LFU portion of cache, move key to front of LFU by shifting everything left
		} else {
			keyVal := arc.cacheOrder[index]
			for i := index; i < len(arc.cacheOrder)-1; i++ {
				arc.cacheOrder[i] = arc.cacheOrder[i+1]
			}
			arc.cacheOrder[len(arc.cacheOrder)-1] = keyVal
		}
		return val.Value, true
	}

	_, prsB1 := arc.b1[key]
	if prsB1 {
		// evict a key from T2, expand T1, and put the key into B2
		evictedKey := arc.cacheOrder[arc.splitIndex+1]
		if evictedKey != "" {
			arc.cacheOrder[arc.splitIndex+1] = ""
			arc.currentUsage -= (len(evictedKey) + len(arc.T[evictedKey].Value))
			arc.len -= 1
			delete(arc.T, evictedKey)

			if arc.b2Size >= arc.limit {
				b2EvictedKey := arc.b2CacheOrder[len(arc.b2CacheOrder)-1]
				delete(arc.b2, b2EvictedKey)
				arc.b2Size--
			}
			for i := len(arc.b2CacheOrder) - 1; i > 0; i-- {
				arc.b2CacheOrder[i] = arc.b2CacheOrder[i-1]
			}
			arc.b2CacheOrder[0] = evictedKey
			arc.b2[evictedKey] = evictedKey
			arc.b2Size++
		}
		arc.splitIndex++

		return nil, false

	}

	_, prsB2 := arc.b2[key]
	if prsB2 {
		// evict a key from T1, expand T2, and put the key into B1
		evictedKey := arc.cacheOrder[arc.splitIndex]
		if evictedKey != "" {
			arc.cacheOrder[arc.splitIndex] = ""
			arc.currentUsage -= len(evictedKey) + len(arc.T[evictedKey].Value)
			arc.len -= 1
			delete(arc.T, evictedKey)

			if arc.b1Size >= arc.limit {
				b1EvictedKey := arc.b1CacheOrder[len(arc.b1CacheOrder)-1]
				delete(arc.b1, b1EvictedKey)
				arc.b1Size--
			}
			for i := len(arc.b1CacheOrder) - 1; i > 0; i-- {
				arc.b1CacheOrder[i] = arc.b1CacheOrder[i-1]
			}
			arc.b1CacheOrder[0] = evictedKey
			arc.b1[evictedKey] = evictedKey
			arc.b1Size++
		}
		arc.splitIndex--

		return nil, false
	}

	return nil, false
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (arc *ARC) Remove(key string) (value []byte, ok bool) {
	// arc.lock.Lock()
	// defer arc.lock.Unlock()
	val, ok := arc.T[key]

	if !ok {
		return nil, false
	}

	index := arc.GetIndex(key)
	evictedKey := arc.cacheOrder[index]

	if index == -1 {
		log.Fatalf("key not found")
	} else {
		// if in the LRU part of the cache
		if index < arc.splitIndex+1 {
			for i := index; i < arc.splitIndex; i++ {
				arc.cacheOrder[i] = arc.cacheOrder[i+1]
			}
			arc.cacheOrder[arc.splitIndex] = ""

			// add to b1
			if arc.b1Size >= arc.limit {
				b1EvictedKey := arc.b1CacheOrder[len(arc.b1CacheOrder)-1]
				delete(arc.b1, b1EvictedKey)
				arc.b1Size--
			}
			for i := len(arc.b1CacheOrder) - 1; i > 0; i-- {
				arc.b1CacheOrder[i] = arc.b1CacheOrder[i-1]
			}
			arc.b1CacheOrder[0] = evictedKey
			arc.b1[evictedKey] = evictedKey
			arc.b1Size++

			// if in the LFU part of the cache
		} else {
			for i := index; i > arc.splitIndex+1; i-- {
				arc.cacheOrder[i] = arc.cacheOrder[i-1]
			}
			arc.cacheOrder[arc.splitIndex+1] = ""

			if arc.b2Size >= arc.limit {
				b2EvictedKey := arc.b2CacheOrder[len(arc.b2CacheOrder)-1]
				delete(arc.b2, b2EvictedKey)
				arc.b2Size--
			}
			for i := len(arc.b2CacheOrder) - 1; i > 0; i-- {
				arc.b2CacheOrder[i] = arc.b2CacheOrder[i-1]
			}
			arc.b2CacheOrder[0] = evictedKey
			arc.b2[evictedKey] = evictedKey
			arc.b2Size++
		}
	}
	delete(arc.T, key)

	arc.currentUsage -= (len(val.Value) + len(val.Key))
	arc.len -= 1

	return val.Value, true
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (arc *ARC) Set(key string, value []byte) bool {
	// arc.lock.Lock()
	// defer arc.lock.Unlock()

	if len(value)+len(key) > arc.limit {
		return false
	}

	val, ok := arc.T[key]

	if !ok {
		arc.currentUsage += len(key) + len(value)
		arc.len += 1

		switchList := false //boolean to switch between empyting LRU and LFU

		for arc.RemainingStorage() < 0 {
			k := arc.splitIndex
			var evict_key string
			for !switchList {
				evict_key = arc.cacheOrder[k] // Evicting from L1, LRU
				if evict_key != "" {
					break
				}
				k--
				if k < 0 {
					switchList = true
					k = arc.splitIndex
					break
				}
			}

			for switchList {
				evict_key = arc.cacheOrder[k+1]
				if evict_key != "" {
					break
				}
				k++
			}

			_, ok := arc.Remove(evict_key)
			if !ok {
				log.Fatalf("Remove failed in Set")
			}
		}

		arc.T[key] = &entry{Key: key, Value: value}

		temp := key
		for i := 0; i <= arc.splitIndex; i++ {
			if arc.cacheOrder[i] == "" {
				arc.cacheOrder[i] = temp
				break
			}
			inner_temp := arc.cacheOrder[i]
			arc.cacheOrder[i] = temp
			temp = inner_temp

			if i == arc.splitIndex && temp != "" {
				if arc.b1Size >= arc.limit {
					b1EvictedKey := arc.b1CacheOrder[len(arc.b1CacheOrder)-1]
					delete(arc.b1, b1EvictedKey)
					arc.b1Size--
				}
				for i := len(arc.b1CacheOrder) - 1; i > 0; i-- {
					arc.b1CacheOrder[i] = arc.b1CacheOrder[i-1]
				}
				arc.b1CacheOrder[0] = temp
				arc.b1[temp] = temp
				arc.b1Size++

				arc.currentUsage -= len(temp) + len(arc.T[temp].Value)
				arc.len -= 1

				delete(arc.T, temp)
			}
		}
		arc.cacheOrder[0] = key

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

package arc

import (
	"container/list"
	"fmt"
)

// An LRU is a fixed-size in-memory cache with least-recently-used eviction
type LRU struct {
	limit        int
	currentUsage int
	pointers     list.List
	// whatever fields you want here
}
type KeyPairs struct {
	Key   string
	Value []byte
}

// NewLRU returns a pointer to a new LRU with a capacity to store limit bytes
func NewLru(limit int) *LRU {
	return &LRU{limit: limit, currentUsage: 0}
}

// MaxStorage returns the maximum number of bytes this LRU can store
func (lru *LRU) MaxStorage() int {
	return lru.limit
}

// RemainingStorage returns the number of unused bytes available in this LRU
func (lru *LRU) RemainingStorage() int {
	return lru.limit - lru.currentUsage
}

func (lru *LRU) listIterate() {
	node := lru.pointers.Front()
	fmt.Println(node.Value.(KeyPairs).Key + string(node.Value.(KeyPairs).Value))

	for node.Next() != nil {
		fmt.Println(node.Next().Value.(KeyPairs).Key + string(node.Next().Value.(KeyPairs).Value))
		node = node.Next()
	}
}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lru *LRU) Get(key string) (value []byte, ok bool) {
	if lru.pointers.Len() == 0 {
		return nil, false
	}

	node := lru.pointers.Front()
	tempValue := node.Value

	if tempValue.(KeyPairs).Key == key {
		return tempValue.(KeyPairs).Value, true
	}

	for node.Next() != nil {
		tempValue := node.Next().Value
		if tempValue.(KeyPairs).Key == key {
			lru.pointers.MoveToFront(node.Next())
			return tempValue.(KeyPairs).Value, true
		}
		node = node.Next()
	}

	return nil, false
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lru *LRU) Remove(key string) (value []byte, ok bool) {
	if lru.pointers.Len() == 0 {
		return nil, false
	}

	node := lru.pointers.Front()
	if node.Value.(KeyPairs).Key == key {
		removeBytes := len(node.Value.(KeyPairs).Key) + len(node.Value.(KeyPairs).Value)
		lru.currentUsage -= removeBytes
		lru.pointers.MoveToBack(node)
		return lru.pointers.Remove(node).(KeyPairs).Value, true
	}

	for node.Next() != nil {
		if node.Next().Value.(KeyPairs).Key == key {
			removeBytes := len(node.Next().Value.(KeyPairs).Key) + len(node.Next().Value.(KeyPairs).Value)
			lru.currentUsage -= removeBytes
			lru.pointers.MoveToBack(node.Next())
			return lru.pointers.Remove(node.Next()).(KeyPairs).Value, true
		}
		node = node.Next()
	}

	return nil, false
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lru *LRU) Set(key string, value []byte) bool {
	byteSize := (len(key) + len(value))
	if byteSize > lru.limit {
		return false
	}

	found, _ := lru.Get(key)
	if found != nil {
		node := lru.pointers.Front()
		if node.Value.(KeyPairs).Key == key {
			lru.currentUsage += (len(value) - len(node.Value.(KeyPairs).Value))
			node.Value = KeyPairs{Key: key, Value: value}
			return true
		}

		for node.Next() != nil {
			if node.Next().Value.(KeyPairs).Key == key {
				lru.currentUsage += (len(value) - len(node.Next().Value.(KeyPairs).Value))
				node.Next().Value = KeyPairs{Key: key, Value: value}
				lru.pointers.MoveToFront(node.Next())
				return true
			}
			node = node.Next()
		}
	}

	if byteSize > lru.RemainingStorage() {
		for byteSize > lru.RemainingStorage() {
			removed := len(lru.pointers.Back().Value.(KeyPairs).Key) + len(lru.pointers.Back().Value.(KeyPairs).Value)
			lru.currentUsage -= removed
			lru.pointers.Remove(lru.pointers.Back())
		}
	}

	lru.pointers.PushFront(KeyPairs{Key: key, Value: value})
	lru.currentUsage += byteSize
	return true
}

// Len returns the number of bindings in the LRU.
func (lru *LRU) Len() int {
	return lru.pointers.Len()
}

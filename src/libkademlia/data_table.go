/*
These functions except Init and Finalize are thread safe
*/

package libkademlia

import (
	"errors"
	"sync"
)

// HashTable :
type DataTable struct {
	Table  map[ID]VanashingDataObject
	Parent *Kademlia
	Mutex  sync.Mutex
}

// Init : Not thread safe, should be called only once. Must be called before all other functions can work
func (tab *DataTable) Init(Parent *Kademlia) error {
	tab.Table = make(map[ID]VanashingDataObject)
	tab.Parent = Parent
	return nil
}

// Find :
func (tab *DataTable) Find(key ID) (V VanashingDataObject, err error) {
	tab.Mutex.Lock()
	v, ok := tab.Table[key]
	tab.Mutex.Unlock()
	if ok {
		return v, nil
	}
	return V, errors.New("Key not found")
}

// Add :
func (tab *DataTable) Add(key ID, V VanashingDataObject) error {
	tab.Mutex.Lock()
	_, ok := tab.Table[key]
	tab.Table[key] = V
	tab.Mutex.Unlock()
	if ok {
		return errors.New("Already in table")
	}
	return nil
}

// Remove :
func (tab *DataTable) Remove(key ID, V VanashingDataObject) error {
	tab.Mutex.Lock()
	_, ok := tab.Table[key]
	if ok {
		delete(tab.Table, key)
	}
	tab.Mutex.Unlock()
	if !ok {
		return errors.New("Not in table")
	}
	return nil
}

/*
These functions except Init and Finalize are thread safe
*/

package libkademlia

import (
	"errors"
	"sync"
	"time"
)

// HashTable :
type DataTable struct {
	Table  map[ID]VanashingDataObject
	Expire map[ID]time.Time
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
	exp, eok := tab.Expire[key]
	tab.Mutex.Unlock()
	if ok {
		if eok {
			if time.Now().After(exp) {
				tab.Remove(key)
				return V, errors.New("Object expired")
			}
		}
		return v, nil
	}
	return V, errors.New("Key not found")
}

// Add :
func (tab *DataTable) Add(key ID, V VanashingDataObject) error {
	return tab.AddEx(key, V, -1)
}

// AddEx : Add with expiration time
func (tab *DataTable) AddEx(key ID, V VanashingDataObject, exp_sec int64) error {
	tab.Mutex.Lock()
	_, ok := tab.Table[key]
	tab.Table[key] = V
	if exp_sec > 0 {
		tab.Expire[key] = time.Now().Add(time.Duration(exp_sec * 1000000000))
	}
	tab.Mutex.Unlock()
	if ok {
		return errors.New("Already in table")
	}
	return nil
}

// Remove :
func (tab *DataTable) Remove(key ID) error {
	tab.Mutex.Lock()
	_, ok := tab.Table[key]
	if ok {
		_, ok = tab.Expire[key]
		if ok {
			delete(tab.Expire, key)
		}
		delete(tab.Table, key)
	}
	tab.Mutex.Unlock()
	if !ok {
		return errors.New("Not in table")
	}
	return nil
}

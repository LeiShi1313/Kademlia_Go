/*
These functions are not thread safe, should not be called by other routine
*/

package libkademlia

import (
	"errors"
	"fmt"
)

const (
	HASH_TABLE_EVENT_ADD    = 1
	HASH_TABLE_EVENT_FIND   = 2
	HASH_TABLE_EVENT_REMOVE = 3
)

// HashTable :
type HashTable struct {
	Table     map[ID]HashTableEntry
	Self      *Kademlia
	EventChan chan HashTableEvent
}

// HashTableEntry :
type HashTableEntry struct {
	Key   ID
	Value []byte
}

// HashTableEvent :
type HashTableEvent struct {
	EventID int
	Arg     HashTableEventArg
	Ret     chan error
}

// HashTableEventArg :
type HashTableEventArg struct {
	Key   *ID
	Value **[]byte
}

// Dispatcher :
func (tab *HashTable) Dispatcher() {
	var Event HashTableEvent
	var Ret error
	running := true
	for running {
		event, ok := <-tab.EventChan
		if ok {
			Event = event
			switch Event.EventID {
			case HASH_TABLE_EVENT_ADD:
				Ret = tab.AddCore(Event.Arg)
				break
			case HASH_TABLE_EVENT_FIND:
				Ret = tab.FindCore(Event.Arg)
				break
			case HASH_TABLE_EVENT_REMOVE:
				Ret = tab.RemoveCore(Event.Arg)
				break
			default:
				fmt.Printf("Err: unknown command\n")
			}
			// return value
			Event.Ret <- Ret
		} else {
			running = false
		}
	}
}

// Delegate :
func (tab *HashTable) Delegate(EventID int, Arg HashTableEventArg) error {
	retchan := make(chan error)
	E := HashTableEvent{EventID, Arg, retchan}
	tab.EventChan <- E
	err, ok := <-E.Ret
	if ok {
		return err
	}
	return errors.New("Channel break")
}

// FindCore :
func (tab *HashTable) FindCore(Arg HashTableEventArg) error {
	_, ok := tab.Table[*(Arg.Key)]
	if ok {
		E := tab.Table[*(Arg.Key)]
		T := make([]byte, len(E.Value))
		*(Arg.Value) = &T
		for i := 0; i < len(E.Value); i++ {
			T[i] = E.Value[i]
		}
		return nil
	}
	return errors.New("Key not found")
}

// AddCore :
func (tab *HashTable) AddCore(Arg HashTableEventArg) error {
	tab.Table[*(Arg.Key)] = HashTableEntry{*(Arg.Key), **(Arg.Value)}
	return nil
}

// RemoveCore : FIND_NODE
func (tab *HashTable) RemoveCore(Arg HashTableEventArg) error {
	_, ok := tab.Table[*(Arg.Key)]
	if ok {
		delete(tab.Table, *(Arg.Key))
		return nil
	}
	return errors.New("Key not found")
}

/*
These functions are not thread safe, should not be called by other routine
*/

package libkademlia

import (
	"errors"
	"fmt"
)

const (
	ROUTING_TABLE_EVENT_UPDATE            = 1
	ROUTING_TABLE_EVENT_FIND_NEAREST_NODE = 2
	ROUTING_TABLE_EVENT_LOOK_UP           = 3
	ROUTING_TABLE_EVENT_FINALIZE          = 4
)

// RoutingTable : one more bucket for exactly the same, not used
type RoutingTable struct {
	Buckets   [b + 1]Bucket
	EventChan chan RountingTableEvent
	Self      *Kademlia
}

// RountingTableEvent :
type RountingTableEvent struct {
	EventID int
	Arg     RountingTableEventArg
	Ret     chan error
}

// RountingTableEventArg :
type RountingTableEventArg struct {
	ID *ID
	C  *Contact
	CS **[]Contact
}

// Dispatcher :
func (tab *RoutingTable) Dispatcher() {
	var Event RountingTableEvent
	var Ret error
	running := true
	for running {
		event, ok := <-tab.EventChan
		if ok {
			Event = event
			switch Event.EventID {
			case ROUTING_TABLE_EVENT_UPDATE:
				Ret = tab.UpdateCore(Event.Arg)
				break
			case ROUTING_TABLE_EVENT_FIND_NEAREST_NODE:
				Ret = tab.FindNearestNodeCore(Event.Arg)
				break
			case ROUTING_TABLE_EVENT_LOOK_UP:
				Ret = tab.LookUpCore(Event.Arg)
				break
			case ROUTING_TABLE_EVENT_FINALIZE:
				running = false
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
func (tab *RoutingTable) Delegate(EventID int, Arg RountingTableEventArg) error {
	retchan := make(chan error)
	E := RountingTableEvent{EventID, Arg, retchan}
	tab.EventChan <- E
	err, ok := <-E.Ret
	if ok {
		return err
	}
	return errors.New("Channel break")
}

// UpdateCore :
func (tab *RoutingTable) UpdateCore(Arg RountingTableEventArg) error {
	C := *(Arg.C)
	dist := (tab.Self.NodeID.Xor(C.NodeID)).PrefixLenEx()
	if dist < b {
		err := tab.Buckets[dist].MoveFront(C)
		if err != nil { // Not in list
			if tab.Buckets[dist].size < k { // Not full
				tab.Buckets[dist].PushBack(C)
				return nil
			}
			H, _ := tab.Buckets[dist].Top()
			_, err = tab.Self.DoInternalPing(H.Host, H.Port)
			if err != nil { // Can ping head
				return errors.New("Bucket full")
			}
			tab.Buckets[dist].Pop()
			tab.Buckets[dist].PushBack(C)
			return nil
		}
	} else {
		return errors.New("Can not add self to table")
	}
	return nil
}

// FindNearestNodeCore :
func (tab *RoutingTable) FindNearestNodeCore(Arg RountingTableEventArg) error {
	var C []Contact
	id := *Arg.ID
	dist := (tab.Self.NodeID.Xor(id)).PrefixLenEx()
	for j := dist; j < b; j++ {
		if len(C) < k {
			for i := 0; i < tab.Buckets[j].size && len(C) < k; i++ {
				T, _ := tab.Buckets[j].Get(i)
				C = append(C, T)
			}
		} else {
			break
		}
	}
	for j := dist - 1; j > -1; j-- {
		if len(C) < k {
			for i := 0; i < tab.Buckets[j].size && len(C) < k; i++ {
				T, _ := tab.Buckets[j].Get(i)
				C = append(C, T)
			}
		} else {
			break
		}
	}
	*Arg.CS = &C
	return nil
}

// LookUpCore : ID to Contact
func (tab *RoutingTable) LookUpCore(Arg RountingTableEventArg) error {
	id := *(Arg.ID)
	dist := (tab.Self.NodeID.Xor(id)).PrefixLenEx()
	C, err := tab.Buckets[dist].Find(id)
	*(Arg.C) = C
	return err
}

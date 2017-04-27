/*
These functions are not thread safe, should not be called by other routine
*/

package libkademlia

import (
	"errors"
	"fmt"
)

// RoutingTable : one more for exactly the same, not used
type RoutingTable [b + 1]Bucket

const (
	TABLE_EVENT_UPDATE            = 1
	TABLE_EVENT_FIND_NEAREST_NODE = 2
)

// TableEvent :
type TableEvent struct {
	EventID int
	Arg     TableEventArg
	Ret     chan error
}

// TableEventArg :
type TableEventArg struct {
	CIn   *Contact
	CSIn  *[]Contact
	IDIn  *ID
	VIn   *[]byte
	COut  *Contact
	CSOut *[]Contact
	IDOut *ID
	VOut  *[]byte
}

// TableDispatcher :
func (kad *Kademlia) TableDispatcher() {
	var Event TableEvent
	var Ret error
	running := true
	for running {
		event, ok := <-kad.TableRoutingCh
		if ok {
			Event = event
			switch Event.EventID {
			case TABLE_EVENT_UPDATE:
				Ret = kad.TableUpdateCore(Event.Arg)
				break
			case TABLE_EVENT_FIND_NEAREST_NODE:
				Ret = kad.FindNearestNodeCore(Event.Arg)
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

// TableDelegate :
func (kad *Kademlia) TableDelegate(EventID int, Arg TableEventArg) error {
	retchan := make(chan error)
	E := TableEvent{EventID, Arg, retchan}
	kad.TableRoutingCh <- E
	err, ok := <-E.Ret
	if ok {
		return err
	}
	return errors.New("Channel break")
}

// TableInit : Not thread safe, should be called only once
func (kad *Kademlia) TableInit() error {
	for i := 0; i < k; i++ {
		kad.Table[i].Init()
	}
	kad.TableRoutingCh = make(chan TableEvent)
	go kad.TableDispatcher()
	return nil
}

// TableUpdateCore :
func (kad *Kademlia) TableUpdateCore(Arg TableEventArg) error {
	C := *(Arg.CIn)
	dist := (kad.NodeID.Xor(C.NodeID)).PrefixLen()
	if dist < b {
		err := kad.Table[dist].MoveFront(C)
		if err != nil { // Not in list
			if kad.Table[dist].size < k { // Not full
				kad.Table[dist].PushBack(C)
				return nil
			}
			H, _ := kad.Table[dist].Top()
			_, err = kad.DoPing(H.Host, H.Port)
			if err != nil { // Can ping head
				return errors.New("Bucket full")
			}
			kad.Table[dist].Pop()
			kad.Table[dist].PushBack(C)
			return nil
		}
	} else {
		return errors.New("Can not add self to table")
	}
	return nil
}

// FindNearestNodeCore :
func (kad *Kademlia) FindNearestNodeCore(Arg TableEventArg) error {
	var C []Contact
	id := *Arg.IDIn
	dist := (kad.NodeID.Xor(id)).PrefixLen()
	for j := dist; j < b; j++ {
		if len(C) < k {
			for i := 0; i < kad.Table[j].size && len(C) < k; i++ {
				T, _ := kad.Table[j].Get(i)
				C = append(C, T)
			}
		} else {
			break
		}
	}
	for j := dist - 1; j > -1; j++ {
		if len(C) < k {
			for i := 0; i < kad.Table[j].size && len(C) < k; i++ {
				T, _ := kad.Table[j].Get(i)
				C = append(C, T)
			}
		} else {
			break
		}
	}
	Arg.CSOut = &C
	return nil
}

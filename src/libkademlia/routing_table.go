package libkademlia

import (
	"fmt"
)

const (
	ROUTING_EVENT_UPDATE = 1
)

type RoutingTable [b]Bucket

type RoutingEvent struct {
	EventId int
	In      chan Contact
	Out     chan Contact
}

func (k *Kademlia) TableDispatcher() {
	var event RoutingEvent
	running := true
	for running {
		e, ok := <-k.RoutingCh
		event = e
		var args []Contact
		hasargs := true
		for (hasargs) {
			arg, more := <- event.In
			if more {
				append(args, arg)
			}
			else{
				hasargs = false
			}
		}
		switch Event.EventId {
		case ROUTING_EVENT_UPDATE:
			

			break
		default:
			fmt.Printf("Err: unknown command\n")
		}
	}
}

func (k *Kademlia) TableUpdate(args []Contact){

}



func (k *Kademlia) Table_Init() error {
	for i := 0; i < k; i++ {
		k.Table[i].head = 0;
		k.Table[i].tail = 0;
	}

	go k.TableDispatcher()
	return nil
}

package libkademlia

import (
	"fmt"
	"net"
)

const (
	ROUTING_EVENT_UPDATE = 1
)

type Bucket struct {
	Entries [k]Contact
	head    int
	tail    int
}

type RoutingTable [b]Bucket

type RoutingEvent struct{
	EventId int
	In chan Contact
	Out chan Contact
}

func (k *Kademlia) TableDispatcher() {
	var Event RoutingEvent;
	running := true
	while(running){
		Event <- k.RoutingCh
		switch(Event.EventId){
			case ROUTING_EVENT_UPDATE:
				fmt.Println("OS X.")
				break
			default:
				fmt.Printf("%s.", os)
		}
	}
}

func (k *Kademlia) Table_Init() (error) {
	// TODO: Implement
	return nil
}

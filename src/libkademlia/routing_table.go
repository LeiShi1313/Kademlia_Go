/*
These functions except Init and Finalize are thread safe
*/

package libkademlia

// Init : Not thread safe, should be called only once. Must be called before all other functions can work
func (tab *RoutingTable) Init(Self *Kademlia) error {
	for i := 0; i < k; i++ {
		tab.Buckets[i].Init()
	}
	tab.EventChan = make(chan RountingTableEvent)
	tab.Self = Self
	go tab.Dispatcher()
	return nil
}

// Finalize : Not thread safe, should be called only once. Must be called before program exit. All functions can't be called after Finalize
func (tab *RoutingTable) Finalize(Self *Kademlia) error {
	close(tab.EventChan)
	return nil
}

// UpdateN :
func (tab *RoutingTable) UpdateN(C []Contact) (err []error) {
	for i := 0; i < len(C); i++ {
		err = append(err, tab.Update(C[i]))
	}
	return
}

// Update :
func (tab *RoutingTable) Update(C Contact) error {
	E := RountingTableEventArg{nil, &C, nil}
	return tab.Delegate(ROUTING_TABLE_EVENT_UPDATE, E)
}

// FindNearestNode : FIND_NODE
func (tab *RoutingTable) FindNearestNode(id ID) (C []Contact, num int, err error) {
	E := RountingTableEventArg{&id, nil, &C}
	ret := tab.Delegate(ROUTING_TABLE_EVENT_FIND_NEAREST_NODE, E)
	C = *E.CS
	return C, len(C), ret
}

// LookUp : ID to Contact
func (tab *RoutingTable) LookUp(id ID) (C Contact, err error) {
	E := RountingTableEventArg{&id, &C, nil}
	ret := tab.Delegate(ROUTING_TABLE_EVENT_LOOK_UP, E)
	C = *E.C
	return C, ret
}
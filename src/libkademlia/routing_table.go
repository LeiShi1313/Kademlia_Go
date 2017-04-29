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
	var T *[]Contact
	E := RountingTableEventArg{&id, nil, &T}
	ret := tab.Delegate(ROUTING_TABLE_EVENT_FIND_NEAREST_NODE, E)
	C = **(E.CS)
	return C, len(C), ret
}

// LookUp : ID to Contact
func (tab *RoutingTable) LookUp(id ID) (C Contact, err error) {
	var T Contact
	E := RountingTableEventArg{&id, &T, nil}
	ret := tab.Delegate(ROUTING_TABLE_EVENT_LOOK_UP, E)
	C = T
	return C, ret
}

// Size :
func (tab *RoutingTable) Size() int {
	ret := 0
	for i := 0; i < k; i++ {
		ret += tab.Buckets[i].size
	}
	return ret
}

// Info : Return size of each bucket
func (tab *RoutingTable) Info() []int {
	info := make([]int, b+1)
	for i, bucket := range tab.Buckets {
		info[i] = bucket.size
	}
	return info
}

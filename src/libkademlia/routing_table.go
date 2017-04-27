/*
These functions are thread safe
*/

package libkademlia

// TableUpdateN :
func (kad *Kademlia) TableUpdateN(C []Contact) (err []error) {
	for i := 0; i < len(C); i++ {
		err = append(err, kad.TableUpdate(C[i]))
	}
	return
}

// TableUpdate :
func (kad *Kademlia) TableUpdate(C Contact) error {
	E := TableEventArg{&C, nil, nil, nil, nil, nil, nil, nil}
	return kad.TableDelegate(TABLE_EVENT_UPDATE, E)
}

// FindNearestNode : FIND_NODE
func (kad *Kademlia) FindNearestNode(id ID) (C []Contact, num int, err error) {
	E := TableEventArg{nil, nil, &id, nil, nil, &C, nil, nil}
	ret := kad.TableDelegate(TABLE_EVENT_UPDATE, E)
	C = *E.CSOut
	return C, len(C), ret
}

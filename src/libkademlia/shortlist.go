/*
These functions are not thread safe, should not be called by other routine
*/

package libkademlia

import (
	"errors"
	"sort"
)

// ShortList :
type ShortList struct {
	Entries    map[ID]ShortListEntry
	ClosetNode *ShortListEntry
	Target     ID
	Parent     *Kademlia
}

// ShortListEntry :
type ShortListEntry struct {
	Conn   Contact
	Dist   int
	Active bool
}

type EntryList []ShortListEntry

func (L EntryList) Len() int {
	return len(L)
}
func (L EntryList) Swap(i, j int) {
	L[i], L[j] = L[j], L[i]
}
func (L EntryList) Less(i, j int) bool {
	return L[i].Dist < L[j].Dist
}

// Init : Not thread safe, should be called only once. Must be called before all other functions can work
func (l *ShortList) Init(Self *Kademlia) error {
	l.Parent = Self
	l.Entries = make(map[ID]ShortListEntry)
	l.ClosetNode = nil // ClosetNode undefined
	return nil
}

// Add : Not thread safe
func (l *ShortList) Add(C Contact) error {
	_, ok := l.Entries[C.NodeID]
	if ok {
		return errors.New("Already in list")
	}
	E := ShortListEntry{C, C.NodeID.Xor(l.Target).PrefixLenEx(), false}
	l.Entries[C.NodeID] = E
	return nil
}

// MAdd : Not thread safe
func (l *ShortList) MAdd(C []Contact) (err []error) {
	for i := 0; i < len(C); i++ {
		err = append(err, l.Add(C[i]))
	}
	return err
}

// Remove : Not thread safe
func (l *ShortList) Remove(id ID) error {
	_, ok := l.Entries[id]
	if !ok {
		return errors.New("Not in list")
	}
	delete(l.Entries, id)
	return nil
}

// Get : Not thread safe
func (l *ShortList) Find(id ID) (E ShortListEntry, err error) {
	var ok bool
	E, ok = l.Entries[id]
	if !ok {
		return E, errors.New("Not in list")
	}
	return E, nil;
}

// Get : Not thread safe
func (l *ShortList) FindContact(id ID) (C Contact, err error) {
	E, ok := l.Find(id)
	if !ok {
		return C, errors.New("Not in list")
	}
	return E.Conn, nil
}

// GetAll : Not thread safe
func (l *ShortList) GetAll() (C []Contact) {
	for _, E := range l.Entries {
		C = append(C, E.Conn)
	}
	return C
}

// GetActive : Not thread safe
func (l *ShortList) GetActive() (C []Contact) {
	for _, E := range l.Entries {
		if E.Active {
			C = append(C, E.Conn)
		}
	}
	return C
}

// GetNearest : Not thread safe
func (l *ShortList) GetNearest(n int) C []Contact {
	E := l.GetActive()
	sort.Sort(EntryList(E))
	for i := 0; i < n && i < len(E); i++{
		C = append(C, E[i].)
	}
	return C
}

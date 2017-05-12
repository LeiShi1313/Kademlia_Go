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
	Entries          map[ID]ShortListEntry
	ClosetNode       *ShortListEntry
	ClosetActiveNode *ShortListEntry
	Target           ID
	Parent           *Kademlia
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
func (l *ShortList) Init(Self *Kademlia, target ID) error {
	l.Parent = Self
	l.Entries = make(map[ID]ShortListEntry)
	l.ClosetNode = nil       // ClosetNode undefined
	l.ClosetActiveNode = nil // ClosetNode undefined
	l.Target = target
	return nil
}

// Add : Not thread safe
func (l *ShortList) Add(C Contact) error {
	_, ok := l.Entries[C.NodeID]
	if ok {
		return errors.New("Already in list")
	}
	E := ShortListEntry{C, C.NodeID.Xor(l.Target).PrefixLenEx(), false}
	if l.ClosetNode == nil || l.ClosetNode.Dist > E.Dist {
		l.ClosetNode = &E
	}
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

// Find : Not thread safe
func (l *ShortList) Find(id ID) (E ShortListEntry, err error) {
	var ok bool
	E, ok = l.Entries[id]
	if !ok {
		return E, errors.New("Not in list")
	}
	return E, nil
}

// FindContact : Not thread safe
func (l *ShortList) FindContact(id ID) (Contact, error) {
	E, err := l.Find(id)
	return E.Conn, err
}

// GetAllEntry : Not thread safe
func (l *ShortList) GetAllEntry() (Ent []ShortListEntry) {
	for _, E := range l.Entries {
		Ent = append(Ent, E)
	}
	return Ent
}

// GetAllContact : Not thread safe
func (l *ShortList) GetAllContact() (C []Contact) {
	for _, E := range l.Entries {
		C = append(C, E.Conn)
	}
	return C
}

// GetActiveEntry : Not thread safe
func (l *ShortList) GetActiveEntry() (Ent []ShortListEntry) {
	for _, E := range l.Entries {
		if E.Active {
			Ent = append(Ent, E)
		}
	}
	return Ent
}

// GetActiveContact : Not thread safe
func (l *ShortList) GetActiveContact() (C []Contact) {
	for _, E := range l.Entries {
		if E.Active {
			C = append(C, E.Conn)
		}
	}
	return C
}

// GetInactiveEntry : Not thread safe
func (l *ShortList) GetInactiveEntry() (Ent []ShortListEntry) {
	for _, E := range l.Entries {
		if !E.Active {
			Ent = append(Ent, E)
		}
	}
	return Ent
}

// GetInactiveContact : Not thread safe
func (l *ShortList) GetInactiveContact() (C []Contact) {
	for _, E := range l.Entries {
		if !E.Active {
			C = append(C, E.Conn)
		}
	}
	return C
}

// GetNearestN : Get nearest n nodes that's not proven active
func (l *ShortList) GetNearestN(n int) (C []Contact) {
	E := l.GetInactiveEntry()
	sort.Sort(EntryList(E))
	for i := 0; i < n && i < len(E); i++ {
		C = append(C, E[i].Conn)
	}
	return C
}

// Size : Size of short list
func (l *ShortList) Size() int {
	return len(l.Entries)
}

// ActiveSize : Number of active nodes
func (l *ShortList) ActiveSize() int {
	ret := 0
	for _, E := range l.Entries {
		if E.Active {
			ret++
		}
	}
	return ret
}

// SetActive : Number of active nodes
func (l *ShortList) SetActive(id ID) error {
	E, ok := l.Entries[id]
	if !ok {
		return errors.New("Not in list")
	}
	if l.ClosetActiveNode == nil || l.ClosetActiveNode.Dist > E.Dist {
		l.ClosetNode = &E
	}
	E.Active = true
	l.Entries[id] = E
	return nil
}

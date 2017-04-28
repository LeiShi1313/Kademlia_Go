/*
These functions except Init and Finalize are thread safe
*/

package libkademlia

// Init : Not thread safe, should be called only once. Must be called before all other functions can work
func (tab *HashTable) Init(Self *Kademlia) error {
	tab.Table = make(map[ID]HashTableEntry)
	tab.Self = Self
	tab.EventChan = make(chan HashTableEvent)
	go tab.Dispatcher()
	return nil
}

// Finalize : Not thread safe, should be called only once. Must be called before program exit. All functions can't be called after Finalize
func (tab *HashTable) Finalize() error {
	close(tab.EventChan)
	return nil
}

// Find :
func (tab *HashTable) Find(key ID) (V []byte, err error) {
	var varp *[]byte
	E := HashTableEventArg{&key, &varp}
	err = tab.Delegate(HASH_TABLE_EVENT_FIND, E)
	if err == nil {
		V = **(E.Value)
	}
	return V, err
}

// Add : Adding existing key overwrites the value
func (tab *HashTable) Add(key ID, value []byte) error {
	var varp *[]byte
	varp = &value
	E := HashTableEventArg{&key, &varp}
	return tab.Delegate(HASH_TABLE_EVENT_ADD, E)
}

// Remove : FIND_NODE
func (tab *HashTable) Remove(key ID) error {
	E := HashTableEventArg{&key, nil}
	return tab.Delegate(HASH_TABLE_EVENT_REMOVE, E)
}

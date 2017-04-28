package libkademlia

// NOTE: All functions are not thread safe

import (
	"errors"
)

// Bucket :
type Bucket struct {
	Entries [k]Contact
	head    int
	size    int
}

// Init :
func (bkt *Bucket) Init() {
	bkt.head = 0
	bkt.size = 0
}

// PushBack :
func (bkt *Bucket) PushBack(C Contact) (err error) {
	if bkt.size < k {
		bkt.Entries[(bkt.head+bkt.size)%k] = C
		bkt.size++
		return nil
	}
	return errors.New("Bucket full")
}

// Pop :
func (bkt *Bucket) Pop() (C Contact, err error) {
	if bkt.size > 0 {
		C = bkt.Entries[bkt.head]
		bkt.head = (bkt.head + 1) % k
		return C, nil
	}
	return C, errors.New("Bucket empty")
}

// Top :
func (bkt *Bucket) Top() (C Contact, err error) {
	return bkt.Get(0)
}

// Get :
func (bkt *Bucket) Get(idx int) (C Contact, err error) {
	if bkt.size > idx {
		C = bkt.Entries[(bkt.head+idx)%k]
		return C, nil
	}
	return C, errors.New("Bucket empty")
}

// MoveFront :
func (bkt *Bucket) MoveFront(C Contact) error {
	var i, j int
	for i = 0; i < bkt.size; i++ {
		if bkt.Entries[(i+bkt.size)%k].NodeID.Equals(C.NodeID) {
			break
		}
	}
	if i < bkt.size {
		for j = i + 1; j < bkt.size; j++ {
			T := bkt.Entries[(j-1+bkt.size)%k]
			bkt.Entries[(j-1+bkt.size)%k] = bkt.Entries[(j+bkt.size)%k]
			bkt.Entries[(j+bkt.size)%k] = T
		}
		return nil
	}
	return errors.New("Contact not found")
}

// Find :
func (bkt *Bucket) Find(id ID) (C Contact, err error) {
	var i int
	for i = 0; i < bkt.size; i++ {
		if bkt.Entries[(i+bkt.head)%k].NodeID.Equals(id) {
			C = bkt.Entries[(i+bkt.head)%k]
			err = nil
			return C, err
		}
	}
	return C, errors.New("ID not in bucket")
}

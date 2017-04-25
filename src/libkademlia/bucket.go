// NOTE: All functions are not thread safe
package libkademlia

type Bucket struct {
	Entries [k]Contact
	head    int
	tail    int
}

func (self *Bucket) Init() error {
	self.head = 0
	self.tail = 0
}

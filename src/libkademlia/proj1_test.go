package libkademlia

import (
	"bytes"
	"net"
	"strconv"
	"testing"
	//"time"
)

func StringToIpPort(laddr string) (ip net.IP, port uint16, err error) {
	hostString, portString, err := net.SplitHostPort(laddr)
	if err != nil {
		return
	}
	ipStr, err := net.LookupHost(hostString)
	if err != nil {
		return
	}
	for i := 0; i < len(ipStr); i++ {
		ip = net.ParseIP(ipStr[i])
		if ip.To4() != nil {
			break
		}
	}
	portInt, err := strconv.Atoi(portString)
	port = uint16(portInt)
	return
}

func TestPing(t *testing.T) {
	instance1 := NewKademlia("localhost:7890")
	instance2 := NewKademlia("localhost:7891")
	host2, port2, _ := StringToIpPort("localhost:7891")
	contact2, err := instance2.FindContact(instance2.NodeID)
	if err != nil {
		t.Error("A node cannot find itself's contact info")
	}
	contact2, err = instance2.FindContact(instance1.NodeID)
	if err == nil {
		t.Error("Instance 2 should not be able to find instance " +
			"1 in its buckets before ping instance 1")
	}
	instance1.DoPing(host2, port2)
	contact2, err = instance1.FindContact(instance2.NodeID)
	if err != nil {
		t.Error("Instance 2's contact not found in Instance 1's contact list")
		return
	}
	wrong_ID := NewRandomID()
	_, err = instance2.FindContact(wrong_ID)
	if err == nil {
		t.Error("Instance 2 should not be able to find a node with the wrong ID")
	}

	contact1, err := instance2.FindContact(instance1.NodeID)
	if err != nil {
		t.Error("Instance 1's contact not found in Instance 2's contact list")
		return
	}
	if contact1.NodeID != instance1.NodeID {
		t.Error("Instance 1 ID incorrectly stored in Instance 2's contact list")
	}
	if contact2.NodeID != instance2.NodeID {
		t.Error("Instance 2 ID incorrectly stored in Instance 1's contact list")
	}
	return
}

func TestStore(t *testing.T) {
	// test Dostore() function and LocalFindValue() function
	instance1 := NewKademlia("localhost:7892")
	instance2 := NewKademlia("localhost:7893")
	host2, port2, _ := StringToIpPort("localhost:7893")
	instance1.DoPing(host2, port2)
	contact2, err := instance1.FindContact(instance2.NodeID)
	if err != nil {
		t.Error("Instance 2's contact not found in Instance 1's contact list")
		return
	}
	key := NewRandomID()
	value := []byte("Hello World")
	err = instance1.DoStore(contact2, key, value)
	if err != nil {
		t.Error("Can not store this value")
	}
	storedValue, err := instance2.LocalFindValue(key)
	if err != nil {
		t.Error("Stored value not found!")
	}
	if !bytes.Equal(storedValue, value) {
		t.Error("Stored value did not match found value")
	}
	return
}

func TestFindNode(t *testing.T) {
	// tree structure;
	// A->B->tree
	/*
	         C
	      /
	  A-B -- D
	      \
	         E
	*/
	instance1 := NewKademlia("localhost:7894")
	instance2 := NewKademlia("localhost:7895")
	host2, port2, _ := StringToIpPort("localhost:7895")
	instance1.DoPing(host2, port2)
	contact2, err := instance1.FindContact(instance2.NodeID)
	if err != nil {
		t.Error("Instance 2's contact not found in Instance 1's contact list")
		return
	}
	tree_node := make([]*Kademlia, 10)
	for i := 0; i < 10; i++ {
		address := "localhost:" + strconv.Itoa(7896+i)
		tree_node[i] = NewKademlia(address)
		host_number, port_number, _ := StringToIpPort(address)
		instance2.DoPing(host_number, port_number)
	}
	key := NewRandomID()
	contacts, err := instance1.DoFindNode(contact2, key)
	if err != nil {
		t.Error("Error doing FindNode")
	}

	if contacts == nil || len(contacts) == 0 {
		t.Error("No contacts were found")
	}

	var j int
	for i := 0; i < len(contacts); i++ {
		for j = 0; j < len(tree_node); j++ {
			if contacts[i].NodeID.Equals(tree_node[j].SelfContact.NodeID) {
				break
			}
		}
		if j == len(tree_node) {
			if !(contacts[i].NodeID.Equals(instance1.SelfContact.NodeID)) {
				t.Error("Wrong contacts")
			}
		}
	}

	// TODO: Check that the correct contacts were stored
	//       (and no other contacts)

	return
}

func TestFindValue(t *testing.T) {
	// tree structure;
	// A->B->tree
	/*
	         C
	      /
	  A-B -- D
	      \
	         E
	*/
	instance1 := NewKademlia("localhost:7926")
	instance2 := NewKademlia("localhost:7927")
	host2, port2, _ := StringToIpPort("localhost:7927")
	instance1.DoPing(host2, port2)
	contact2, err := instance1.FindContact(instance2.NodeID)
	if err != nil {
		t.Error("Instance 2's contact not found in Instance 1's contact list")
		return
	}

	tree_node := make([]*Kademlia, 10)
	for i := 0; i < 10; i++ {
		address := "localhost:" + strconv.Itoa(7928+i)
		tree_node[i] = NewKademlia(address)
		host_number, port_number, _ := StringToIpPort(address)
		instance2.DoPing(host_number, port_number)
	}

	key := NewRandomID()
	value := []byte("Hello world")
	err = instance2.DoStore(contact2, key, value)
	if err != nil {
		t.Error("Could not store value")
	}

	// Given the right keyID, it should return the value
	foundValue, contacts, err := instance1.DoFindValue(contact2, key)
	if !bytes.Equal(foundValue, value) {
		t.Error("Stored value did not match found value")
	}

	//Given the wrong keyID, it should return k nodes.
	wrongKey := NewRandomID()
	foundValue, contacts, err = instance1.DoFindValue(contact2, wrongKey)
	if contacts == nil || len(contacts) < 10 || len(contacts) > 11 {
		t.Error("Searching for a wrong ID did not return contacts")
	}

	var j int
	for i := 0; i < len(contacts); i++ {
		for j = 0; j < len(tree_node); j++ {
			if contacts[i].NodeID.Equals(tree_node[j].SelfContact.NodeID) {
				break
			}
		}
		if j == len(tree_node) {
			if !(contacts[i].NodeID.Equals(instance1.SelfContact.NodeID)) {
				t.Error("Wrong contacts")
			}
		}
	}

	// TODO: Check that the correct contacts were stored
	//       (and no other contacts)
}

func TestHashTable(t *testing.T) {
	var H HashTable
	ida := NewRandomID()
	idb := NewRandomID()
	var self *Kademlia

	va := []byte("Value A")
	vb := []byte("Value B")
	H.Init(self)
	H.Add(ida, va)
	V, err := H.Find(idb)
	if err == nil {
		t.Error("Id b shouldn't be in table")
	}
	V, err = H.Find(ida)
	if err != nil {
		t.Error("Id a should be in table")
	}
	if len(V) != len(va) {
		t.Error("Value a mismatch")
	}
	for i := 0; i < len(V); i++ {
		if V[i] != va[i] {
			t.Error("Value a mismatch")
		}
	}
	H.Add(idb, vb)
	V, err = H.Find(ida)
	if err != nil {
		t.Error("Id a should be in table")
	}
	V, err = H.Find(idb)
	if err != nil {
		t.Error("Id b should be in table")
	}
	if len(V) != len(vb) {
		t.Error("Value a mismatch")
	}
	for i := 0; i < len(V); i++ {
		if V[i] != vb[i] {
			t.Error("Value b mismatch")
		}
	}
	V, err = H.Find(ida)
	if err != nil {
		t.Error("Id a should be in table")
	}
	if len(V) != len(va) {
		t.Error("Value a mismatch")
	}
	for i := 0; i < len(V); i++ {
		if V[i] != va[i] {
			t.Error("Value a mismatch")
		}
	}
	H.Remove(ida)
	V, err = H.Find(ida)
	if err == nil {
		t.Error("Id a shouldn't be in table")
	}
}

func TestFullBucket(t *testing.T) {
	instance1 := NewKademlia("localhost:10001")
	instance2 := NewKademlia("localhost:10002")
	host2, port2, _ := StringToIpPort("localhost:10002")
	instance1.DoPing(host2, port2)
	// instance2Size, instance2Info := instance2.GetRoutingTableInfo()
	// fmt.Printf("%v\n%v\n", instance2Size, instance2Info)
	tree_node := make([]*Kademlia, 40)
	firstNodeID := instance2.NodeID
	nodeId := CopyID(firstNodeID)
	nodeId[13] = 128
	for i := 0; i < 40; i++ {
		if i < 20 {
			nodeId = nodeId.Increse(1)
		} else {
			nodeId[14] = 128
			nodeId = nodeId.Increse(1)
		}
		// fmt.Printf("%v  ", nodeId.Xor(firstNodeID))
		address := "localhost:" + strconv.Itoa(10003+i)
		tree_node[i] = NewKademliaWithId(address, nodeId)
		host_number, port_number, _ := StringToIpPort(address)
		instance2.DoPing(host_number, port_number)
		// instance2Size, instance2Info := instance2.GetRoutingTableInfo()
		// fmt.Printf("%v\n%v\n", instance2Size, instance2Info)
	}

	instance2Size, _ := instance2.GetRoutingTableInfo()
	// fmt.Printf("%v\n%v\n", instance2Size, instance2Info)
	if instance2Size != 21 {
		t.Error("Full bucket panic")
	}
}


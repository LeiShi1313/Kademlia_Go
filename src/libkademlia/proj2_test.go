package libkademlia

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	//"time"
)

func TestNodeLeave(t *testing.T) {
	instance1 := NewKademlia("localhost:20001")
	instance2 := NewKademlia("localhost:20002")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	host2, port2, _ := StringToIpPort("localhost:20002")
	_, err := instance1.DoPing(host2, port2)
	if err != nil {
		t.Error("Can't ping active peer")
	}
	host2, port2, _ = StringToIpPort("localhost:20001")
	_, err = instance2.DoPing(host2, port2)
	if err != nil {
		t.Error("Can't ping active peer")
	}

	host3, port3, _ := StringToIpPort("localhost:20003")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	_, err = instance1.DoPing(host3, port3)
	if err == nil {
		t.Error("Can ping downed peer")
	}
}

func TestIterativeFindNode(t *testing.T) {
	instance1 := NewKademlia("localhost:7101")
	instance2 := NewKademlia("localhost:7102")
	host2, port2, _ := StringToIpPort("localhost:7102")
	_, err := instance1.DoPing(host2, port2)
	if err != nil {
		t.Error("Can't ping instance1")
	}
	if instance2Size, _ := instance2.GetRoutingTableInfo(); instance2Size != 1 {
		t.Error("Kademlia.GetRoutingTableInfo return incorrect size")
	}
	tree_node := make([]*Kademlia, 40)
	for i := 0; i < 40; i++ {
		address := "localhost:" + strconv.Itoa(7103+i)
		tree_node[i] = NewKademlia(address)
		host_number, port_number, _ := StringToIpPort(address)
		_, err = instance2.DoPing(host_number, port_number)
		if err != nil {
			t.Error("Can't ping instance" + strconv.Itoa(i+3))
		}
	}
	if instance2Size, _ := instance2.GetRoutingTableInfo(); instance2Size < 30 {
		t.Error("Kademlia.GetRoutingTableInfo return incorrect size")
	}
	id := NewRandomID()
	contacts, err := instance2.DoIterativeFindNode(id)
	if err != nil {
		t.Error(fmt.Sprint(err))
	}
	if len(contacts) != 20 {
		t.Error("Iterative store return insufficient number of contacts")
	}
}

func TestIterativeStore(t *testing.T) {
	instance1 := NewKademlia("localhost:30001")
	instance2 := NewKademlia("localhost:30002")
	host2, port2, _ := StringToIpPort("localhost:30002")
	_, err := instance1.DoPing(host2, port2)
	if err != nil {
		t.Error("Can't ping instance1")
	}
	if instance2Size, _ := instance2.GetRoutingTableInfo(); instance2Size != 1 {
		t.Error("Kademlia.GetRoutingTableInfo return incorrect size")
	}
	tree_node := make([]*Kademlia, 40)
	for i := 0; i < 40; i++ {
		address := "localhost:" + strconv.Itoa(30003+i)
		tree_node[i] = NewKademlia(address)
		host_number, port_number, _ := StringToIpPort(address)
		_, err = instance2.DoPing(host_number, port_number)
		if err != nil {
			t.Error("Can't ping instance" + strconv.Itoa(i+3))
		}
	}
	if instance2Size, _ := instance2.GetRoutingTableInfo(); instance2Size < 30 {
		t.Error("Kademlia.GetRoutingTableInfo return incorrect size")
	}
	key := NewRandomID()
	val := []byte("Hi there")
	instance2.DoIterativeStore(key, val)
	count := 0
	for i := 0; i < 40; i++ {
		result, err := tree_node[i].LocalFindValue(key)
		if err != nil {
			// fmt.Println(i, ":", err)
		} else {
			if bytes.Equal(result, val) {
				count++
			}
			// fmt.Println(i, ":", string(result[:]))
		}

	}
	if count < 3 {
		t.Error("IterativeStore failed, cannot get enough data")
	}
}

func TestIterativeFindValue(t *testing.T) {
	instance1 := NewKademlia("localhost:40001")
	instance2 := NewKademlia("localhost:40002")
	host2, port2, _ := StringToIpPort("localhost:40002")
	_, err := instance1.DoPing(host2, port2)
	if err != nil {
		t.Error("Can't ping instance1")
	}
	instance2Size, _ := instance2.GetRoutingTableInfo()
	// fmt.Printf("%v\n%v\n", instance2Size, instance2Info)
	if instance2Size != 1 {
		t.Error("Kademlia.GetRoutingTableInfo return incorrect size")
	}
	tree_node := make([]*Kademlia, 40)
	for i := 0; i < 40; i++ {
		address := "localhost:" + strconv.Itoa(40003+i)
		tree_node[i] = NewKademlia(address)
		host_number, port_number, _ := StringToIpPort(address)
		_, err = instance2.DoPing(host_number, port_number)
		if err != nil {
			t.Error("Can't ping instance" + strconv.Itoa(i+3))
		}
	}
	instance2Size, _ = instance2.GetRoutingTableInfo()
	// fmt.Printf("%v\n%v\n", instance2Size, instance2Info)

	if instance2Size < 30 {
		t.Error("Kademlia.GetRoutingTableInfo return incorrect size")
	}
	key := NewRandomID()
	val := []byte("Hi there, I'm finding something")
	instance2.DoIterativeStore(key, val)
	result, err := instance2.DoIterativeFindValue(key)
	if err != nil {
		t.Error(fmt.Sprint(err))
	}
	if !bytes.Equal(result, val) {
		t.Error("IterativeStore failed")
		t.Error(fmt.Sprint("Except: ", val))
		t.Error(fmt.Sprint("Got: ", result))
	}
}

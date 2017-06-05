package libkademlia

import (
	"bytes"
	"fmt"
	"net/rpc"
	"strconv"
	"testing"
)

func TestGetVDORPC(t *testing.T) {
	instance1 := NewKademlia("localhost:9400")
	instance2 := NewKademlia("localhost:9401")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	host2, port2, _ := StringToIpPort("localhost:9401")
	_, err := instance1.DoPing(host2, port2)
	if err != nil {
		t.Error("Can't ping instance2")
	}
	treeNode := make([]*Kademlia, 10)
	for i := 0; i < 10; i++ {
		address := "localhost:" + strconv.Itoa(9402+i)
		treeNode[i] = NewKademlia(address)
		hostNumber, portNumber, _ := StringToIpPort(address)
		_, err = instance2.DoPing(hostNumber, portNumber)
		if err != nil {
			t.Error("Can't ping instance" + strconv.Itoa(i+3))
		}
	}

	key, _ := IDFromString("Hello")
	data := []byte("World!")
	numberKeys := byte(10)
	threshold := byte(7)
	VDO := instance2.Vanish(key, data, numberKeys, threshold, 0)
	if VDO.NumberKeys <= 0 {
		t.Error("Vanish failed!")
	}

	peerStr := host2.String() + ":" + strconv.Itoa(int(port2))
	portStr := fmt.Sprint(port2)
	client, err := rpc.DialHTTPPath("tcp", peerStr, rpc.DefaultRPCPath+portStr)

	if err != nil {
		t.Error("Can't dial instance2")
	}
	msgID := NewRandomID()
	req := GetVDORequest{instance1.SelfContact, key, msgID}
	var reply GetVDOResult
	err = client.Call("KademliaRPC.GetVDO", req, &reply)
	if err != nil {
		t.Error("GetVDO RPC failed!")
	}
	if reply.Err.Msg != "" {
		t.Error(fmt.Println("RPC Error: ", reply.Err.Msg))
	}
	if reply.VDO.AccessKey != VDO.AccessKey || bytes.Compare(reply.VDO.Ciphertext, VDO.Ciphertext) != 0 {
		t.Error("VDO not expected!")
	}
}

func TestStoreVDO(t *testing.T) {
	instance1 := NewKademlia("localhost:9500")
	instance2 := NewKademlia("localhost:9501")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	host2, port2, _ := StringToIpPort("localhost:9501")
	_, err := instance1.DoPing(host2, port2)
	if err != nil {
		t.Error("Can't ping instance 2")
	}
	treeNode := make([]*Kademlia, 10)
	for i := 0; i < 10; i++ {
		address := "localhost:" + strconv.Itoa(9502+i)
		treeNode[i] = NewKademlia(address)
		hostNumber, portNumber, _ := StringToIpPort(address)
		_, err = instance2.DoPing(hostNumber, portNumber)
		if err != nil {
			t.Error("Can't ping instance" + strconv.Itoa(i+3))
		}
	}
	key, _ := IDFromString("Hello")
	data := []byte("World!")
	numberKeys := byte(10)
	threshold := byte(7)
	VDO := instance2.Vanish(key, data, numberKeys, threshold, 0)
	if VDO.NumberKeys <= 0 {
		t.Error("Vanish failed!")
	}
	ret := instance2.Unvanish(instance2.NodeID, key)
	if bytes.Compare(ret, data) != 0 {
		t.Error("Unvanish failed!")
		t.Error(fmt.Sprintln("Expect: ", data))
		t.Error(fmt.Sprintln("GOT: ", ret))
	}
}

func TestRetriveVDOFromOtherNode(t *testing.T) {
	instance1 := NewKademlia("localhost:9600")
	instance2 := NewKademlia("localhost:9601")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	host2, port2, _ := StringToIpPort("localhost:9601")
	_, err := instance1.DoPing(host2, port2)
	if err != nil {
		t.Error("Can't ping instance 2")
	}
	treeNode := make([]*Kademlia, 10)
	for i := 0; i < 10; i++ {
		address := "localhost:" + strconv.Itoa(9602+i)
		treeNode[i] = NewKademlia(address)
		hostNumber, portNumber, _ := StringToIpPort(address)
		_, err = instance2.DoPing(hostNumber, portNumber)
		if err != nil {
			t.Error("Can't ping instance" + strconv.Itoa(i+3))
		}
	}
	key, _ := IDFromString("Hello")
	data := []byte("World!")
	numberKeys := byte(10)
	threshold := byte(7)
	VDO := instance2.Vanish(key, data, numberKeys, threshold, 0)
	if VDO.NumberKeys <= 0 {
		t.Error("Vanish failed!")
	}

	treeNode[0].DoStoreVDO(key, VDO)
	instance2.DT.Remove(key)

	if _, err := instance2.DT.Find(key); err == nil {
		t.Error("Datatable delete failed!")
	}

	remoteVDO, err := treeNode[0].DT.Find(key)
	if err != nil {
		t.Error("Remote node find key failed!")
	}

	if remoteVDO.AccessKey != VDO.AccessKey || bytes.Compare(remoteVDO.Ciphertext, VDO.Ciphertext) != 0 {
		t.Error("Remote node return wrong VDO!")
	}

	ret := instance2.Unvanish(instance2.NodeID, key)
	if bytes.Compare(ret, data) != 0 {
		t.Error("Unvanish from remote node failed!")
		t.Error(fmt.Sprintln("Expect: ", data))
		t.Error(fmt.Sprintln("GOT: ", ret))
	}
}

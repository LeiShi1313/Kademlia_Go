package libkademlia

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
)

const (
	alpha = 3
	b     = 8 * IDBytes
	k     = 20
)

// Kademlia type. You can put whatever state you need in this.
type Kademlia struct {
	NodeID      ID
	SelfContact Contact
	RT          RoutingTable
	HT          HashTable
}

func NewKademliaWithId(laddr string, nodeID ID) *Kademlia {
	k := new(Kademlia)
	k.NodeID = nodeID

	// TODO: Initialize other state here as you add functionality.
	k.RT.Init(k)
	k.HT.Init(k)
	// Set up RPC server
	// NOTE: KademliaRPC is just a wrapper around Kademlia. This type includes
	// the RPC functions.

	s := rpc.NewServer()
	s.Register(&KademliaRPC{k})
	hostname, port, err := net.SplitHostPort(laddr)
	if err != nil {
		return nil
	}
	s.HandleHTTP(rpc.DefaultRPCPath+port,
		rpc.DefaultDebugPath+port)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Listen: ", err)
	}

	// Run RPC server forever.
	go http.Serve(l, nil)

	// Add self contact
	hostname, port, _ = net.SplitHostPort(l.Addr().String())
	if hostname == "::" {
		hostname = GetOutboundIP()
	}
	port_int, _ := strconv.Atoi(port)
	ipAddrStrings, err := net.LookupHost(hostname)
	var host net.IP
	for i := 0; i < len(ipAddrStrings); i++ {
		host = net.ParseIP(ipAddrStrings[i])
		if host.To4() != nil {
			break
		}
	}
	gob.Register(errors.New(""))
	k.SelfContact = Contact{k.NodeID, host, uint16(port_int)}
	return k
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

func NewKademlia(laddr string) *Kademlia {
	return NewKademliaWithId(laddr, NewRandomID())
}

type ContactNotFoundError struct {
	id  ID
	msg string
}

func (e *ContactNotFoundError) Error() string {
	return fmt.Sprintf("%x %s", e.id, e.msg)
}

func (k *Kademlia) FindContact(nodeId ID) (*Contact, error) {
	// TODO: Search through contacts, find specified ID
	// Find contact with provided ID
	if nodeId == k.SelfContact.NodeID {
		return &k.SelfContact, nil
	}
	contact, err := k.RT.LookUp(nodeId)
	return &contact, err
}

func (k *Kademlia) GetRoutingTableInfo() (total int, info []int) {
	info = k.RT.Info()
	total = k.RT.Size()
	return
}

type CommandFailed struct {
	msg string
}

func (e *CommandFailed) Error() string {
	return fmt.Sprintf("%s", e.msg)
}

func (k *Kademlia) Finalize() {
	k.RT.Finalize()
	k.HT.Finalize()

}

func (k *Kademlia) dial(host net.IP, port uint16) (*rpc.Client, error) {
	peerStr := host.String() + ":" + strconv.Itoa(int(port))
	portStr := fmt.Sprint(port)
	return rpc.DialHTTPPath("tcp", peerStr, rpc.DefaultRPCPath+portStr)
}

func (k *Kademlia) DoPing(host net.IP, port uint16) (*Contact, error) {
	client, err := k.dial(host, port)

	if err != nil {
		return nil, err
	}
	var reply PongMessage
	err = client.Call("KademliaRPC.Ping", PingMessage{k.SelfContact, NewRandomID()}, &reply)
	if err == nil {
		k.RT.Update(reply.Sender)
	}
	return &reply.Sender, err
}

/*
NOTE: This function can only be used within routing table core
*/
func (k *Kademlia) DoInternalPing(host net.IP, port uint16) (*Contact, error) {
	client, err := k.dial(host, port)

	if err != nil {
		return nil, err
	}
	var reply PongMessage
	err = client.Call("KademliaRPC.Ping", PingMessage{k.SelfContact, NewRandomID()}, &reply)
	if err == nil {
		k.RT.UpdateInternal(reply.Sender)
	}
	return &reply.Sender, err
}

func (k *Kademlia) DoStore(contact *Contact, key ID, value []byte) error {
	// TODO: Implement
	client, err := k.dial(contact.Host, contact.Port)
	if err != nil {
		return err
	}
	var reply StoreResult
	err = client.Call("KademliaRPC.Store", StoreRequest{k.SelfContact, NewRandomID(), key, value}, &reply)
	if err != nil {
		return err
	}
	return reply.Err
}

func (k *Kademlia) DoFindNode(contact *Contact, searchKey ID) ([]Contact, error) {
	// TODO: Implement
	client, err := k.dial(contact.Host, contact.Port)
	if err != nil {
		return nil, err
	}
	var reply FindNodeResult
	msgId := NewRandomID()
	err = client.Call("KademliaRPC.FindNode", FindNodeRequest{k.SelfContact, msgId, searchKey}, &reply)
	if err != nil {
		return nil, err
	}
	if reply.MsgID != msgId {
		return nil, &CommandFailed{"MsgId inconsitent"}
	}

	return reply.Nodes, nil
}

func (k *Kademlia) DoFindValue(contact *Contact,
	searchKey ID) (value []byte, contacts []Contact, err error) {
	// TODO: Implement
	client, err := k.dial(contact.Host, contact.Port)
	if err != nil {
		return nil, nil, err
	}
	var reply FindValueResult
	err = client.Call("KademliaRPC.FindValue", FindValueRequest{k.SelfContact, NewRandomID(), searchKey}, &reply)
	return reply.Value, reply.Nodes, reply.Err
}

func (k *Kademlia) LocalFindValue(searchKey ID) ([]byte, error) {
	// TODO: Implement

	return k.HT.Find(searchKey)
}

func (k *Kademlia) DoFindNodeAsync(contact *Contact, searchKey ID) (*rpc.Call, error) {
	// TODO: Implement
	client, err := k.dial(contact.Host, contact.Port)
	if err != nil {
		return nil, err
	}
	var reply FindNodeResult
	msgId := NewRandomID()
	return client.Go("KademliaRPC.FindNode", FindNodeRequest{k.SelfContact, msgId, searchKey}, &reply, nil), nil
}

func (k *Kademlia) doFindValueAsync(contact *Contact, key ID) (*rpc.Call, error) {
	client, err := k.dial(contact.Host, contact.Port)
	if err != nil {
		return nil, err
	}
	var reply FindValueResult
	msgId := NewRandomID()
	findValueRequest := FindValueRequest{k.SelfContact, msgId, key}
	return client.Go("KademliaRPC.FindValue", findValueRequest, &reply, nil), nil
}


func (k *Kademlia) DoFindNodeWait(Call *rpc.Call) ([]Contact, error) {
	if Call == nil {
		return nil, &CommandFailed{"Null handle"}
	}
	Ret := *(<-Call.Done)
	err := Ret.Error
	if err != nil {
		return nil, err
	}
	req, ok := Ret.Args.(FindNodeRequest)
	if !ok {
		return nil, &CommandFailed{"Unknown Error"}
	}
	reply, ok2 := Ret.Args.(FindNodeResult)
	if !ok2 {
		return nil, &CommandFailed{"Unknown Error"}
	}
	if reply.MsgID != req.MsgID {
		return nil, &CommandFailed{"MsgId inconsitent"}
	}
	return reply.Nodes, nil
}

// For project 2!
func (kad *Kademlia) DoIterativeFindNode(id ID) (C []Contact, e error) {
	list := new(ShortList)
	list.Init(kad, id)
	initnodes, _, err := kad.RT.FindNearestNode(id)
	if err != nil {
		return nil, err
	}
	list.MAdd(initnodes)

	quit := false
	for !quit {
		rpchwnd := make([]*rpc.Call, 0)
		alphacontacts := list.GetNearestN(alpha)
		for i := 0; i < len(alphacontacts); i++ {
			hwnd, _ := kad.DoFindNodeAsync(&alphacontacts[i], id)
			rpchwnd = append(rpchwnd, hwnd)
		}
		// TODO: Is it the closet node or closet active node?
		olddist := list.ClosetNode.Dist
		for i := 0; i < len(alphacontacts); i++ {
			Ret, err := kad.DoFindNodeWait(rpchwnd[i])
			if err != nil {
				/* Not responding */
				list.Remove(alphacontacts[i].NodeID)
			} else {
				list.SetActive(alphacontacts[i].NodeID)
				list.MAdd(Ret)
			}
		}
		newdist := list.ClosetNode.Dist

		if list.ActiveSize() >= k {
			quit = true
			continue
		}

		if newdist >= olddist {
			quit = true
			// TODO: Spec not clear
			// Should we terminate here, or keep contact those node not contacted?
			continue
		}
	}

	active := list.GetActiveContact()
	if len(active) > k {
		active = active[:k]
	}

	return active, nil
}

func (k *Kademlia) DoIterativeStore(key ID, value []byte) (received []Contact, e error) {
	C, err := k.DoIterativeFindNode(key)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(C); i++ {
		ret := k.DoStore(&C[i], key, value)
		if ret == nil {
			received = append(received, C[i])
		}
	}

	return received, nil
}

type FindValueResultPair struct {
	res FindValueResult
	index int
}

func (k *Kademlia) rpcCallFanIn(chans []chan *rpc.Call, out chan FindValueResultPair) {
	for i, ch := range chans {
		go func(i int, in chan *rpc.Call) {
			for ret := range in {
				err := ret.Error
				if err != nil {
					out<-FindValueResultPair{
						FindValueResult{NewRandomID(),nil,nil, err},
						i}
				}
				req, ok := ret.Args.(FindValueRequest)
				if !ok {
					out<-FindValueResultPair{
						FindValueResult{NewRandomID(), nil, nil,
							&CommandFailed{"Unknown Error"}},
						i}
				}
				res, ok := ret.Args.(FindValueResult)
				if !ok {
					out<-FindValueResultPair{
						FindValueResult{NewRandomID(), nil, nil,
							&CommandFailed{"Unknown Error"}},
						i}
				}
				if req.MsgID != res.MsgID {
					out<-FindValueResultPair{
						FindValueResult{NewRandomID(), nil, nil,
							&CommandFailed{"MsgId inconsitent"}},
						i}
				}
				out<-FindValueResultPair{res, i}
			}
		}(i, ch)
	}
}

func (kadamlia *Kademlia) DoIterativeFindValue(key ID) (value []byte, err error) {
	list := new(ShortList)
	list.Init(kadamlia, key)
	initNode, _, err := kadamlia.RT.FindNearestNode(key)
	if err != nil {
		return nil, err
	}
	list.MAdd(initNode)

	quit := false
	for !quit {
		rpcChan := make([]chan *rpc.Call, 0)
		alphacontacts := list.GetNearestN(alpha)
		for _, contact := range alphacontacts {
			handle, _ := kadamlia.doFindValueAsync(&contact, key)
			rpcChan = append(rpcChan, handle.Done)
		}
		count := len(rpcChan)
		out := make(chan FindValueResultPair)
		kadamlia.rpcCallFanIn(rpcChan, out)

		oldDist := list.ClosetNode.Dist
		for count != 0 {
			select {
			case res := <-out:
				// TODO deal with result
				if res.res.Err != nil {
					list.Remove(alphacontacts[res.index].NodeID)
				}
				count -= 1
			}
		}
		newDist := list.ClosetNode.Dist

		if list.ActiveSize() >= k {
			quit = true
			continue
		}
		if newDist >= oldDist {
			quit = true
			continue
		}
	}


	return nil, &CommandFailed{"Not implemented"}
}

// For project 3!
func (k *Kademlia) Vanish(data []byte, numberKeys byte,
	threshold byte, timeoutSeconds int) (vdo VanashingDataObject) {
	return
}

func (k *Kademlia) Unvanish(searchKey ID) (data []byte) {
	return nil
}

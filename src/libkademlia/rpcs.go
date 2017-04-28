package libkademlia

// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"net"
)

type KademliaRPC struct {
	kademlia *Kademlia
}

// Host identification.
type Contact struct {
	NodeID ID
	Host   net.IP
	Port   uint16
}

///////////////////////////////////////////////////////////////////////////////
// PING
///////////////////////////////////////////////////////////////////////////////
type PingMessage struct {
	Sender Contact
	MsgID  ID
}

type PongMessage struct {
	MsgID  ID
	Sender Contact
}

func (k *KademliaRPC) Ping(ping PingMessage, pong *PongMessage) error {
	pong.MsgID = CopyID(ping.MsgID)
	// Specify the sende
	pong.Sender = k.kademlia.SelfContact
	// Update contact, etc
	k.kademlia.RT.Update(ping.Sender)
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// STORE
///////////////////////////////////////////////////////////////////////////////
type StoreRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
	Value  []byte
}

type StoreResult struct {
	MsgID ID
	Err   error
}

func (k *KademliaRPC) Store(req StoreRequest, res *StoreResult) error {
	res.MsgID = CopyID(req.MsgID)
	res.Err = k.kademlia.HT.Add(req.Key, req.Value)
	// Update contact
	k.kademlia.RT.Update(req.Sender)
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// FIND_NODE
///////////////////////////////////////////////////////////////////////////////
type FindNodeRequest struct {
	Sender Contact
	MsgID  ID
	NodeID ID
}

type FindNodeResult struct {
	MsgID ID
	Nodes []Contact
	Err   error
}

func (k *KademliaRPC) FindNode(req FindNodeRequest, res *FindNodeResult) error {
	//var resultCount int
	// Fill up result
	res.MsgID = CopyID(req.MsgID)
	res.Nodes, _, res.Err = k.kademlia.RT.FindNearestNode(req.NodeID)
	// Update contact
	k.kademlia.RT.Update(req.Sender)
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// FIND_VALUE
///////////////////////////////////////////////////////////////////////////////
type FindValueRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
}

// If Value is nil, it should be ignored, and Nodes means the same as in a
// FindNodeResult.
type FindValueResult struct {
	MsgID ID
	Value []byte
	Nodes []Contact
	Err   error
}

func (k *KademliaRPC) FindValue(req FindValueRequest, res *FindValueResult) error {
	// Fill up result
	res.MsgID = CopyID(req.MsgID)
	res.Nodes, _, res.Err = k.kademlia.RT.FindNearestNode(req.Key)
	res.Value, res.Err = k.kademlia.HT.Find(req.Key)
	if res.Err != nil {
		res.Value = nil
	}
	//update contact
	k.kademlia.RT.Update(req.Sender)
	return nil
}

// For Project 3

type GetVDORequest struct {
	Sender Contact
	VdoID  ID
	MsgID  ID
}

type GetVDOResult struct {
	MsgID ID
	VDO   VanashingDataObject
}

func (k *KademliaRPC) GetVDO(req GetVDORequest, res *GetVDOResult) error {
	// TODO: Implement.
	return nil
}

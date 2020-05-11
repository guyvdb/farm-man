package node



import (

)

// A node is a node on the network. A node has a unique serial number.
// When a node is first placed on the network, it uses a temporary serial
// number until the server assigns it a permanent serial number.


type Node struct {
	SerialNumber string
	
}

func NewNode(serial string) *Node {
	return &Node{
		SerialNumber: serial,
	}
}

func FindNode(serial string) *Node {
}




func (n *Node) Send(cmd 


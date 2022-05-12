package server

import (
	"bytes"
	K "consistent-hashing/key"
	"sort"
)

type Server struct {
	Name         string
	VirtualNodes []*VirtualNode

	// uuid
	Key          string
	RingPosition []byte
}

type VirtualNode struct {
	Key             string
	RingPosition    []byte
	OriginKey       string //origin is the key of original server
	KeyList         []*K.Key
	OriginServer    *Server
	NextVirtualNode *VirtualNode
}

// Sort the virtualNodes of a server
func (s *Server) SortServerVirtualNodes() {
	sort.Slice(s.VirtualNodes, func(i, j int) bool {
		return bytes.Compare(s.VirtualNodes[i].RingPosition, s.VirtualNodes[j].RingPosition) == -1
	})
}

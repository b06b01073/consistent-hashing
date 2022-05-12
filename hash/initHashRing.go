package hash

import (
	"bytes"
	K "consistent-hashing/key"
	"consistent-hashing/server"
	"crypto/sha1"
	"fmt"
	"os"
	"sort"

	"github.com/google/uuid"
)

const (
	initServerNumber  = 3
	initKeyNumber     = 10
	virtualNodeNumber = 4
)

type HashRing struct {
	Servers      []*server.Server
	Keys         []*K.Key
	VirtualNodes []*server.VirtualNode
}

func GetRing() *HashRing {

	if initServerNumber < 1 {
		fmt.Printf("You need to have at least one server.\n")
		os.Exit(1)
	}

	// slices are sorted by their position on hash ring
	hashRing := &HashRing{}

	hashRing.initServers()
	hashRing.initKeys()

	hashRing.distributeKey()

	return hashRing
}

// initialize a slice of server
func (h *HashRing) initServers() {
	serverSlice := make([]*server.Server, initServerNumber)

	for i := 0; i < initServerNumber; i++ {
		serverName := fmt.Sprintf("Server%d", i)
		serverSlice[i] = h.initServer(serverName)
	}

	h.sortVirtualNodes()
	h.getVirtualNodesLinkedList()
	h.Servers = serverSlice
}

//initialize a single server
func (h *HashRing) initServer(serverName string) *server.Server {
	key := uuid.New().String()
	ringPosition := getRingPosition(key)
	virtualNodes := make([]*server.VirtualNode, virtualNodeNumber)

	s := &server.Server{
		Name:         serverName,
		Key:          key,
		RingPosition: ringPosition,
	}

	for j := 0; j < virtualNodeNumber; j++ {
		virtualNodeKey := uuid.New().String()
		virtualNodeRingPosition := getRingPosition(virtualNodeKey)

		virtualNode := &server.VirtualNode{
			Key:          virtualNodeKey,
			RingPosition: virtualNodeRingPosition,
			OriginServer: s,
		}

		virtualNodes[j] = virtualNode
	}

	s.VirtualNodes = virtualNodes
	s.SortServerVirtualNodes()
	h.VirtualNodes = append(h.VirtualNodes, s.VirtualNodes...)

	return s
}

func (h *HashRing) initKeys() {
	keySlice := make([]*K.Key, initKeyNumber)

	for i := 0; i < initKeyNumber; i++ {
		key := uuid.New().String()
		ringPosition := getRingPosition(key)
		keySlice[i] = &K.Key{
			Key:          key,
			RingPosition: ringPosition,
		}
	}

	sort.Slice(keySlice, func(i int, j int) bool {
		return bytes.Compare(keySlice[i].RingPosition, keySlice[j].RingPosition) == -1
	})

	h.Keys = keySlice
}

func getRingPosition(key string) []byte {

	hashFunc := sha1.New()
	hashFunc.Write([]byte(key))
	ringPosition := hashFunc.Sum(nil)

	return ringPosition
}

func (h *HashRing) distributeKey() {

	virtualNodesIndex := 0
	keyIndex := 0

	// users in a server is sorted by this method
	for virtualNodesIndex < len(h.VirtualNodes) && keyIndex < len(h.Keys) {
		res := bytes.Compare(h.Keys[keyIndex].RingPosition, h.VirtualNodes[virtualNodesIndex].RingPosition)

		if res == -1 {
			h.VirtualNodes[virtualNodesIndex].KeyList = append(h.VirtualNodes[virtualNodesIndex].KeyList, h.Keys[keyIndex])

			keyIndex++
		} else {
			virtualNodesIndex++
		}
	}

	//剩下的key position > 最大的server position，根據ring的結構剩下的keys由ring上的第一個virtual node負責
	for ; keyIndex < len(h.Keys); keyIndex++ {
		h.VirtualNodes[0].KeyList = append(h.VirtualNodes[0].KeyList, h.Keys[keyIndex])
	}
}

func (h *HashRing) getVirtualNodesLinkedList() {
	n := len(h.VirtualNodes)
	for i := 0; i < n; i++ {
		h.VirtualNodes[i].NextVirtualNode = h.VirtualNodes[(i+1)%n]
	}
}

// Sort the virtualNodes of the hash ring
func (h *HashRing) sortVirtualNodes() {
	sort.Slice(h.VirtualNodes, func(i int, j int) bool {
		return bytes.Compare(h.VirtualNodes[i].RingPosition, h.VirtualNodes[j].RingPosition) == -1
	})
}

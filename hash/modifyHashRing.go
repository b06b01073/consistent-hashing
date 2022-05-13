package hash

import (
	"bytes"
	"consistent-hashing/server"
	"fmt"
)

func (h *HashRing) AddServer(serverName string) {
	server := h.initServer(serverName)
	h.Servers = append(h.Servers, server)
	h.sortVirtualNodes()
	server.SortServerVirtualNodes()
	h.getVirtualNodesLinkedList()

	//redistribute方式: 從插入的節點逆時針找到下一個節點，把中間遇到的user加到自己的KeyList並從順時針方向會遇到的下一個節點刪除，可以用slicing解決
	for _, virtualNode := range server.VirtualNodes {
		nextVirtualNode := virtualNode.NextVirtualNode

		/**
		找到第一個ringPosition嚴格大於virtualNode ringPosition的key
		loop完以後i為第一個ringPosition大於virtualNode
		所以keyList[:i]應該要歸屬於新的virtualNode
		**/
		for i := 0; i < len(nextVirtualNode.KeyList); i++ {
			if bytes.Compare(virtualNode.RingPosition, nextVirtualNode.KeyList[i].RingPosition) == -1 {
				virtualNode.KeyList = nextVirtualNode.KeyList[:i]
				nextVirtualNode.KeyList = nextVirtualNode.KeyList[i:]
				break
			}
		}
	}
}

func (h *HashRing) RemoveServer(key string) {
	if len(h.Servers) == 1 {
		fmt.Printf("Cannot remove all servers.\n")
		return
	}

	newHashRingVirtualNodes := make([]*server.VirtualNode, 0)

	flag := false

	for _, virtualNode := range h.VirtualNodes {
		if virtualNode.OriginServer.Key == key {
			virtualNode.NextVirtualNode.KeyList = append(virtualNode.KeyList, virtualNode.NextVirtualNode.KeyList...)
			flag = true
		} else {
			newHashRingVirtualNodes = append(newHashRingVirtualNodes, virtualNode)
		}
	}

	for i, server := range h.Servers {
		if server.Key == key {
			h.Servers = append(h.Servers[:i], h.Servers[i+1:]...)
			break
		}
	}

	if flag {
		h.VirtualNodes = newHashRingVirtualNodes
		h.getVirtualNodesLinkedList()
	} else {
		fmt.Printf("Key not found...\n")
	}
}

func (h *HashRing) ListServerInfo() {
	totalKeys := 0
	for _, server := range h.Servers {
		fmt.Println("------------------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Server name: %s\n", server.Name)
		fmt.Printf("Virtual nodes: %d\n", len(server.VirtualNodes))
		fmt.Printf("Server key: %s\n", server.Key)
		fmt.Printf("Server Position: %v\n\n", server.RingPosition)

		keyLen := 0
		for _, virtualNode := range server.VirtualNodes {
			keyLen += len(virtualNode.KeyList)
		}
		fmt.Printf("Connected by %d keys\n", keyLen)
		totalKeys += keyLen

		fmt.Printf("Virtual nodes connect: ")
		for _, v := range server.VirtualNodes {
			fmt.Printf("%d ", len(v.KeyList))
		}
		fmt.Printf("keys respectively\n")
	}
	fmt.Printf("\nTotal Keys: %d", totalKeys)
	fmt.Printf("\n\n")
}

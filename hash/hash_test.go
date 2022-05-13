package hash

import (
	"bytes"
	"math/rand"
	"testing"
)

const testingServerName = "testingServer"

var reset = "\033[0m"
var green = "\033[32m"
var red = "\033[31m"

func TestAddServer(t *testing.T) {
	hashRing := GetRing()
	hashRing.AddServer(testingServerName)
	if !isVirtaulNodesSorted(hashRing, t) {
		return
	}

	if !isTotalKeysNumberCorrect(hashRing, t) {
		return
	}
	t.Log(colorize(green, "Check pass"))
}

func TestRemoveServer(t *testing.T) {
	hashRing := GetRing()

	//remove a random server
	serverIndex := rand.Intn(len(hashRing.Servers))
	hashRing.RemoveServer(hashRing.Servers[serverIndex].Key)

	if !isVirtaulNodesSorted(hashRing, t) {
		return
	}

	if !isTotalKeysNumberCorrect(hashRing, t) {
		return
	}

	t.Log(colorize(green, "Check pass"))
}

func TestRemoveTilOneServer(t *testing.T) {
	hashRing := GetRing()
	n := len(hashRing.Servers)

	for i := n; i >= 1; i-- {
		serverIndex := rand.Intn(i)
		hashRing.RemoveServer(hashRing.Servers[serverIndex].Key)

		if !isTotalKeysNumberCorrect(hashRing, t) {
			t.Error(colorize(red, "Check failed: Virtual nodes are not sorted!"))
			return
		}

		if !isVirtaulNodesSorted(hashRing, t) {
			return
		}
	}

	t.Log(colorize(green, "Check pass"))
}

func isTotalKeysNumberCorrect(hashRing *HashRing, t *testing.T) bool {
	sum := 0
	for _, virtualNode := range hashRing.VirtualNodes {
		sum += len(virtualNode.KeyList)
	}

	flag := (sum == config.InitKeyNumber)
	if !flag {
		t.Error(colorize(red, "Check failed: Total number of keys is incorrect!"))
	}
	return flag
}
func isVirtaulNodesSorted(hashRing *HashRing, t *testing.T) bool {
	flag := true
	n := len(hashRing.VirtualNodes)
	for i := 0; i < n; i++ {
		curVirtualNode := hashRing.VirtualNodes[i]
		nextVirtualNode := hashRing.VirtualNodes[(i+1)%n]

		if curVirtualNode.NextVirtualNode != nextVirtualNode {
			flag = false
		}

		if i+1 < n && bytes.Compare(hashRing.VirtualNodes[i].RingPosition, hashRing.VirtualNodes[i+1].RingPosition) != -1 {
			flag = false
		}

		if !flag {
			break
		}
	}

	if !flag {
		t.Error(colorize(red, "Check failed: Virtual nodes are not sorted!"))
	}
	return flag
}

func colorize(color string, s string) string {
	return color + s + reset
}

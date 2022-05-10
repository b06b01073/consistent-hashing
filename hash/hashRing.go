package hash

import (
	"bytes"
	"consistent-hashing/server"
	"consistent-hashing/user"
	"crypto/sha1"
	"fmt"
	"sort"

	"github.com/google/uuid"
)

type HashRing struct {
	Servers []server.Server
	Users   []user.User
}

func GetRing(serverNum int, userNum int) *HashRing {

	// slices are sorted by their position on hash ring
	servers := getServers(serverNum)
	users := getUsers(userNum)

	hashRing := &HashRing{
		Servers: servers,
		Users:   users,
	}
	hashRing.distributeUsers()

	return hashRing
}

func (h *HashRing) AddServer(serverName string) {
	key := uuid.New().String()
	ringPosition := getRingPosition(key)

	h.Servers = append(h.Servers, server.Server{
		Name:         serverName,
		Key:          key,
		RingPosition: ringPosition,
		UserList:     make([]user.User, 0, 10),
	})

	sort.Slice(h.Servers, func(i int, j int) bool {
		return bytes.Compare(h.Servers[i].RingPosition, h.Servers[j].RingPosition) == -1
	})
	//redistribute方式: 從插入的節點順時針找到下一個節點，把中間遇到的user加到自己的userList並從順時針方向會遇到的下一個節點刪除，可以用slicing解決
	h.distributeUsers()
}

func (h *HashRing) RemoveServer(key string) {
	for i, server := range h.Servers {
		if server.Key == key {
			//redistribute方式: 從被刪除的節點往順時鐘方向找到下一個節點，把userList全部交給他
			h.Servers = append(h.Servers[:i], h.Servers[i+1:]...)
		}
	}
	h.distributeUsers()
}

func getServers(serverNum int) []server.Server {
	serverSlice := make([]server.Server, serverNum)

	for i := 0; i < serverNum; i++ {
		key := uuid.New().String()
		ringPosition := getRingPosition(key)
		serverName := fmt.Sprintf("Server %d", i)

		serverSlice[i] = server.Server{
			Name:         serverName,
			Key:          key,
			RingPosition: ringPosition,
			UserList:     make([]user.User, 0, 10),
		}
	}

	sort.Slice(serverSlice, func(i int, j int) bool {
		return bytes.Compare(serverSlice[i].RingPosition, serverSlice[j].RingPosition) == -1
	})

	return serverSlice
}

func getUsers(userNum int) []user.User {
	userSlice := make([]user.User, userNum)

	for i := 0; i < userNum; i++ {
		key := uuid.New().String()
		ringPosition := getRingPosition(key)
		userName := fmt.Sprintf("User %d", i)
		userSlice[i] = user.User{
			Name:         userName,
			Key:          key,
			RingPosition: ringPosition,
		}
	}

	sort.Slice(userSlice, func(i int, j int) bool {
		return bytes.Compare(userSlice[i].RingPosition, userSlice[j].RingPosition) == -1
	})

	return userSlice
}

func getRingPosition(key string) []byte {

	hashFunc := sha1.New()
	hashFunc.Write([]byte(key))
	ringPosition := hashFunc.Sum(nil)

	return ringPosition
}

func (h *HashRing) ListServerInfo() {
	for _, server := range h.Servers {
		fmt.Println("------------------------------------------------------------------------------------------------------------------------")
		// fmt.Printf("Server name: %s\n", server.Name)
		// fmt.Printf("Server key: %s\n", server.Key)
		// fmt.Printf("Server Position: %v\n\n", server.RingPosition)
		fmt.Printf("Connected by %d users\n", len(server.UserList))

		// Print when long info is suggest
		// for _, user := range server.UserList {
		// 	fmt.Printf("Username: %s, Key: %s, Position: %v\n", user.Name, user.Key, user.RingPosition)
		// }
	}
}

func (h *HashRing) distributeUsers() {

	for i := 0; i < len(h.Servers); i++ {
		h.Servers[i].UserList = make([]user.User, 0, 10)
	}

	serverIndex := 0
	userIndex := 0

	// users in a server is sorted by this method
	for serverIndex < len(h.Servers) && userIndex < len(h.Users) {
		res := bytes.Compare(h.Users[userIndex].RingPosition, h.Servers[serverIndex].RingPosition)

		if res == -1 {
			h.Servers[serverIndex].UserList = append(h.Servers[serverIndex].UserList, h.Users[userIndex])
			userIndex++
		} else {
			serverIndex++
		}
	}

	//剩下的user position > 最大的server position，根據ring的結構剩下的users由ring上的第一個server負責
	for ; userIndex < len(h.Users) && len(h.Servers) > 0; userIndex++ {
		h.Servers[0].UserList = append(h.Servers[0].UserList, h.Users[userIndex])
	}
}

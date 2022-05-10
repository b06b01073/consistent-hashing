package server

import "consistent-hashing/user"

type Server struct {
	Name string
	// VirtualNodes []string

	// uuid
	Key          string
	RingPosition []byte
	UserList     []user.User
}

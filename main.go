package main

import (
	cli "consistent-hashing/cli"
	hash "consistent-hashing/hash"
)

const (
	initServerNumber = 100
	initUser         = 100000
)

func main() {
	hashRing := hash.GetRing(initServerNumber, initUser)
	cli.Cli(hashRing)
}
